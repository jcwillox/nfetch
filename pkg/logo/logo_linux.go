package logo

import (
	_ "embed"
	"github.com/shirou/gopsutil/host"
	"strings"
)

//go:embed logos/debian.txt
var debian string

//go:embed logos/ubuntu.txt
var ubuntu string

//go:embed logos/kali.txt
var kali string

//go:embed logos/alpine.txt
var alpine string

//go:embed logos/proxmox.txt
var proxmox string

//go:embed logos/linux.txt
var linux string

func getLogo(logo string) (string, []int) {
	logo, colors := fetchLogo(logo)
	if logo != "" {
		return logo, colors
	}
	// fallback to platform
	platform, _, _, _ := host.PlatformInformation()
	logo, colors = fetchLogo(platform)
	if logo != "" {
		return logo, colors
	}
	// fallback to generic linux logo
	return linux, []int{15, 8, 3}
}

func fetchLogo(logo string) (string, []int) {
	switch {
	case strings.HasPrefix(logo, "debian"):
		return debian, []int{1, 7, 3}
	case strings.HasPrefix(logo, "ubuntu"):
		return ubuntu, []int{1, 7, 3}
	case strings.HasPrefix(logo, "kali"):
		return kali, []int{4, 8}
	case strings.HasPrefix(logo, "alpine"):
		return alpine, []int{4, 5, 7, 6}
	case strings.HasPrefix(logo, "proxmox"):
		return proxmox, []int{7, 202}
	}
	return "", nil
}
