package logo

import (
	_ "embed"
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

//go:embed logos/linux.txt
var linux string

func getLogo(logo string) (string, []int) {
	switch {
	case strings.HasPrefix(logo, "debian"):
		return debian, []int{1, 7, 3}
	case strings.HasPrefix(logo, "ubuntu"):
		return ubuntu, []int{1, 7, 3}
	case strings.HasPrefix(logo, "kali"):
		return kali, []int{4, 8}
	case strings.HasPrefix(logo, "alpine"):
		return alpine, []int{4, 5, 7, 6}
	}
	return linux, []int{15, 8, 3}
}
