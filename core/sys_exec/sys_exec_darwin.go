//go:build darwin
// +build darwin

package sys_exec

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

func BuildCmd(cmdStr string) (cmd *exec.Cmd, err error) {
	if cmdStr == "" {
		return nil, errors.New("command string is empty")
	}
	args := strings.Split(cmdStr, " ")
	return exec.Command(args[0], args[1:]...), nil
}

func SetPgid(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}
}
