package sysinfo

import (
	"github.com/jcwillox/nfetch/pkg/utils"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

type ShellInfo struct {
	Name    string
	Exe     string
	Version string
}

func Shell() ShellInfo {
	// fetch parent process as this should be the 'shell'
	proc, _ := process.NewProcess(int32(os.Getppid()))
	name, _ := proc.Name()
	exe, _ := proc.Exe()

	if exe == "" {
		// fallback to fetching from environment
		exe = os.Getenv("SHELL")
		if exe != "" {
			name = filepath.Base(exe)
		}
	} else if name == "" {
		name = filepath.Base(exe)
	}

	// strip file extension
	name = utils.GetFileName(name)

	name, version, _ := shell(name, exe)
	if version == "" {
		version, _ = shellVersion(name, exe)
	}

	return ShellInfo{
		Name:    name,
		Exe:     exe,
		Version: version,
	}
}

func shellVersion(name string, exe string) (string, error) {
	regexVersion, _ := regexp.Compile(`(\d\.[\d.]*)`)
	var version []byte
	var err error

	switch name {
	case "sh", "ash", "dash", "es":
		return "", nil
	case "bash":
		version, err = exec.Command(exe, "-c", "printf %s \"$BASH_VERSION\"").Output()
	case "ksh":
		version, err = exec.Command(exe, "-c", "printf %s \"$KSH_VERSION\"").Output()
	case "osh":
		version, err = exec.Command(exe, "-c", "printf %s \"$OIL_VERSION\"").Output()
	case "tcsh":
		version, err = exec.Command(exe, "-c", "printf %s \"$tcsh\"").Output()
	case "go":
		version, err = exec.Command(exe, "version").Output()
	default:
		version, err = exec.Command(exe, "--version").Output()
	}

	if err != nil {
		return "", err
	}

	return string(regexVersion.Find(version)), nil
}
