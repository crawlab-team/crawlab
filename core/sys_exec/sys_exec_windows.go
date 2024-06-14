//go:build windows
// +build windows

package sys_exec

import "os/exec"

func BuildCmd(cmdStr string) *exec.Cmd {
	return exec.Command("cmd", "/C", cmdStr)
}
