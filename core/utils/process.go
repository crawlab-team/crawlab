/*
 * Copyright (c) 2024. Core Digital Limited
 * 版权所有 (c) 2024. 重庆科锐数研科技有限公司 (Core Digital Limited)
 * All rights reserved. 保留所有权利。
 *
 * 该软件由 重庆科锐数研科技有限公司 (Core Digital Limited) 开发，未经明确书面许可，任何人不得使用、复制、修改或分发该软件的任何部分。
 * This software is developed by Core Digital Limited. No one is permitted to use, copy, modify, or distribute this software without explicit written permission.
 *
 * 许可证：
 * 该软件仅供授权使用。授权用户有权在授权范围内使用、复制、修改和分发该软件。
 * License:
 * This software is for authorized use only. Authorized users are permitted to use, copy, modify, and distribute this software within the scope of their authorization.
 *
 * 免责声明：
 * 该软件按“原样”提供，不附带任何明示或暗示的担保，包括但不限于对适销性和适用于特定目的的担保。在任何情况下，版权持有者或其许可方对因使用该软件而产生的任何损害或其他责任概不负责。
 * Disclaimer:
 * This software is provided "as is," without any express or implied warranties, including but not limited to warranties of merchantability and fitness for a particular purpose. In no event shall the copyright holder or its licensors be liable for any damages or other liability arising from the use of this software.
 *
 */

package utils

import (
	"github.com/crawlab-team/crawlab/trace"
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
