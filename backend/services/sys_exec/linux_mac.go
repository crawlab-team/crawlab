// +build !windows

package sys_exec

import (
	"os/exec"
	"syscall"
)

func BuildCmd(cmdStr string) *exec.Cmd {
	return exec.Command("sh", "-c", cmdStr)
}

func Setpgid(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}
}

func KillProcess(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
