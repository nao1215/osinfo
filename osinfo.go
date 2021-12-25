//
// osinfo/osinfo.go
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
)

type utsname struct {
	sys     string
	node    string
	release string
	version string
	machine string
	domain  string
}

type macProductInfo struct {
	name     string
	version  string
	buildVer string
}

type OsInfo struct {
	Os            string
	Distro        string
	KernelName    string
	KernelVer     string
	KernelMachine string
	mac           macProductInfo
}

func Get() OsInfo {
	utsname := uts()

	osinfo := OsInfo{
		Os:            operatingSystem(utsname.sys),
		Distro:        distribution(operatingSystem(utsname.sys), utsname.sys, utsname.release),
		KernelName:    utsname.sys,
		KernelVer:     utsname.release,
		KernelMachine: utsname.machine,
		mac:           getMacProductInfo(),
	}

	fmt.Println(osinfo)
	return osinfo
}

func utsToString(f [65]int8) string {
	out := make([]byte, 0, 64)
	for _, v := range f[:] {
		if v == 0 {
			break
		}
		out = append(out, uint8(v))
	}
	return string(out)
}