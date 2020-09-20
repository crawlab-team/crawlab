// +build windows

package sys_exec

import (
	"os/exec"
)

func BuildCmd(cmdStr string) *exec.Cmd {
	return exec.Command("cmd", "/C", cmdStr)
}

func Setpgid(cmd *exec.Cmd) {
	return
}

func KillProcess(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
	}
	return nil
}
