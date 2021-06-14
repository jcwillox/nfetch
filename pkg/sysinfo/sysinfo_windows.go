// +build windows

package sysinfo

import (
	"github.com/StackExchange/wmi"
	"github.com/miekg/dns"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/registry"
	"net"
	"strings"
)

const wqlBaseboard = "SELECT Manufacturer, Product FROM Win32_BaseBoard"
const wqlModel = "SELECT Manufacturer, Model FROM Win32_ComputerSystem"
const wqlGPU = "SELECT Name FROM Win32_VideoController"
const wqlPageFile = "SELECT CurrentUsage, AllocatedBaseSize FROM Win32_PageFileUsage"

type Win32Baseboard struct {
	Manufacturer *string
	Product      *string
}

type Win32Model struct {
	Manufacturer *string
	Model        *string
}

type Win32GPU struct {
	Name *string
}

type Win32PageFile struct {
	CurrentUsage      *uint64
	AllocatedBaseSize *uint64
}

func init() {
	sWbemServices, err := wmi.InitializeSWbemServices(wmi.DefaultClient)
	if err != nil {
		panic(err)
	}
	wmi.DefaultClient.SWbemServicesClient = sWbemServices
}

func CPUName() (string, error) {
	return ReadRegistryString(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, "ProcessorNameString")
}

func Distro() string {
	platform, _, _, err := host.PlatformInformation()
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(platform, "Microsoft ")
}

func Motherboard() (Win32Baseboard, error) {
	var win32BaseboardDescriptions []Win32Baseboard
	if err := wmi.Query(wqlBaseboard, &win32BaseboardDescriptions); err != nil {
		return Win32Baseboard{}, err
	}
	if len(win32BaseboardDescriptions) > 0 {
		return win32BaseboardDescriptions[0], nil
	}
	return Win32Baseboard{}, nil
}

func Model() (Win32Model, error) {
	var win32ModelDescriptions []Win32Model
	if err := wmi.Query(wqlModel, &win32ModelDescriptions); err != nil {
		return Win32Model{}, err
	}
	if len(win32ModelDescriptions) > 0 {
		return win32ModelDescriptions[0], nil
	}
	return Win32Model{}, nil
}

func Theme() (uint64, uint64, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		return 0, 0, err
	}
	defer key.Close()

	sysTheme, _, err := key.GetIntegerValue("SystemUsesLightTheme")
	if err != nil {
		return 0, 0, err
	}

	appTheme, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return 0, 0, err
	}

	return sysTheme, appTheme, nil
}

func GPU() (Win32GPU, error) {
	var win32GPUDescriptions []Win32GPU
	if err := wmi.Query(wqlGPU, &win32GPUDescriptions); err != nil {
		return Win32GPU{}, err
	}
	if len(win32GPUDescriptions) > 0 {
		return win32GPUDescriptions[0], nil
	}
	return Win32GPU{}, nil
}

func Swap() (mem.SwapMemoryStat, error) {
	var win32PageFileDescriptions []Win32PageFile
	if err := wmi.Query(wqlPageFile, &win32PageFileDescriptions); err != nil {
		return mem.SwapMemoryStat{}, err
	}
	if len(win32PageFileDescriptions) > 0 {
		stat := mem.SwapMemoryStat{
			Total: *win32PageFileDescriptions[0].AllocatedBaseSize * 1048576,
			Used:  *win32PageFileDescriptions[0].CurrentUsage * 1048576,
		}
		stat.Free = stat.Total - stat.Used
		stat.UsedPercent = float64(stat.Used) / float64(stat.Total) * 100
		return stat, nil
	}

	return mem.SwapMemoryStat{}, nil
}

func publicIP() (net.IP, error) {
	target := "myip.opendns.com"
	server := "resolver1.opendns.com:53"

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	r, _, err := c.Exchange(&m, server)
	if err != nil {
		return nil, err
	}

	for _, answer := range r.Answer {
		record := answer.(*dns.A)
		return record.A, nil
	}

	return nil, nil
}
