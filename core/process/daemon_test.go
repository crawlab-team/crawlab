package process

import (
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

func TestDaemon(t *testing.T) {
	d := NewProcessDaemon(func() *exec.Cmd {
		return exec.Command("echo", "hello")
	})
	err := d.Start()
	require.Nil(t, err)

	d = NewProcessDaemon(func() *exec.Cmd {
		return exec.Command("return", "1")
	})
	err = d.Start()
	require.NotNil(t, err)
}
