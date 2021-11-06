package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jcwillox/emerald"
	"github.com/jcwillox/nfetch/pkg/lines"
	"github.com/jcwillox/nfetch/pkg/sysinfo"
	"github.com/k0kubun/pp/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strings"
)

var pprint *pp.PrettyPrinter
var rawArgs = []string{
	lines.LineTitle,
	"distro",
	lines.LineOS,
	lines.LineHost,
	"model",
	lines.LineKernel,
	lines.LineMotherboard,
	lines.LineBios,
	lines.LineUptime,
	lines.LinePkgs,
	lines.LineShell,
	lines.LineResolution,
	lines.LineTerminal,
	lines.LineTheme,
	lines.LineCPU,
	"cores",
	"temps",
	lines.LineGPU,
	lines.LineUsage,
	lines.LineMemory,
	lines.LineSwap,
	lines.LineDisk,
	lines.LineBattery,
	lines.LineLocale,
	lines.LineWeather,
	lines.LineLocalIP,
	lines.LinePublicIP,
	"wallpaper",
	lines.LineColorbar,
}

var rawCmd = &cobra.Command{
	Use:       "raw <info>",
	Hidden:    false,
	ValidArgs: rawArgs,
	Args:      cobra.ExactValidArgs(1),
	Annotations: map[string]string{
		"Info": strings.Join(rawArgs, ", "),
	},
	Run: func(cmd *cobra.Command, args []string) {
		pprint = pp.New()
		pprint.SetDecimalUint(true)

		switch args[0] {
		case lines.LineTitle:
			printInfo(lines.Title())
		case lines.LineOS:
			printInfo(lines.OS(lines.LineConfig{}))
		case "distro":
			printInfo(sysinfo.Distro())
		case lines.LineHost:
			printInfo(host.Info())
		case "model":
			printInfo(sysinfo.Model())
		case lines.LineKernel:
			printInfo(host.KernelVersion())
		case lines.LineMotherboard:
			printInfo(sysinfo.Motherboard())
		case lines.LineBios:
			printInfo(sysinfo.Bios())
		case lines.LineUptime:
			printInfo(host.Uptime())
		case lines.LinePkgs:
			printInfo(sysinfo.Pkgs(nil))
		case lines.LineShell:
			printInfo(sysinfo.Shell())
		case lines.LineResolution:
			printInfo(sysinfo.Resolution())
		case lines.LineTerminal:
			printInfo(sysinfo.Terminal())
		case lines.LineTheme:
			printInfo(lines.Theme(lines.LineConfig{}))
		case lines.LineCPU:
			printInfo(sysinfo.CPU())
		case "cores":
			printInfo(cpu.Info())
		case "temps":
			printInfo(host.SensorsTemperatures())
		case lines.LineGPU:
			printInfo(sysinfo.GPU())
		case lines.LineUsage:
			printInfo(sysinfo.Usage())
		case lines.LineMemory:
			printInfo(sysinfo.Memory())
		case lines.LineSwap:
			printInfo(sysinfo.Swap())
		case lines.LineDisk:
			printInfo(sysinfo.Disk())
		case lines.LineBattery:
			printInfo(sysinfo.Battery())
		case lines.LineLocale:
			printInfo(lines.Locale(lines.LineConfig{}))
			printInfo(sysinfo.Locale())
		case lines.LineWeather:
			printInfo(sysinfo.Weather())
		case lines.LineLocalIP:
			printInfo(sysinfo.LocalIP())
		case lines.LinePublicIP:
			printInfo(sysinfo.PublicIP())
		case "wallpaper":
			printInfo(sysinfo.Wallpaper())
		case lines.LineColorbar:
			for _, line := range lines.Colorbar() {
				emerald.Println(line)
			}
		}
	},
}

func printInfo(a interface{}, err ...error) {
	if result, ok := a.(string); ok {
		fmt.Println(result)
	} else if result, ok := a.(uint64); ok {
		fmt.Println(result)
	} else if result, ok := a.(net.IP); ok {
		fmt.Println(result)
	} else if emerald.ColorEnabled {
		if len(err) > 0 && err[0] != nil {
			pprint.Println(err[0])
		} else {
			pprint.Println(a)
		}
	} else {
		if len(err) > 0 && err[0] != nil {
			fmt.Fprintln(os.Stderr, "error:", err[0])
		} else {
			result, _ := json.MarshalIndent(a, "", "  ")
			os.Stdout.Write(result)
		}
	}
}

func init() {
	rootCmd.AddCommand(rawCmd)
	rawCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

Info: 
  {{.Annotations.Info}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`)
}
