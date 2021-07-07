package sysinfo

import (
	"bytes"
	"fmt"
	"os/exec"
)

func shell(name string, exe string) (string, string, error) {
	if name == "powershell" {
		name = "PowerShell"

		output, err := exec.Command(exe, "-NoProfile", "-NoLogo", "-NonInteractive", "-Command", "$PSVersionTable.PSVersion.ToString()").Output()
		if err != nil {
			return name, "", err
		}

		return name, string(bytes.TrimSuffix(output, []byte("\r\n"))), nil
	}

	// use friendly names
	switch name {
	case "pwsh":
		name = "PowerShell"
	case "cmd":
		name = "Command Prompt"
	}

	// get version from file details
	ver, err := GetFileVersion(exe)
	if err != nil {
		return name, "", err
	}

	var version string
	if ver.Build != 0 {
		version = fmt.Sprintf("%d.%d.%d.%d", ver.Major, ver.Minor, ver.Patch, ver.Build)
	} else {
		version = fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
	}

	return name, version, err
}
