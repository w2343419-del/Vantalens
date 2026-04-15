package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"vantalens/talentwriter/internal/auth"
	"vantalens/talentwriter/internal/config"
	"vantalens/talentwriter/internal/models"
)

func HandleControlPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/platform/control" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cfg := config.GetConfig()
	writerURL := config.LocalhostURL(9091, "/platform/backend")
	if cfg != nil && cfg.WriterPort > 0 {
		writerURL = config.LocalhostURL(cfg.WriterPort, "/platform/backend")
	}
	_, _ = w.Write([]byte(ControlHTML(writerURL)))
}

func HandleWriterPageRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/platform/backend" {
		http.NotFound(w, r)
		return
	}
	cfg := config.GetConfig()
	writerURL := config.LocalhostURL(9091, "/platform/backend")
	if cfg != nil && cfg.WriterPort > 0 {
		writerURL = config.LocalhostURL(cfg.WriterPort, "/platform/backend")
	}
	http.Redirect(w, r, writerURL, http.StatusTemporaryRedirect)
}

func HandleBackendPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/platform/backend" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cfg := config.GetConfig()
	controlURL := config.LocalhostURL(9090, "/platform/control")
	if cfg != nil && cfg.ControlPort > 0 {
		controlURL = config.LocalhostURL(cfg.ControlPort, "/platform/control")
	}
	_, _ = w.Write([]byte(DashboardHTML("2.0.0", controlURL)))
}

func HandleControlStatus(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	hugoPath := ""
	adminSet := false
	if cfg != nil {
		hugoPath = cfg.HugoPath
		adminSet = cfg.AdminToken != ""
	}

	frontend := runCommand(hugoPath, 25*time.Second, "hugo", "version")
	backend := map[string]interface{}{
		"service":            "online",
		"platform":           runtime.GOOS,
		"hugo_path":          hugoPath,
		"admin_token_config": adminSet,
	}

	RespondJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]interface{}{
		"frontend": frontend,
		"backend":  backend,
	}})
}

type controlCommandRequest struct {
	Scope  string `json:"scope"`
	Action string `json:"action"`
}

func HandleControlCommand(w http.ResponseWriter, r *http.Request) {
	if !auth.RequireAuth(w, r) {
		return
	}

	var req controlCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid request"})
		return
	}

	req.Scope = strings.ToLower(strings.TrimSpace(req.Scope))
	req.Action = strings.ToLower(strings.TrimSpace(req.Action))
	if req.Scope == "" || req.Action == "" {
		RespondJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "scope and action are required"})
		return
	}

	cfg := config.GetConfig()
	hugoPath := "."
	if cfg != nil && strings.TrimSpace(cfg.HugoPath) != "" {
		hugoPath = cfg.HugoPath
	}

	result, err := executeControlCommand(cfg, hugoPath, req.Scope, req.Action)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error(), Data: result})
		return
	}

	RespondJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: result})
}

func executeControlCommand(cfg *config.Config, hugoPath, scope, action string) (map[string]interface{}, error) {
	switch scope {
	case "frontend":
		switch action {
		case "check":
			res := runCommand(hugoPath, 25*time.Second, "hugo", "version")
			if !res.Success {
				return map[string]interface{}{"scope": scope, "action": action, "result": res}, fmt.Errorf("hugo unavailable: %s", res.Output)
			}
			return map[string]interface{}{"scope": scope, "action": action, "result": res}, nil
		case "build":
			res := runCommand(hugoPath, 2*time.Minute, "hugo", "--minify")
			if !res.Success {
				return map[string]interface{}{"scope": scope, "action": action, "result": res}, fmt.Errorf("frontend build failed")
			}
			return map[string]interface{}{"scope": scope, "action": action, "result": res}, nil
		default:
			return nil, fmt.Errorf("unsupported frontend action: %s", action)
		}
	case "backend":
		switch action {
		case "check":
			return map[string]interface{}{
				"scope":  scope,
				"action": action,
				"result": map[string]interface{}{"service": "online", "time": time.Now().Format(time.RFC3339)},
			}, nil
		case "routes":
			return map[string]interface{}{
				"scope":  scope,
				"action": action,
				"result": []string{"/platform/control", "/platform/backend", "/api/posts", "/api/comments", "/api/settings", "/health", "/api/health"},
			}, nil
		case "stop_writer":
			writerPort := 9091
			if cfg != nil && cfg.WriterPort > 0 {
				writerPort = cfg.WriterPort
			}
			res := stopListenerOnPort(writerPort)
			if !res.Success {
				return map[string]interface{}{"scope": scope, "action": action, "result": res}, fmt.Errorf("failed to stop writer listener")
			}
			return map[string]interface{}{"scope": scope, "action": action, "result": res}, nil
		case "stop_control":
			controlPort := 9090
			if cfg != nil && cfg.ControlPort > 0 {
				controlPort = cfg.ControlPort
			}
			result := map[string]interface{}{
				"scope":  scope,
				"action": action,
				"result": map[string]interface{}{
					"message": "control service will stop shortly",
					"port":    controlPort,
				},
			}
			go func() {
				time.Sleep(400 * time.Millisecond)
				os.Exit(0)
			}()
			return result, nil
		default:
			return nil, fmt.Errorf("unsupported backend action: %s", action)
		}
	default:
		return nil, fmt.Errorf("unsupported scope: %s", scope)
	}
}

func stopListenerOnPort(port int) commandResult {
	if port <= 0 {
		return commandResult{Success: false, Output: "invalid port"}
	}
	if runtime.GOOS == "windows" {
		script := "$conns = Get-NetTCPConnection -LocalPort " + strconv.Itoa(port) + " -State Listen -ErrorAction SilentlyContinue; if ($conns) { $pids = $conns | Select-Object -ExpandProperty OwningProcess -Unique; Stop-Process -Id $pids -Force; Write-Output ('Stopped PIDs: ' + ($pids -join ', ')) } else { Write-Output 'No listener found.' }"
		return runCommand(".", 12*time.Second, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	}
	cmd := "pids=$(lsof -t -iTCP:" + strconv.Itoa(port) + " -sTCP:LISTEN 2>/dev/null); if [ -n \"$pids\" ]; then kill -9 $pids; echo \"Stopped PIDs: $pids\"; else echo \"No listener found.\"; fi"
	return runCommand(".", 12*time.Second, "sh", "-c", cmd)
}

type commandResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
}

func runCommand(dir string, timeout time.Duration, name string, args ...string) commandResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if output == "" && err != nil {
		output = err.Error()
	}

	if ctx.Err() != nil {
		return commandResult{Success: false, Output: "command timeout"}
	}
	if err != nil {
		return commandResult{Success: false, Output: output}
	}
	return commandResult{Success: true, Output: output}
}
