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

var cfgFile string

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
		lines.RenderLines(offset, viper.Get("lines").([]interface{}))
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.config/nfetch/config.yaml)")
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
