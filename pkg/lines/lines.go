package lines

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"nfetch/internal/color"
	. "nfetch/pkg"
	"nfetch/pkg/sysinfo"
	"strings"
)

const (
	LineTitle       = "title"
	LineDashes      = "dashes"
	LineBlank       = "blank"
	LineColorbar    = "colorbar"
	LineOS          = "os"
	LineHost        = "host"
	LineKernel      = "kernel"
	LineUptime      = "uptime"
	LineMotherboard = "motherboard"
	LinePkgs        = "pkgs"
	LineShell       = "shell"
	LineResolution  = "resolution"
	LineTerminal    = "terminal"
	LineTheme       = "theme"
	LineCPU         = "cpu"
	LineGPU         = "gpu"
	LineUsage       = "cpu_usage"
	LineMemory      = "memory"
	LineSwap        = "swap"
	LineDisk        = "disk"
	LineBattery     = "battery"
	LineLocale      = "locale"
	LineWeather     = "weather"
	LineLocalIP     = "local_ip"
	LinePublicIP    = "public_ip"
)

var funcMap = map[string]func(config LineConfig) (string, error){
	LineOS:          OS,
	LineHost:        Model,
	LineKernel:      Kernel,
	LineMotherboard: Motherboard,
	LineUptime:      Uptime,
	//LinePkgs:
	//LineShell:
	//LineResolution:
	//LineTerminal:
	LineTheme:  Theme,
	LineCPU:    CPU,
	LineGPU:    GPU,
	LineUsage:  Usage,
	LineMemory: Memory,
	LineSwap:   Swap,
	//LineBattery:
	//LineLocale:
	LineWeather:  Weather,
	LineLocalIP:  LocalIP,
	LinePublicIP: PublicIP,
}

var defaultTitleMap = map[string]string{
	LineOS:          "OS",
	LineHost:        "Host",
	LineKernel:      "Kernel",
	LineMotherboard: "Motherboard",
	LineUptime:      "Uptime",
	LinePkgs:        "Packages",
	LineShell:       "Shell",
	LineResolution:  "Resolution",
	LineTerminal:    "Terminal",
	LineTheme:       "Theme",
	LineCPU:         "CPU",
	LineGPU:         "GPU",
	LineUsage:       "CPU Usage",
	LineMemory:      "Memory",
	LineSwap:        "Swap",
	LineWeather:     "Weather",
	LineLocalIP:     "Local IP",
	LinePublicIP:    "Public IP",
}

var allLines = []interface{}{
	LineTitle,
	LineDashes,
	LineOS,
	LineHost,
	LineKernel,
	LineMotherboard,
	LineUptime,
	LinePkgs,
	LineShell,
	LineResolution,
	LineTerminal,
	LineTheme,
	LineCPU,
	LineGPU,
	LineUsage,
	LineMemory,
	LineSwap,
	LineDisk,
	LineBattery,
	LineLocale,
	LineWeather,
	LineLocalIP,
	LinePublicIP,
	LineBlank,
	LineColorbar,
}

func Title() (string, error) {
	return fmt.Sprintf("%s@%s", aurora.Colorize(sysinfo.Username(), color.Colors.C1), aurora.Colorize(sysinfo.HostInfo().Hostname, color.Colors.C1)), nil
}

func Dashes() (string, error) {
	return strings.Repeat("-", len(sysinfo.Username())+len(sysinfo.HostInfo().Hostname)+1), nil
}

func OS(config LineConfig) (string, error) {
	name := sysinfo.Distro()
	info := sysinfo.HostInfo()
	if strings.Contains(info.KernelVersion, "microsoft") {
		name += " on Windows 10"
	}
	return fmt.Sprintf("%s [%s]", name, sysinfo.HostInfo().KernelArch), nil
}

func Kernel(config LineConfig) (string, error) {
	return sysinfo.HostInfo().KernelVersion, nil
}

func Uptime(config LineConfig) (string, error) {
	return ToHumanTime(sysinfo.HostInfo().Uptime), nil
}

func Usage(config LineConfig) (string, error) {
	usage, err := sysinfo.Usage()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.f%% (%d processes)", usage, sysinfo.HostInfo().Procs), nil
}

func CPU(config LineConfig) (string, error) {
	name, err := sysinfo.CPUName()
	if err != nil {
		return "", err
	}
	return name, err
}

func Memory(config LineConfig) (string, error) {
	mem, err := sysinfo.Memory()
	if err != nil {
		return "", err
	}
	used, unit := BytesToHuman(float64(mem.Used), 2, "GiB")
	total, unitTotal := BytesToHuman(float64(mem.Total), 2, "GiB")
	return fmt.Sprintf("%s %s / %s %s (%.f%%)", used, unit, total, unitTotal, mem.UsedPercent), err
}

// Disk info is an exception to the rule and uses a different return type
// as it is handled manually so it can render multiple lines
func Disk(config LineConfig) (title []string, content []string, err error) {
	disks, err := sysinfo.Disk()
	if err != nil {
		return
	}

	for _, stat := range disks {
		val, unit := BytesToHuman(float64(stat.Used), 1, "TiB")
		max, unitMax := BytesToHuman(float64(stat.Total), 1, "TiB")

		title = append(title, fmt.Sprintf("Disk (%s)", stat.Path))
		content = append(content, fmt.Sprintf("%s %s / %s %s (%.f%%)", val, unit, max, unitMax, stat.UsedPercent))
	}
	return
}

func Swap(config LineConfig) (string, error) {
	swap, err := sysinfo.Swap()
	if err != nil {
		return "", err
	}
	used, unit := BytesToHuman(float64(swap.Used), 2, "GiB")
	total, unitTotal := BytesToHuman(float64(swap.Total), 2, "GiB")
	return fmt.Sprintf("%s %s / %s %s (%.f%%)", used, unit, total, unitTotal, swap.UsedPercent), err
}

func Weather(config LineConfig) (string, error) {
	return sysinfo.Weather()
}

func LocalIP(config LineConfig) (string, error) {
	ip, err := sysinfo.LocalIP()
	if err != nil {
		return "", err
	}
	return ip.String(), err
}

func PublicIP(config LineConfig) (string, error) {
	ip, err := sysinfo.PublicIP()
	if err != nil {
		return "", err
	}
	return ip.String(), err
}
