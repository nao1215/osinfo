//go:build linux

//
// osinfo/osinfo_linux.go
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
	u := syscall.Utsname{}
	err := syscall.Uname(&u)
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
	return macProductInfo{
		name:     "This is not mac",
		version:  "No version information",
		buildVer: "No build information"}
}
