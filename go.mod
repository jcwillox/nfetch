module github.com/jcwillox/nfetch

go 1.17

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46
	github.com/Xuanwo/go-locale v1.0.0
	github.com/anthonynsimon/bild v0.13.0
	github.com/distatus/battery v0.10.0
	github.com/jcwillox/emerald v0.2.0
	github.com/k0kubun/pp/v3 v3.0.7
	github.com/kbinani/screenshot v0.0.0-20210326165202-b96eb3309bb0
	github.com/miekg/dns v1.1.42
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/reujab/wallpaper v0.0.0-20201124162023-c3898ec30d2c
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.0.0-20211103235746-7861aae1554b
	golang.org/x/text v0.3.6
)

require (
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/bi-zone/go-ole v1.2.5 // indirect
	github.com/bi-zone/wmi v1.1.4 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/gen2brain/shm v0.0.0-20200228170931-49f9650110c5 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/scjalliance/comshim v0.0.0-20190308082608-cf06d2532c4e // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v0.0.0-20201203080718-1454fab16a06 // indirect
)

replace github.com/k0kubun/pp/v3 v3.0.7 => github.com/k0kubun/pp/v3 v3.0.8-0.20210415165650-b87d88f85b84

replace (
	github.com/StackExchange/wmi => github.com/bi-zone/wmi v1.1.4
	github.com/bi-zone/wmi v1.1.4 => github.com/jeffreystoke/wmi v1.1.5-0.20201112194144-6556453f893c
)
