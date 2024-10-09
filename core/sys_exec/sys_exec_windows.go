//go:build windows
// +build windows

package sys_exec

import (
	"errors"
	"os/exec"
	"strings"
)

func BuildCmd(cmdStr string) (cmd *exec.Cmd, err error) {
	if cmdStr == "" {
		return nil, errors.New("command string is empty")
	}
	args := strings.Split(cmdStr, " ")
	return exec.Command(args[0], args[1:]...), nil
}
