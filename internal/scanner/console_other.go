//go:build !windows

package scanner

import "os/exec"

func hideConsoleWindow(cmd *exec.Cmd) {}
