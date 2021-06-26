// +build linux

package sysinfo

import (
	"bufio"
	"context"
	"github.com/shirou/gopsutil/mem"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CPUName() (string, error) {
	return "unknown", nil
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

func Distro() string {
	out, err := exec.Command("pveversion").Output()
	if err == nil {
		return "Proxmox VE " + strings.SplitN(string(out), "/", 3)[1]
	}
	name := getPrettyName("/etc/lsb-release")
	if name != "" {
		return name
	}
	out, err = exec.Command("lsb_release", "-sd").Output()
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
