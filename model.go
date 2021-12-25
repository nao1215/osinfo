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

import "os/exec"

func model(os string) string {
	model := ""
	switch os {
	case "Linux":
		model = getLinuxModelName()
	}
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
