package logo

import (
	_ "embed"
	"strings"
)

//go:embed logos/windows_10.txt
var windows10 string

//go:embed logos/windows_11.txt
var windows11 string

//go:embed logos/windows_legacy.txt
var windowsLegacy string

func GetDistroLogo(logo string) (string, []string) {
	if strings.HasPrefix(logo, "windows 10") {
		return windows10, []string{"14"}
	} else if strings.HasPrefix(logo, "windows 11") {
		return windows11, []string{"12", "14", "4"}
	}
	return windowsLegacy, []string{"1", "2", "4", "3"}
}
