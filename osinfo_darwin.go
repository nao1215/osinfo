//go:build darwin

//
// osinfo/osinfo_darwin.go
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
	"syscall"
)

func uts() utsname {
	u := unix.Utsname{}
	err := unix.Uname(&u)
	if err != nil {
		return utsname{}
	}

	uname := utsname{
		sys:     utsToString(u.Sysname),
		node:    utsToString(u.Nodename),
		release: utsToString(u.Release),
		version: utsToString(u.Version),
		machine: utsToString(u.Machine),
		domain:  utsToString(u.Domainname),
	}
	return uname
}

func getMacProductInfo() macProductInfo {
	result, err := exec.Command("sw_vers").Output()
	if err != nil {
		return macProductInfo{}
	}
	productInfo := strings.Split(string(result), "\n")
	nameLine:=strings.Split(productInfo[0], ":")
	verLine:=strings.Split(productInfo[1], ":")
	buildVerLine:=strings.Split(productInfo[2], ":")

	return macProductInfo{
		name: strings.TrimSpace(nameLine[1])
		version:strings.TrimSpace(verLine[1])
		buildVer:strings.TrimSpace(buildVerLine[1])
	}
}
