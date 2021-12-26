//
// osinfo/example/example.go
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
package main

import (
	"fmt"

	"github.com/nao1215/osinfo"
)

func main() {
	osinfo := osinfo.Get()
	printAllInfo(osinfo)
}

func printAllInfo(info osinfo.OsInfo) {
	fmt.Println("OS                  :" + info.Os)
	fmt.Println("Distribution        :" + info.Distro)
	fmt.Println("Model(Host)         :" + info.Model)
	fmt.Println("Kenel name          :" + info.KernelName)
	fmt.Println("Kernel version      :" + info.KernelVer)
	fmt.Println("Kernel architecture :" + info.KernelArch)
	//fmt.Println(info.Mac.name)
	//fmt.Println(info.Mac.version)
	//fmt.Println(info.Mac.buildVer)
}
