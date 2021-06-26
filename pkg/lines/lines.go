package lines

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/shirou/gopsutil/host"
	"golang.org/x/text/language/display"
	"nfetch/internal/color"
	"nfetch/pkg/sysinfo"
	. "nfetch/pkg/utils"
	"runtime"
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
	LinePkgs:        Pkgs,
	//LineShell:
	LineResolution: Resolution,
	//LineTerminal:
	LineTheme:    Theme,
	LineCPU:      CPU,
	LineGPU:      GPU,
	LineUsage:    Usage,
	LineMemory:   Memory,
	LineSwap:     Swap,
	LineBattery:  Battery,
	LineLocale:   Locale,
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
	LineBattery:     "Battery",
	LineLocale:      "Locale",
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

func Title() string {
	return fmt.Sprintf("%s@%s", aurora.Colorize(sysinfo.Username(), color.Colors.C1), aurora.Colorize(sysinfo.Hostname(), color.Colors.C1))
}

func Dashes() string {
	return strings.Repeat("-", len(sysinfo.Username())+len(sysinfo.Hostname())+1)
}

func Colorbar() []string {
	return []string{
		"\x1b[0;40m   \x1b[0;41m   \x1b[0;42m   \x1b[0;43m   \x1b[0;44m   \x1b[0;45m   \x1b[0;46m   \x1b[0;47m   \x1b[0m",
		"\x1b[0;100m   \x1b[0;101m   \x1b[0;102m   \x1b[0;103m   \x1b[0;104m   \x1b[0;105m   \x1b[0;106m   \x1b[0;107m   \x1b[0m",
	}
}

func OS(config LineConfig) (string, error) {
	name := sysinfo.Distro()
	if runtime.GOOS != "windows" {
		kernelVersion, err := host.KernelVersion()
		if err != nil {
			return "", err
		}
		if strings.Contains(kernelVersion, "microsoft") {
			name += " on Windows 10"
		}
	}
	kernelArch, err := host.KernelArch()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s [%s]", name, kernelArch), nil
}

func Model(config LineConfig) (string, error) {
	model, err := sysinfo.Model()
	if err != nil {
		return "", err
	}
	result := model.Manufacturer
	if model.Model != "" {
		if model.Manufacturer != "" {
			result += " "
		}
		result += model.Model
	}
	return result, nil
}

func Kernel(config LineConfig) (string, error) {
	kernelVersion, err := host.KernelVersion()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		kernelVersion = StripToEnd(kernelVersion, " Build")
	}
	return kernelVersion, nil
}

func Motherboard(config LineConfig) (string, error) {
	info, err := sysinfo.Motherboard()
	if err != nil || info == (sysinfo.MotherboardInfo{}) {
		return "", err
	}
	return info.Manufacturer + " " + info.Product, nil
}

func Uptime(config LineConfig) (string, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return "", err
	}
	return ToHumanTime(uptime), nil
}

func Pkgs(config LineConfig) (string, error) {
	show := config.GetStringSlice("include")
	exclude := config.GetStringSlice("exclude")
	if len(exclude) > 0 {
		// we know that the slice will be at least this long
		show = make([]string, 0, len(sysinfo.AllPkgManagers)-len(exclude))
		for _, manager := range sysinfo.AllPkgManagers {
			if !StringSliceContains(exclude, manager) {
				show = append(show, manager)
			}
		}
	} else {
		if len(show) == 0 {
			show = sysinfo.AllPkgManagers
		}
	}

	pkgMap := sysinfo.Pkgs(show)
	if len(pkgMap) == 0 {
		return "(none)", nil
	}

	result := make([]string, 0, len(pkgMap))
	for _, pkg := range show {
		amt, ok := pkgMap[pkg]
		if !ok {
			continue
		}
		result = append(result, fmt.Sprintf("%d (%s)", amt, pkg))
	}
	return strings.Join(result, ", "), nil
}

func Resolution(config LineConfig) (string, error) {
	displays := sysinfo.Resolution()
	if len(displays) == 0 {
		return "(none)", nil
	}
	parts := make([]string, len(displays))
	for i, bounds := range displays {
		parts[i] = fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
	}
	return strings.Join(parts, ", "), nil
}

func Usage(config LineConfig) (string, error) {
	usage, err := sysinfo.Usage()
	if err != nil {
		return "", err
	}
	procs, err := sysinfo.NumProcs()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.f%% (%d processes)", usage, procs), nil
}

func CPU(config LineConfig) (string, error) {
	name, err := sysinfo.CPUName()
	if err != nil {
		return "", err
	}
	return name, err
}

func GPU(config LineConfig) (string, error) {
	return sysinfo.GPU()
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

func Battery(config LineConfig) (string, error) {
	bt, err := sysinfo.Battery()
	if err != nil {
		return "", err
	}
	if bt == nil {
		return "(none)", nil
	}
	return fmt.Sprintf("%.f%% (%s)", bt.Current/bt.Full*100, bt.State), nil
}

func Locale(config LineConfig) (string, error) {
	locale, err := sysinfo.Locale()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s [%s]", display.English.Tags().Name(locale), locale.String()), nil
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
