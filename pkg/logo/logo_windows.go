package logo

import (
	_ "embed"
	"github.com/spf13/viper"
	"nfetch/pkg/sysinfo"
	"strings"
)

//go:embed logos/windows_10.txt
var windows10 string

//go:embed logos/windows_11.txt
var windows11 string

//go:embed logos/windows_legacy.txt
var windowsLegacy string

func GetLogo() (string, []int) {
	logo := viper.GetString("logo")
	if logo == "" {
		logo = sysinfo.Distro()
	}
	if strings.HasPrefix(logo, "Windows 10") {
		return windows10, []int{4}
	} else if strings.HasPrefix(logo, "Windows 11") {
		return windows11, []int{12, 14, 4}
	} else {
		return windowsLegacy, []int{1, 2, 4, 3}
	}
}
