//
// osinfo/model.go
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
	"os/exec"
	"regexp"
	"strings"
)

func model(os string, kernelArch string) string {
	model := ""
	switch os {
	case "Linux":
		model = getLinuxModelName()
	case "Mac OS X", "macOS":
		model = getMacModelName()
	case "iPhone OS":
		model = getiPhoneModelName(kernelArch)
	}
	model = strings.ReplaceAll(model, "\n", "")
	return model
}

func getLinuxModelName() string {
	model := ""
	if isAndroid() {
		model = androidModelName()
	} else if hasBoardInfoFile() {
		model = boardInfo()
	} else if hasProductInfoFile() {
		model = productInfo()
	} else if hasFirmwareInfoFile() {
		model = firmwareInfo()
	} else if hasSysinfoModelFile() {
		model = sysinfoModelFile()
	}
	return model
}

func getMacModelName() string {
	if isHackintosh() {
		return hackintoshModelName()
	}
	return macModelName()
}

func getiPhoneModelName(kernelArch string) string {
	model := "Unknown"
	series := map[string]string{
		"iPad1,1":            "iPad",
		"iPad2,[1-4]":        "iPad 2",
		"iPad3,[1-3]":        "iPad 3",
		"iPad3,[4-6]":        "iPad 4",
		"iPad6,1[12]":        "iPad 5",
		"iPad7,[5-6]":        "iPad 6",
		"iPad7,1[12]":        "iPad 7",
		"iPad11,[67]":        "iPad 8",
		"iPad4,[1-3]":        "iPad Air",
		"iPad5,[3-4]":        "iPad Air 2",
		"iPad11,[3-4]":       "iPad Air 3",
		"iPad13,[1-2]":       "iPad Air 4",
		"iPad6,[7-8]":        "iPad Pro (12.9 Inch)",
		"iPad6,[3-4]":        "iPad Pro (9.7 Inch)",
		"iPad7,[1-2]":        "iPad Pro 2 (12.9 Inch)",
		"iPad7,[3-4]":        "iPad Pro (10.5 Inch)",
		"iPad8,[1-4]":        "iPad Pro (11 Inch)",
		"iPad8,[5-8]":        "iPad Pro 3 (12.9 Inch)",
		"iPad8,9 | iPad8,10": "iPad Pro 4 (11 Inch)",
		"iPad8,1[1-2]":       "iPad Pro 4 (12.9 Inch)",
		"iPad2,[5-7]":        "iPad mini",
		"iPad4,[4-6]":        "iPad mini 2",
		"iPad4,[7-9]":        "iPad mini 3",
		"iPad5,[1-2]":        "iPad mini 4",
		"iPad11,[1-2]":       "iPad mini 5",
		"iPhone1,1":          "iPhone",
		"iPhone1,2":          "iPhone 3G",
		"iPhone2,1":          "iPhone 3GS",
		"iPhone3,[1-3]":      "iPhone 4",
		"iPhone4,1":          "iPhone 4S",
		"iPhone5,[1-2]":      "iPhone 5",
		"iPhone5,[3-4]":      "iPhone 5c",
		"iPhone6,[1-2]":      "iPhone 5s",
		"iPhone7,2":          "iPhone 6",
		"iPhone7,1":          "iPhone 6 Plus",
		"iPhone8,1":          "iPhone 6s",
		"iPhone8,2":          "iPhone 6s Plus",
		"iPhone8,4":          "iPhone SE",
		"iPhone9,[13]":       "iPhone 7",
		"iPhone9,[24]":       "iPhone 7 Plus",
		"iPhone10,[14]":      "iPhone 8",
		"iPhone10,[25]":      "iPhone 8 Plus",
		"iPhone10,[36]":      "iPhone X",
		"iPhone11,2":         "iPhone XS",
		"iPhone11,[46]":      "iPhone XS Max",
		"iPhone11,8":         "iPhone XR",
		"iPhone12,1":         "iPhone 11",
		"iPhone12,3":         "iPhone 11 Pro",
		"iPhone12,5":         "iPhone 11 Pro Max",
		"iPhone12,8":         "iPhone SE 2020",
		"iPhone13,1":         "iPhone 12 Mini",
		"iPhone13,2":         "iPhone 12",
		"iPhone13,3":         "iPhone 12 Pro",
		"iPhone13,4":         "iPhone 12 Pro Max",
		"iPod1,1":            "iPod touch",
		"ipod2,1":            "iPod touch 2G",
		"ipod3,1":            "iPod touch 3G",
		"ipod4,1":            "iPod touch 4G",
		"ipod5,1":            "iPod touch 5G",
		"ipod7,1":            "iPod touch 6G",
		"iPod9,1":            "iPod touch 7G"}

	for k, v := range series {
		match, _ := regexp.MatchString(k, kernelArch)
		if match {
			model = v
			break
		}
	}
	return model
}

func androidModelName() string {
	modelName := ""
	brand, brandErr := exec.Command("getprop", "ro.product.brand").Output()
	model, modelErr := exec.Command("getprop", "ro.product.model").Output()

	if brandErr != nil && modelErr == nil {
		modelName = string(model)
	} else if brandErr == nil && modelErr != nil {
		modelName = string(brand)
	} else if brandErr == nil && modelErr == nil {
		modelName = string(brand) + " " + string(model)
	}
	return modelName
}

func hackintoshModelName() string {
	out, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err != nil {
		return "Hackintosh"
	}
	return "Hackintosh (SMBIOS: " + string(out)
}

func macModelName() string {
	out, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err != nil {
		return "Macintosh"
	}
	return string(out)
}

func hasBoardInfoFile() bool {
	return isFile("/sys/devices/virtual/dmi/id/board_vendor") ||
		isFile("/sys/devices/virtual/dmi/id/board_name")
}

func hasProductInfoFile() bool {
	return isFile("/sys/devices/virtual/dmi/id/board_vendor") ||
		isFile("/sys/devices/virtual/dmi/id/board_name")
}

func hasFirmwareInfoFile() bool {
	return isFile("/sys/firmware/devicetree/base/model")
}

func hasSysinfoModelFile() bool {
	return isFile("/tmp/sysinfo/model")
}

func boardInfo() string {
	return readFile("/sys/devices/virtual/dmi/id/board_vendor") + " " +
		readFile("/sys/devices/virtual/dmi/id/board_name")
}

func productInfo() string {
	return readFile("/sys/devices/virtual/dmi/id/product_name") + " " +
		readFile("/sys/devices/virtual/dmi/id/product_version")
}

func firmwareInfo() string {
	return readFile("/sys/firmware/devicetree/base/model")
}

func sysinfoModelFile() string {
	return readFile("/tmp/sysinfo/model")
}

func isHackintosh() bool {
	out, err := exec.Command("kextstat").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "FakeSMC") || strings.Contains(string(out), "VirtualSMC")
}
