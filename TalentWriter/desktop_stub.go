//go:build !desktop

package main

import "fmt"

// launchDesktopShell is a no-op in non-desktop builds.
func launchDesktopShell(url, title string, width, height int) (bool, error) {
	return false, fmt.Errorf("desktop shell is only available with -tags desktop")
}
