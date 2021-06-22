// +build linux

package sysinfo

import (
	"nfetch/pkg/utils"
	"strings"
)

var AllPkgManagers = []string{
	"dpkg",
	"apk",
	"kiss",
	"cpt-list",
	"snap",
}

func CountPkgs(manager string) (int, error) {
	switch manager {
	case "dpkg":
		return countCmdLines(0, func(line string) bool { return strings.HasPrefix(line, "ii") }, "dpkg", "--list")
	case "apk":
		return countCmdLines(0, nil, "apk", "info")
	case "kiss":
		return countCmdLines(0, nil, "kiss", "l")
	case "cpt-list":
		return countCmdLines(0, nil, "cpt-list")
	case "snap":
		if proc, _ := utils.GetProcessByName("snapd"); proc == nil {
			return 0, nil
		}
		return countCmdLines(-1, nil, "snap", "list")
	}
	return 0, nil
}
