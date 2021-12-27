//
// osinfo/utility.go
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
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

const (
	// Unknown : Don't use this bits
	Unknown os.FileMode = 1 << (9 - iota)
	// Readable : readable bits
	Readable
	// Writable : writable bits
	Writable
	// Executable : executable bits
	Executable
)

func isFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (!stat.IsDir())
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (stat.IsDir())
}

func IsReadable(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && ((stat.Mode() & Readable) != 0)
}

func hasEnvVar(environmentVar string) bool {
	return !emptyStr(os.Getenv(environmentVar))
}

func emptyStr(str string) bool {
	return str == ""
}

func readFile(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func existCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func removeStringByRegexp(str string, pattern string) string {
	rep := regexp.MustCompile(pattern)
	return rep.ReplaceAllString(str, "")
}
