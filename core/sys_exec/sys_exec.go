package sys_exec

import (
	"bufio"
	"github.com/crawlab-team/go-trace"
	"github.com/shirou/gopsutil/process"
	"os/exec"
	"time"
)

type KillProcessOptions struct {
	Timeout time.Duration
	Force   bool
}

func KillProcess(cmd *exec.Cmd, opts *KillProcessOptions) error {
	// process
	p, err := process.NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		return err
	}

	// kill function
	killFunc := func(p *process.Process) error {
		return killProcessRecursive(p, opts.Force)
	}

	if opts.Timeout != 0 {
		// with timeout
		return killProcessWithTimeout(p, opts.Timeout, killFunc)
	} else {
		// without timeout
		return killFunc(p)
	}
}

func killProcessWithTimeout(p *process.Process, timeout time.Duration, killFunc func(*process.Process) error) error {
	go func() {
		if err := killFunc(p); err != nil {
			trace.PrintError(err)
		}
	}()
	for i := 0; i < int(timeout.Seconds()); i++ {
		ok, err := process.PidExists(p.Pid)
		if err == nil && !ok {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return killProcess(p, true)
}

func killProcessRecursive(p *process.Process, force bool) (err error) {
	// children processes
	cps, err := p.Children()
	if err != nil {
		return killProcess(p, force)
	}

	// iterate children processes
	for _, cp := range cps {
		if err := killProcessRecursive(cp, force); err != nil {
			return err
		}
	}

	return nil
}

func killProcess(p *process.Process, force bool) (err error) {
	if force {
		err = p.Kill()
	} else {
		err = p.Terminate()
	}
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func ConfigureCmdLogging(cmd *exec.Cmd, fn func(scanner *bufio.Scanner)) {
	stdout, _ := (*cmd).StdoutPipe()
	stderr, _ := (*cmd).StderrPipe()
	scannerStdout := bufio.NewScanner(stdout)
	scannerStderr := bufio.NewScanner(stderr)
	go fn(scannerStdout)
	go fn(scannerStderr)
}
