# How to use
## Sample code
```
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
	fmt.Println("Kenel name          :" + info.Kernel.Name)
	fmt.Println("Kernel version      :" + info.Kernel.Ver)
	fmt.Println("Kernel architecture :" + info.Kernel.Arch)
	fmt.Println("Uptime              :" + info.Uptime)
	fmt.Println("Mac name            :" + info.Mac.Name)
	fmt.Println("Mac version         :" + info.Mac.Ver)
	fmt.Println("Mac build version   :" + info.Mac.BuildVer)
}
```
## Result
```
OS                  :Linux
Distribution        :Ubuntu Budgie 21.10
Model(Host)         :Gigabyte Technology Co., Ltd. B450 I AORUS PRO WIFI-CF
Kenel name          :Linux
Kernel version      :5.13.0-22-generic
Kernel architecture :x86_64
Uptime              :2 days, 7 hours, 14 minutes
Mac name            :This is not mac
Mac version         :No version information
Mac build version   :No build information
```

# Why did I create the osinfo library
In order to implement the [neofetch](https://github.com/dylanaraps/neofetch) command in golang in another project ([mimixbox](https://github.com/nao1215/mimixbox)), it was necessary to port the function of neofetch (written by shell) to golang.
# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.  
- [GitHub Issue](https://github.com/nao1215/osinfo/issues)
- [Twitter@ARC_AED](https://twitter.com/ARC_AED)
# LICENSE
The osinfo project is licensed under the terms of the Apache License 2.0.  
See LICENSE.
