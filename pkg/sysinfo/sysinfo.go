package sysinfo

import (
	"bytes"
	"github.com/Xuanwo/go-locale"
	"github.com/kbinani/screenshot"
	"github.com/reujab/wallpaper"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/language"
	"image"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func NumProcs() (int, error) {
	pids, err := process.Pids()
	if err != nil {
		return 0, err
	}
	return len(pids), nil
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

func Resolution() []image.Rectangle {
	n := screenshot.NumActiveDisplays()
	displays := make([]image.Rectangle, n)
	for i := 0; i < n; i++ {
		displays[i] = screenshot.GetDisplayBounds(i)
	}
	return displays
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

func Locale() (language.Tag, error) {
	return locale.Detect()
}

func Weather() (string, error) {
	resp, err := http.Get("http://wttr.in/?format=%t+-+%C+(%l)")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body = bytes.TrimPrefix(body, []byte("+"))
	return string(body), nil
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
