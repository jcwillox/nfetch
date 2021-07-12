// +build windows

package sysinfo

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/miekg/dns"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/registry"
	"net"
	"strings"
	"sync"
)

const wqlBaseboard = "SELECT Manufacturer, Product FROM Win32_BaseBoard"
const wqlModel = "SELECT Manufacturer, Model FROM Win32_ComputerSystem"
const wqlCPU = "SELECT Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed FROM Win32_Processor"
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

type Win32CPU struct {
	Name                      *string
	NumberOfCores             *uint64
	NumberOfLogicalProcessors *uint64
	MaxClockSpeed             *uint64
}

type Win32GPU struct {
	Name *string
}

type Win32PageFile struct {
	CurrentUsage      *uint64
	AllocatedBaseSize *uint64
}

var connection *wmi.SWbemServicesConnection
var wmiLock sync.Mutex

func WmiSharedConnection() *wmi.SWbemServicesConnection {
	wmiLock.Lock()
	defer wmiLock.Unlock()
	if connection != nil {
		return connection
	}
	if wmi.DefaultClient.SWbemServicesClient == nil {
		sWbemServices, err := wmi.NewSWbemServices()
		if err != nil {
			panic(fmt.Errorf("Failed creating SWbemServicesClient: %s \n", err))
		}
		// this is not really threadsafe
		wmi.DefaultClient.SWbemServicesClient = sWbemServices
	}
	c, err := wmi.DefaultClient.SWbemServicesClient.ConnectServer()
	if err != nil {
		panic(fmt.Errorf("Failed creating SWbemServicesConnection: %s \n", err))
	}
	connection = c
	return c
}

func Distro() string {
	platform, _, _, err := host.PlatformInformation()
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(platform, "Microsoft ")
}

func Motherboard() (MotherboardInfo, error) {
	var win32BaseboardDescriptions []Win32Baseboard
	conn := WmiSharedConnection()
	if err := conn.Query(wqlBaseboard, &win32BaseboardDescriptions); err != nil {
		return MotherboardInfo{}, err
	}
	if len(win32BaseboardDescriptions) > 0 {
		return MotherboardInfo{
			Manufacturer: *win32BaseboardDescriptions[0].Manufacturer,
			Product:      *win32BaseboardDescriptions[0].Product,
		}, nil
	}
	return MotherboardInfo{}, nil
}

func Model() (ModelInfo, error) {
	var win32ModelDescriptions []Win32Model
	conn := WmiSharedConnection()
	if err := conn.Query(wqlModel, &win32ModelDescriptions); err != nil {
		return ModelInfo{}, err
	}
	if len(win32ModelDescriptions) > 0 {
		return ModelInfo{
			Manufacturer: *win32ModelDescriptions[0].Manufacturer,
			Model:        *win32ModelDescriptions[0].Model,
		}, nil
	}
	return ModelInfo{}, nil
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

func cpuInfo() (CPUInfo, error) {
	var win32CPUDescriptions []Win32CPU
	conn := WmiSharedConnection()
	if err := conn.Query(wqlCPU, &win32CPUDescriptions); err != nil {
		return CPUInfo{}, err
	}
	if len(win32CPUDescriptions) > 0 {
		return CPUInfo{
			Model:   *win32CPUDescriptions[0].Name,
			Cores:   int32(*win32CPUDescriptions[0].NumberOfCores),
			Threads: int32(*win32CPUDescriptions[0].NumberOfLogicalProcessors),
			Mhz:     float64(*win32CPUDescriptions[0].MaxClockSpeed),
		}, nil
	}
	return CPUInfo{}, nil
}

func GPU() (string, error) {
	var win32GPUDescriptions []Win32GPU
	conn := WmiSharedConnection()
	if err := conn.Query(wqlGPU, &win32GPUDescriptions); err != nil {
		return "", err
	}
	if len(win32GPUDescriptions) > 0 {
		return *win32GPUDescriptions[0].Name, nil
	}
	return "", nil
}

func Swap() (mem.SwapMemoryStat, error) {
	var win32PageFileDescriptions []Win32PageFile
	conn := WmiSharedConnection()
	if err := conn.Query(wqlPageFile, &win32PageFileDescriptions); err != nil {
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
