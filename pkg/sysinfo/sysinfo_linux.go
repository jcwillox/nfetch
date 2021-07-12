// +build linux

package sysinfo

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net"
	"nfetch/pkg/utils"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func cpuInfo() (CPUInfo, error) {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		return CPUInfo{}, err
	}

	model := info[0].ModelName
	if model == "" {
		model = info[0].VendorID
	}

	var threads, cores int32 = 0, 0
	for _, stat := range info {
		threads += stat.Cores
		coreId, _ := strconv.Atoi(stat.CoreID)
		if int32(coreId) > cores {
			cores = int32(coreId)
		}
	}

	if cores == 0 {
		// something when wrong so fallback to thread count
		cores = threads
	} else {
		// coreId is 0-based so add 1 for count
		cores += 1
	}

	return CPUInfo{
		Model:   model,
		Cores:   cores,
		Threads: threads,
		Mhz:     info[0].Mhz,
	}, nil
}

func getPrettyName(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		if parts[0] == "PRETTY_NAME" || parts[0] == "DISTRIB_DESCRIPTION" {
			return strings.Trim(parts[1], "\"\n")
		}
	}
	return ""
}

func getProxmoxVersion() string {
	if _, err := exec.LookPath("pveversion"); err != nil {
		return ""
	}

	reader, err := utils.StartCommand("dpkg", "-s", "pve-manager")
	if err == nil {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(line, "Version:") {
				return "Proxmox VE " + strings.TrimRight(strings.TrimPrefix(line, "Version: "), "\n")
			}
		}
	}

	return "Proxmox VE"
}

func Distro() string {
	name := getProxmoxVersion()
	if name != "" {
		return name
	}
	name = getPrettyName("/etc/lsb-release")
	if name != "" {
		return name
	}
	out, err := exec.Command("lsb_release", "-sd").Output()
	if err == nil {
		return strings.TrimRight(string(out), "\n")
	}
	name = getPrettyName("/etc/os-release")
	if name != "" {
		return name
	}
	name = getPrettyName("/usr/lib/os-release")
	if name != "" {
		return name
	}
	name = getPrettyName("/etc/openwrt_release")
	if name != "" {
		return name
	}
	out, err = exec.Command("uname", "-o").Output()
	if err == nil {
		return strings.TrimRight(string(out), "\n")
	}
	return runtime.GOOS
}

func Motherboard() (MotherboardInfo, error) {
	return MotherboardInfo{}, nil
}

func Model() (ModelInfo, error) {
	model, err := model()
	if err != nil {
		return ModelInfo{}, err
	}
	model = strings.TrimSuffix(model, "\n")
	if strings.HasPrefix(model, "Standard PC") {
		model = fmt.Sprintf("KVM/QEMU (%s)", model)
	} else if strings.HasPrefix(model, "OpenBSD") {
		model = fmt.Sprintf("vmm (%s)", model)
	}
	return ModelInfo{Model: model}, nil
}

func model() (string, error) {
	model, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name")
	if err == nil {
		model = bytes.TrimSuffix(model, []byte("\n"))
		productVersion, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_version")
		if err == nil {
			model = append(model, ' ')
			model = append(model, productVersion...)
		}
		return string(model), nil
	}
	model, err = os.ReadFile("/sys/firmware/devicetree/base/model")
	if err == nil {
		return string(model), nil
	}
	model, err = os.ReadFile("/tmp/sysinfo/model")
	if err == nil {
		return string(model), nil
	}
	return "(unknown)", nil
}

func GPU() (string, error) {
	return "(not implemented)", nil
}

func Swap() (*mem.SwapMemoryStat, error) {
	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	return swap, nil
}

func publicIP() (net.IP, error) {
	target := "myip.opendns.com"
	server := "resolver1.opendns.com:53"

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, server)
		},
	}

	ip, err := r.LookupIP(context.Background(), "ip", target)
	if len(ip) > 0 {
		return ip[0], err
	}

	return nil, err
}
