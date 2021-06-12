package logo

import (
	_ "embed"
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
	distro := sysinfo.Distro()
	switch {
	case strings.HasPrefix(distro, "Ubuntu"):
		return ubuntu, []int{1, 7, 3}
	case strings.HasPrefix(distro, "Kali"):
		return kali, []int{4, 8}
	case strings.HasPrefix(distro, "Alpine"):
		return alpine, []int{4, 5, 7, 6}
	}
	return "", nil
}
