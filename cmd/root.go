package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nfetch/internal/color"
	. "nfetch/pkg"
	"nfetch/pkg/lines"
	"nfetch/pkg/logo"
	"nfetch/pkg/sysinfo"
)

var cfg struct {
	File string
	All  bool
}

//go:embed config.yaml
var defaultConfig []byte

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "nfetch",
	Short:   "A truly cross-platform alternative to neofetch",
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		// prefetch host info
		sysinfo.HostInfo()

		HideCursor()
		DisableLineWrap()
		defer ShowCursor()
		defer EnableLineWrap()

		logoString, logoColors := logo.GetLogo()

		if !color.NoColor {
			color.SetColors(logoColors...)
		}

		logoWidth, logoHeight := logo.PrintLogo(logoString)

		offset := CursorRight(logoWidth + 3)
		CursorUp(logoHeight)

		// render lines
		var showLines []interface{}
		if cfg.All {
			showLines = lines.AllLines
		} else {
			var ok bool
			showLines, ok = viper.Get("lines").([]interface{})
			if !ok {
				showLines = lines.AllLines
			}
		}
		lines.RenderLines(offset, showLines)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfg.File, "config", "", "config file (default ~/.config/nfetch/config.yaml)")
	rootCmd.Flags().BoolVarP(&cfg.All, "all", "a", false, "show all info lines")
	rootCmd.Flags().Bool("timing", false, "show time taken for each info line")
	rootCmd.Flags().StringP("logo", "l", "", "override platform specific logo")
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	viper.BindPFlag("timing", rootCmd.Flags().Lookup("timing"))
	viper.BindPFlag("logo", rootCmd.Flags().Lookup("logo"))
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.config/nfetch")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
			if err != nil {
				panic(fmt.Errorf("Failed reading default config: %s \n", err))
			}
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

	}
}
