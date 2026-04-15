package server

import (
    "fmt"
    "runtime"
    "strings"
    "os/exec"
    "time"
)

// StartParentWatcher returns a channel that fires once when parent process exits.
// It helps tie backend lifecycle to the launcher window/process on Windows.
func StartParentWatcher(parentPID int, interval time.Duration) <-chan struct{} {
    done := make(chan struct{}, 1)

    if parentPID <= 0 || interval <= 0 {
        return done
    }

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for range ticker.C {
            if !isProcessAlive(parentPID) {
                done <- struct{}{}
                close(done)
                return
            }
        }
    }()

    return done
}

func isProcessAlive(pid int) bool {
    if pid <= 0 {
        return false
    }

    if runtime.GOOS == "windows" {
        // tasklist output is stable and available on default Windows installations.
        cmd := exec.Command("cmd", "/c", fmt.Sprintf("tasklist /FI \"PID eq %d\" /FO CSV /NH", pid))
        out, err := cmd.CombinedOutput()
        if err != nil {
            return false
        }
        text := strings.TrimSpace(strings.ToLower(string(out)))
        if text == "" {
            return false
        }
        return !strings.Contains(text, "no tasks are running")
    }

    // Unix-like fallback using ps.
    cmd := exec.Command("sh", "-c", fmt.Sprintf("ps -p %d -o pid=", pid))
    out, err := cmd.CombinedOutput()
    if err != nil {
        return false
    }
    return strings.TrimSpace(string(out)) != ""
}
