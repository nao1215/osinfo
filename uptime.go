//
// osinfo/uptime.go
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
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func getUptime(os string) string {
	sec := "0"
	switch os {
	case "Linux", "Windows", "MINIX":
		sec = uptimeSecForLinuxWinMinix()
	case "Mac OS X", "macOS", "iPhone OS", "BSD", "FreeMiNT":
		sec = uptimeSecForAppleBsdFreemint()
	case "Solaris":
		sec = uptimeSecForSolaris()
	case "AIX", "IRIX":
		sec = uptimeSecForAixIrix()
	case "Haiku":
		sec = uptimeSecForHaiku()
	}
	return secToUptime(strings.ReplaceAll(sec, "\n", ""))
}

func secToUptime(sec string) string {
	secInt, _ := strconv.Atoi(sec)

	d := secInt / 60 / 60 / 24
	h := secInt / 60 / 60 % 24
	m := secInt / 60 % 60

	uptime := ""
	if d == 0 {
		// nothing
	} else if d == 1 {
		uptime = uptime + strconv.Itoa(d) + " day"
	} else {
		uptime = uptime + strconv.Itoa(d) + " days"
	}

	if h == 0 {
		// nothing
	} else if h == 1 {
		if !emptyStr(uptime) {
			uptime = uptime + ","
		}
		uptime = uptime + " " + strconv.Itoa(h) + " hour"
	} else {
		if !emptyStr(uptime) {
			uptime = uptime + ","
		}
		uptime = uptime + " " + strconv.Itoa(h) + " hours"
	}

	if m == 0 {
		// nothing
	} else if m == 1 {
		if !emptyStr(uptime) {
			uptime = uptime + ","
		}
		uptime = uptime + " " + strconv.Itoa(m) + " hour"
	} else {
		if !emptyStr(uptime) {
			uptime = uptime + ","
		}
		uptime = uptime + " " + strconv.Itoa(m) + " hours"
	}

	return uptime
}

func canReadUptimeFile() bool {
	return isFile("/proc/uptime") && IsReadable("/proc/uptime")
}

func uptimeSecForLinuxWinMinix() string {
	sec := "0"
	if canReadUptimeFile() {
		sec = readFile("/proc/uptime")
		sec = removeStringByRegexp(sec, "\\..*")
	} else {
		boot, err := exec.Command("date", "-d\"$(uptime -s)\"", "+%s").Output()
		if err != nil {
			return sec
		}
		now, err := exec.Command("date", "+%s").Output()
		if err != nil {
			return sec
		}
		sec = diffNowAndBoot(string(now), string(boot))
	}
	return sec
}

func uptimeSecForAppleBsdFreemint() string {
	sec := "0"
	boot, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
	if err != nil {
		return sec
	}

	now, err := exec.Command("date", "+%s").Output()
	if err != nil {
		return sec
	}

	bootStr := strings.ReplaceAll(string(boot), "{ sec = ", "")
	bootStr = removeStringByRegexp(bootStr, ",.*")
	return diffNowAndBoot(string(now), bootStr)
}

func uptimeSecForSolaris() string {
	time, err := exec.Command("kstat", "-p", "unix:0:system_misc:snaptime").Output()
	if err != nil {
		return "0"
	}
	timeStr := strings.Split(string(time), "\n")
	return removeStringByRegexp(timeStr[1], ".*")
}

func uptimeSecForAixIrix() string {
	// time = 2-04:55:07  <-- 2day, 4hours, 55min, 7sec
	time, err := exec.Command("env", "LC_ALL=POSIX", "ps", "-o", "etime=", "-p", "1").Output()
	if err != nil {
		return "0"
	}

	timeStr := string(time)
	day := "0"
	hour := "0"
	if strings.Contains(timeStr, "-") {
		day = removeStringByRegexp(timeStr, "-.*")
		timeStr = removeStringByRegexp(timeStr, ".*-")
	}

	fmt.Println(timeStr)
	r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`)
	if r.MatchString(timeStr) {
		hour = timeStr[0:2]
		timeStr = timeStr[3:8]
	}

	hour = strings.TrimPrefix(hour, "0")
	timeStr = strings.TrimPrefix(timeStr, "0")
	sec := toSec(day, hour, removeStringByRegexp(timeStr, ":.*"), removeStringByRegexp(timeStr, ".*:"))
	return sec
}

func uptimeSecForHaiku() string {
	time, err := exec.Command("system_time").Output()
	if err != nil {
		return "0"
	}
	t, _ := strconv.Atoi(strings.ReplaceAll(string(time), "\n", ""))
	return strconv.Itoa(t / 1000000)
}

func diffNowAndBoot(now, boot string) string {
	nowInt, _ := strconv.Atoi(strings.ReplaceAll(now, "\n", ""))
	bootInt, _ := strconv.Atoi(strings.ReplaceAll(boot, "\n", ""))
	return strconv.Itoa(nowInt - bootInt)
}

func toSec(day string, hour string, min string, sec string) string {
	min = strings.TrimPrefix(min, "0")
	sec = strings.TrimPrefix(sec, "0")
	d, _ := strconv.Atoi(strings.ReplaceAll(day, "\n", ""))
	h, _ := strconv.Atoi(strings.ReplaceAll(hour, "\n", ""))
	m, _ := strconv.Atoi(strings.ReplaceAll(min, "\n", ""))
	s, _ := strconv.Atoi(strings.ReplaceAll(sec, "\n", ""))

	return strconv.Itoa(d*86400 + h*3600 + m*60 + s)
}
