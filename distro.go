//
// osinfo/distro.go
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
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func distribution(os string, kernelName string, kernelVer string) string {
	distro := "Unknown"

	switch os {
	case "Linux", "BSD", "MINIX":
		distro = getDistroNameForBsdLinuxMinix(kernelName, kernelVer)
	}

	return distro
}

func getDistroNameForBsdLinuxMinix(kernelName string, kernelVer string) string {
	distro := ""
	if isBedrock() {
		distro = bedrock()
	} else if isRedstar() {
		distro = redstar()
	} else if isArmbian() {
		distro = armbian()
	} else if isSiduction() {
		distro = siduction()
	} else if isElbrus() {
		distro = elbrus()
	} else if isProxmoxVE() {
		distro = proxmox()
	} else if hasLsbRelease() {
		distro = distroInfoFromLsbRelease()
	} else if hasReleseFile() {
		distro = distroInfoFromReleaseFile()
	} else if isGoboLinux() {
		distro = gobo()
	} else if isSDE() {
		distro = sde()
	} else if isCrux() {
		distro = crux()
	} else if isSliTaz() {
		distro = slitaz()
	} else if isKSLinux() {
		distro = kslinux()
	} else if isAndroid() {
		distro = android()
	} else if isChromeOS() {
		distro = chromeOS()
	} else if isGuix() {
		distro = guix()
	} else if isOpenBSD(kernelName) {
		distro = openBSD()
	} else {
		distro = othres(kernelName, kernelVer)
	}

	if onWindows(kernelVer) {
		distro = distro + appendWindows()
	} else if onChrome() {
		distro = distro + appendChrome()
	}
	distro = formatDistroStr(distro)

	if isUbuntuFlavor(distro) {
		distro = ubuntuFlavor(distro)
	}
	return distro
}

func formatDistroStr(distro string) string {
	distro = strings.ReplaceAll(distro, "'", "")
	distro = strings.ReplaceAll(distro, "\"", "")
	return strings.ReplaceAll(distro, "NAME=", "")
}

func isBedrock() bool {
	return isFile("/bedrock/etc/bedrock-release") && !hasEnvVar("BEDROCK_RESTRICT")
}

func isRedstar() bool {
	return isFile("/etc/redstar-release")
}

func isArmbian() bool {
	return isFile("/etc/armbian-release")
}

func isSiduction() bool {
	return isFile("/etc/siduction-version")
}

func isElbrus() bool {
	return isFile("/etc/mcst_version")
}

func isProxmoxVE() bool {
	return existCmd("pveversion")
}

func hasLsbRelease() bool {
	return existCmd("lsb_release")
}

func hasReleseFile() bool {
	return isFile("/etc/os-release") || isFile("/usr/lib/os-release") ||
		isFile("/etc/openwrt_release") || isFile("/etc/lsb-release")
}

func isGoboLinux() bool {
	return isFile("/etc/GoboLinuxVersion")
}

func isSDE() bool {
	return isFile("/etc/SDE-VERSION")
}

func isCrux() bool {
	return existCmd("crux")
}

func isSliTaz() bool {
	return isFile("/etc/slitaz-release")
}

func isKSLinux() bool {
	return existCmd("kpt") && existCmd("kpm")
}

func isAndroid() bool {
	return isDir("/system/app/") && isDir("/system/priv-app")
}

func isChromeOS() bool {
	if !isFile("/etc/lsb-release") {
		return false
	}
	release := readFile("/etc/lsb-release")
	return strings.Contains(release, "CHROMEOS")
}

func isGuix() bool {
	return existCmd("guix")
}

func isOpenBSD(kernelName string) bool {
	return kernelName == "OpenBSD"
}

func onWindows(kernelVer string) bool {
	ver := readFile("/proc/version")
	return strings.Contains(kernelVer, "Microsoft") || strings.Contains(ver, "Microsoft")
}

func onChrome() bool {
	ver := readFile("/proc/version")
	return strings.Contains(ver, "chrome-bot") || isFile("/dev/cros_ec")
}

func isUbuntuFlavor(distro string) bool {
	return strings.Contains(distro, "Ubuntu")
}

func bedrock() string {
	distro := "Bedrock Linux"
	relase := readFile("/bedrock/etc/bedrock-release")
	if !emptyStr(relase) {
		distro = distro + " " + relase
	}
	return distro
}

func redstar() string {
	distro := "Red Star OS"
	reg := "[^0-9*]"
	release := readFile("/etc/redstar-release")
	list := regexp.MustCompile(reg).Split(release, -1)

	if len(list) >= 1 {
		return distro + " " + list[1]
	}
	return list[1]
}

func armbian() string {
	release := readFile("/etc/armbian-release")
	releaseList := strings.Split(release, "\n")

	distro := "Armbian"
	distroCode := ""
	distroVer := ""
	for _, v := range releaseList {
		if strings.HasPrefix(v, "DISTRIBUTION_CODENAME") {
			distroCode = getValue(v)
		} else if strings.HasPrefix(v, "VERSION") {
			distroVer = getValue(v)
		}
	}
	if distroVer == "" {
		return distro + " " + distroCode
	}
	return distro + " " + distroCode + " " + distroVer
}

