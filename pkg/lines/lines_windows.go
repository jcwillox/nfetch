// +build windows

package lines

import (
	"fmt"
	. "nfetch/pkg"
	"nfetch/pkg/sysinfo"
)

func Motherboard(config LineConfig) (string, error) {
	info, err := sysinfo.Motherboard()
	if err != nil {
		return "", err
	}
	return *info.Manufacturer + " " + *info.Product, nil
}

func GPU(config LineConfig) (string, error) {
	gpu, err := sysinfo.GPU()
	if err != nil {
		return "", err
	}
	return *gpu.Name, nil
}

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

func Model(config LineConfig) (string, error) {
	model, err := sysinfo.Model()
	if err != nil {
		return "", err
	}
	return *model.Manufacturer + " " + *model.Model, nil
}