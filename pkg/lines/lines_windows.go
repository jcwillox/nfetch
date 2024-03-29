//go:build windows
// +build windows

package lines

import (
	"fmt"
	"github.com/jcwillox/nfetch/pkg/sysinfo"
)

func Theme(config LineConfig) (string, error) {
	sysTheme, appTheme, err := sysinfo.Theme()
	if err != nil {
		return "", err
	}
	var sys string
	var app string
	if sysTheme == 0 {
		sys = "Dark"
	} else {
		sys = "Light"
	}
	if appTheme == 0 {
		app = "Dark"
	} else {
		app = "Light"
	}
	return fmt.Sprintf("System - %s, Apps - %s", sys, app), err
}
