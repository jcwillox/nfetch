// +build windows

package sysinfo

import (
	"github.com/mitchellh/go-homedir"
	"nfetch/pkg/utils"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var AllPkgManagers = []string{
	// "winget",  // disabled by default as its very slow
	"scoop",
	"choco",
}

func CountPkgs(manager string) (int, error) {
	switch manager {
	case "winget":
		return countCmdLines(0, func(line string) bool { return strings.HasSuffix(line, "winget\r\n") }, "cmd", "/c", "winget", "list")
	case "scoop":
		return pkgsScoop()
	case "choco":
		return pkgsChoco()
	}
	return 0, nil
}

func pkgsScoop() (int, error) {
	// assume scoop is in the home directory
	homeDir, err := homedir.Dir()
	if err != nil {
		return 0, err
	}
	pkgs, err := countDirEntries(filepath.Join(homeDir, "scoop/apps/*"), -1)
	if err == nil {
		return pkgs, nil
	}
	// not in home directory, search the path instead
	path, err := exec.LookPath("scoop")
	if err != nil {
		return 0, err
	}
	path = filepath.Dir(filepath.Dir(path))
	return countDirEntries(filepath.Join(homeDir, "apps/*"), -1)
}

func pkgsChoco() (int, error) {
	reader, err := utils.StartCommand("choco", "-l")
	if err != nil {
		return -1, err
	}

	var lastLine string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lastLine = line
	}
	parts := strings.SplitN(lastLine, " ", 2)
	if len(parts) > 1 {
		return strconv.Atoi(parts[0])
	}
	return -1, err
}
