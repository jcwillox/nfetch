module nfetch

go 1.16

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46
	github.com/Xuanwo/go-locale v1.0.0
	github.com/anthonynsimon/bild v0.13.0
	github.com/bi-zone/wmi v1.1.4 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/k0kubun/pp/v3 v3.0.7
	github.com/kbinani/screenshot v0.0.0-20210326165202-b96eb3309bb0
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/miekg/dns v1.1.42
	github.com/reujab/wallpaper v0.0.0-20201124162023-c3898ec30d2c
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	golang.org/x/sys v0.0.0-20210611083646-a4fc73990273
	golang.org/x/text v0.3.6
)

replace github.com/k0kubun/pp/v3 v3.0.7 => github.com/k0kubun/pp/v3 v3.0.8-0.20210415165650-b87d88f85b84

replace (
	github.com/StackExchange/wmi => github.com/bi-zone/wmi v1.1.4
	github.com/bi-zone/wmi v1.1.4 => github.com/jeffreystoke/wmi v1.1.5-0.20201112194144-6556453f893c
)
