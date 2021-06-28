package utils

import (
	"bufio"
	"github.com/mitchellh/go-ps"
	"os"
	"os/exec"
	"path/filepath"
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

func GetProcessByName(name string) (ps.Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processes {
		if p.Executable() == name {
			return p, nil
		}
	}
	return nil, err
}

func GetFileName(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
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

func StripToEnd(s string, sep string) string {
	i := strings.Index(s, sep)
	if i < 0 {
		return s
	}
	if i == 0 {
		return s[0:0]
	}
	return s[0:i]
}

func HasPrefixMulti(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
