package sysinfo

import (
	"fmt"
	"github.com/reujab/wallpaper"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// cache host-info as it's commonly used
// we ignore synchronisation as we expect that host-info will
// always be called before starting any goroutines
var hostInfo *host.InfoStat

func HostInfo() *host.InfoStat {
	if hostInfo != nil {
		return hostInfo
	}
	info, err := host.Info()
	if err != nil {
		panic(fmt.Errorf("Fatal unable to retrieve host info: %s \n", err))
	}
	hostInfo = info
	return hostInfo
}

func Username() string {
	username := os.Getenv("USER")
	if username != "" {
		return username
	}
	username = os.Getenv("USERNAME")
	if username != "" {
		return username
	}
	// extract from homedir
	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Base(homeDir)
	}
	return "unknown"
}

func Usage() (float64, error) {
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return -1, err
	}
	if len(percent) > 0 {
		return percent[0], err
	}
	return -1, nil
}

func Memory() (*mem.VirtualMemoryStat, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return memory, nil
}

func Disk() ([]*disk.UsageStat, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	usages := make([]*disk.UsageStat, 0, 4)
	for _, partition := range partitions {
		if strings.HasPrefix(partition.Mountpoint, "/snap") || strings.HasPrefix(partition.Mountpoint, "/mnt/wsl/") {
			continue
		}
		usage, err := disk.Usage(partition.Mountpoint)
		if err == nil {
			usages = append(usages, usage)
		}

	}

	return usages, err
}

func PublicIP() (net.IP, error) {
	ip, err := publicIP()
	if ip != nil {
		return ip, err
	}

	// if we did not find the ip using dns fallback to using http
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ip, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(string(ip)), err
}

func LocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func Wallpaper() (string, error) {
	return wallpaper.Get()
}
