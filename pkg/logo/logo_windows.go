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

func getLogo(logo string) (string, []int) {
	if strings.HasPrefix(logo, "windows 10") {
		return windows10, []int{4}
	} else if strings.HasPrefix(logo, "windows 11") {
		return windows11, []int{12, 14, 4}
	}
	return windowsLegacy, []int{1, 2, 4, 3}
}
