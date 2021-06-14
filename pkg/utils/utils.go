package utils

import (
	"os"
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
