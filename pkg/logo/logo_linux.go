package logo

import (
	_ "embed"
	"github.com/spf13/viper"
	"nfetch/pkg/sysinfo"
	"strings"
)

//go:embed logos/ubuntu.txt
var ubuntu string

//go:embed logos/kali.txt
var kali string

//go:embed logos/alpine.txt
var alpine string

func GetLogo() (string, []int) {
	logo := viper.GetString("logo")
	if logo == "" {
		logo = sysinfo.Distro()
	}
	switch {
	case strings.HasPrefix(logo, "Ubuntu"):
		return ubuntu, []int{1, 7, 3}
	case strings.HasPrefix(logo, "Kali"):
		return kali, []int{4, 8}
	case strings.HasPrefix(logo, "Alpine"):
		return alpine, []int{4, 5, 7, 6}
	}
	return "", nil
}