func siduction() string {
	distro := "Siduction"
	out, err := exec.Command("lsb_release", "-sic").Output()
	if err != nil {
		return distro
	}
	return distro + " " + string(out)
}

func elbrus() string {
	distro := "OS Elbrus"
	ver := readFile("/etc/mcst_version")
	if !emptyStr(ver) {
		distro = distro + " " + ver
	}
	return distro
}

func proxmox() string {
	distro := "Proxmox VE"
	out, err := exec.Command("pveversion").Output()
	if err != nil {
		return distro
	}

	ver := strings.TrimPrefix(string(out), "pve-manager/")
	rep := regexp.MustCompile(`/.*`)
	ver = rep.ReplaceAllString(ver, "")
	return distro + " " + ver
}

func distroInfoFromLsbRelease() string {
	out, err := exec.Command("lsb_release", "-sd").Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func distroInfoFromReleaseFile() string {
	prettyName := ""
	desc := ""
	codeName := ""
	files := []string{
		"/etc/os-release",
		"/usr/lib/os-release",
		"/etc/openwrt_release",
		"/etc/lsb-release"}

	for _, v := range files {
		if !isFile(v) {
			continue
		}
		for _, line := range strings.Split(readFile(v), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME") {
				prettyName = getValue(line)
			} else if strings.HasPrefix(line, "DISTRIB_DESCRIPTION") {
				desc = getValue(line)
			} else if strings.HasPrefix(line, "UBUNTU_CODENAME") {
				desc = getValue(line)
			}
		}
		if !emptyStr(prettyName) || !emptyStr(desc) || !emptyStr(codeName) {
			break
		}
	}
	if emptyStr(prettyName) {
		return desc + " " + codeName
	}
	return prettyName + " " + codeName
}

func getValue(keyValue string) string {
	rep := regexp.MustCompile(`.*=`)
	return rep.ReplaceAllString(keyValue, "")
}

func gobo() string {
	return "GoboLinux " + readFile("/etc/GoboLinuxVersion")
}

func sde() string {
	return readFile("/etc/SDE-VERSION")
}

func crux() string {
	out, err := exec.Command("pveversion").Output()
	if err != nil {
		return "CRUX"
	}
	return string(out)
}

func slitaz() string {
	return "SliTaz " + readFile("/etc/slitaz-release")
}

func kslinux() string {
	return "KSLinux"
}

func android() string {
	out, err := exec.Command("getprop", "ro.build.version.release").Output()
	if err != nil {
		return "Android"
	}
	return "Android " + string(out)
}

func chromeOS() string {
	return "Chrome OS"
}

func guix() string {
	out, err := exec.Command("guix", "-V").Output()
	if err != nil {
		return "Guix System"
	}
	lines := strings.Split(string(out), "\n")

	return "Guix System " + strings.Split(lines[0], " ")[3]
}

func openBSD() string {
	out, err := exec.Command("sysctl", "-n", "kern.version").Output()
	if err != nil {
		return "OpenBSD"
	}
	elem := strings.Split(string(out), " ")
	return elem[0] + " " + elem[1] + " " + elem[2]
}

func othres(kernelName string, kernelVer string) string {
	distro := ""
	releaseFiles := releaseFiles()
	if len(releaseFiles) != 0 {
		if isFile("/etc/pacbsd-release") {
			distro = "PacBSD"
		}
		return distro
	}

	distro = kernelName + " " + kernelVer
	if strings.Contains(distro, "DragonFly") && !strings.Contains(distro, "DragonFlyBSD") {
		distro = strings.ReplaceAll(distro, "DragonFly", "DragonFlyBSD")
	}

	if isFile("/etc/pcbsd-lang") {
		distro = "PCBSD"
	} else if isFile("/etc/trueos-lang") {
		distro = "TrueOS"
	} else if isFile("/etc/hbsd-update.conf") {
		distro = "HardenedBSD"
	}
	return distro
}

func releaseFiles() []string {
	files, err := os.ReadDir("/etc")
	if err != nil {
		return []string{}
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f := filepath.Join("/etc", file.Name())
		if strings.Contains(f, "release") {
			paths = append(paths, f)
		}
	}
	return paths
}

func appendWindows() string {
	out, err := exec.Command("wmic.exe", "os", "get", "Version").Output()
	if err != nil {
		return " on Windows"
	}
	ver := strings.ReplaceAll(string(out), "Version", "")
	ver = strings.TrimSpace(ver)
	return " on Windows " + ver
}

func appendChrome() string {
	return " on Chrome OS"
}

func ubuntuFlavor(distro string) string {
	flavor := os.Getenv("XDG_CONFIG_DIRS")

	if strings.Contains(flavor, "studio") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Ubuntu Studio")
	} else if strings.Contains(flavor, "plasma") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Kubuntu")
	} else if strings.Contains(flavor, "mate") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Ubuntu MATE")
	} else if strings.Contains(flavor, "xubuntu") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Xubuntu")
	} else if strings.Contains(flavor, "Lubuntu") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Lubuntu")
	} else if strings.Contains(flavor, "budgie") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Ubuntu Budgie")
	} else if strings.Contains(flavor, "cinnamon") {
		distro = strings.ReplaceAll(distro, "Ubuntu", "Ubuntu Cinnamon")
	}
	return distro
}
