//
// osinfo/shell.go
//
// Copyright 2021 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package osinfo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getShell() string {
	shell := filepath.Base(os.Getenv("SHELL"))

	return shell + " " + getShellVer(shell)
}

func getShellVer(shell string) string {
	ver := ""
	switch shell {
	case "bash":
		ver = bashVer()
	case "sh", "ash", "dash", "es":
		//nothing
	case "dtksh", "tksh", "oksh", "mksh", "SKsh":
		ver = kshVer()
	case "osh":
		ver = oshVer()
	case "tcsh":
		ver = tcshVer()
	case "yash":
		ver = yashVer()
	case "nu":
		ver = nuShellVer()
	default:
		ver = otherShell()
	}
	return removeUnusedInfoFromVer(ver)
}

func removeUnusedInfoFromVer(ver string) string {
	ver = removeStringByRegexp(ver, ", version")
	ver = removeStringByRegexp(ver, "options.*")
	ver = removeStringByRegexp(ver, "\\(.*\\)")
	return ver
}

func bashVer() string {
	ver := os.Getenv("BASH_VERSION")
	if emptyStr(ver) {
		version, err := exec.Command("bash", "-c", "printf %s \"$BASH_VERSION\"").Output()
		if err != nil {
			return ""
		}
		ver = string(version)
	}

	return strings.TrimSpace(removeStringByRegexp(ver, "-.*"))
}

func kshVer() string {
	shell := os.Getenv("SHELL")
	version, err := exec.Command(shell, "-c", "printf %s \"$KSH_VERSION\"").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	ver := string(version)
	ver = removeStringByRegexp(ver, " .* KSH")
	return strings.TrimSpace(removeStringByRegexp(ver, "version|Version"))
}

func oshVer() string {
	ver := os.Getenv("OIL_VERSION")
	if emptyStr(ver) {
		version, err := exec.Command("bash", "-c", "printf %s \"$OIL_VERSION\"").Output()
		if err != nil {
			return ""
		}
		ver = string(version)
	}
	return strings.TrimSpace(ver)
}

func tcshVer() string {
	shell := os.Getenv("SHELL")
	version, err := exec.Command(shell, "-c", "printf %s $tcsh").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimSpace(string(version))
}

func yashVer() string {
	shell := os.Getenv("SHELL")
	version, err := exec.Command(shell, "--version").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	verList := strings.Split(string(version), "\n")
	ver := removeStringByRegexp(verList[0], "yash")
	ver = removeStringByRegexp(ver, "Yet another shell")
	ver = removeStringByRegexp(ver, "Version|version|バージョン")
	return strings.TrimSpace(ver)
}

func nuShellVer() string {
	shell := os.Getenv("SHELL")
	verion, err := exec.Command(shell, "-c \"version | get version\"").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	ver := removeStringByRegexp(string(verion), "nu")
	return strings.TrimSpace(ver)
}

func otherShell() string {
	shell := os.Getenv("SHELL")
	version, err := exec.Command(shell, "--version").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	verList := strings.Split(string(version), "\n")
	ver := removeStringByRegexp(verList[0], shell)
	return strings.TrimSpace(ver)
}
