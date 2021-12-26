//
// osinfo/os.go
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

import "regexp"

func operatingSystem(kernelName string) string {
	var os string = "Unknown"
	switch kernelName {
	case "Darwin":
		mac := getMacProductInfo()
		os = mac.Name
	case "SunOS":
		os = "Solaris"
	case "Haiku":
		os = "haiku"
	case "MINIX":
		os = "MINIX"
	case "AIX":
		os = "AIX"
	case "FreeMiNT":
		os = "FreeMiNT"
	default:
		if isIRIX(kernelName) {
			os = "IRIX"
		} else if isLinux(kernelName) {
			os = "Linux"
		} else if isBSD(kernelName) {
			os = "BSD"
		} else if isWindows(kernelName) {
			os = "Windows"
		}
	}
	return os
}

func isIRIX(kernelName string) bool {
	match, _ := regexp.MatchString("IRIX.*", kernelName)
	return match
}

func isLinux(kernelName string) bool {
	match, _ := regexp.MatchString("Linux|GNU.*", kernelName)
	return match
}

func isBSD(kernelName string) bool {
	match, _ := regexp.MatchString(".*BSD|DragonFly|Bitrig", kernelName)
	return match
}

func isWindows(kernelName string) bool {
	match, _ := regexp.MatchString("CYGWIN.*|MSYS.*|MINGW.*", kernelName)
	return match
}
