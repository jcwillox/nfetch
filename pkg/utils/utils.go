package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// StripWSLPath removes all windows paths from the PATH on WSL as it includes many networked windows
// directories which are very slow to access
func StripWSLPath() {
	_, isWSL := os.LookupEnv("WSL_DISTRO_NAME")
	if !isWSL {
		return
	}
	path := os.Getenv("PATH")
	var newPath []string
	for _, line := range strings.Split(path, string(os.PathListSeparator)) {
		if !strings.HasPrefix(line, "/mnt/c") {
			newPath = append(newPath, line)
		}
	}
	os.Setenv("PATH", strings.Join(newPath, string(os.PathListSeparator)))
}

func StartCommand(name string, arg ...string) (*bufio.Reader, error) {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func StringSliceContains(slice []string, s string) bool {
	for _, el := range slice {
		if el == s {
			return true
		}
	}
	return false
}
