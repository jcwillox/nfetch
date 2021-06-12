package logo

import (
	_ "embed"
	"github.com/spf13/viper"
)

//go:embed logos/windows.txt
var windows string

//go:embed logos/windows_legacy.txt
var windowsLegacy string

func GetLogo() (string, []int) {
	logo := viper.GetString("logo")
	if logo == "legacy" {
		return windowsLegacy, []int{1, 2, 4, 3}
	} else {
		return windows, []int{4}
	}
}
