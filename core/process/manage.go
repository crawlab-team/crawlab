package process

import (
	"github.com/crawlab-team/go-trace"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var pidRegexp, _ = regexp.Compile("(?:^|\\s+)\\d+(?:$|\\s+)")

func ProcessIdExists(id int) (ok bool) {
	lines, err := ListProcess(string(rune(id)))
	if err != nil {
		return false
	}
	for _, line := range lines {
		matched := pidRegexp.MatchString(line)
		if matched {
			return true
		}
	}
	return false
}

func ListProcess(text string) (lines []string, err error) {
	if runtime.GOOS == "windows" {
		return listProcessWindow(text)
	} else {
		return listProcessLinuxMac(text)
	}
}

func listProcessWindow(text string) (lines []string, err error) {
	cmd := exec.Command("tasklist", "/fi", text)
	out, err := cmd.CombinedOutput()
	_, ok := err.(*exec.ExitError)
	if !ok {
		return nil, trace.TraceError(err)
	}
	lines = strings.Split(string(out), "\n")
	return lines, nil
}

func listProcessLinuxMac(text string) (lines []string, err error) {
	cmd := exec.Command("ps", "aux")
	out, err := cmd.CombinedOutput()
	_, ok := err.(*exec.ExitError)
	if !ok {
		return nil, trace.TraceError(err)
	}
	_lines := strings.Split(string(out), "\n")
	for _, l := range _lines {
		if strings.Contains(l, text) {
			lines = append(lines, l)
		}
	}
	return lines, nil
}
