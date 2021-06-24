package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nfetch/internal/color"
	. "nfetch/pkg"
	"nfetch/pkg/ioutils"
	"nfetch/pkg/lines"
	"nfetch/pkg/logo"
	"os"
	"path/filepath"
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
		if !color.NoColor {
			HideCursor()
			DisableLineWrap()
			defer ShowCursor()
			defer EnableLineWrap()
		}

		logoString, logoColors := logo.GetLogo()

		if !color.NoColor {
			color.SetColors(logoColors...)
		} else {
			// TODO: allow print logo in no color mode
			logoString = ""
		}

		offset := ""
		logoWidth, logoHeight := logo.PrintLogo(logoString)

		if !color.NoColor && logoString != "" {
			offset = CursorRight(logoWidth + 3)
			CursorUp(logoHeight)
		}

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
		writtenLines := lines.RenderLines(offset, showLines)

		// move cursor back to the bottom
		if !color.NoColor {
			diff := logoHeight - writtenLines
			if diff > 0 {
				CursorDown(diff)
			}
		}
		// print a final blank line
		ioutils.Println()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfg.File, "config", "", "config file (default \"~/.config/nfetch/config.yaml\")")
	rootCmd.PersistentFlags().String("color", "auto", "when to use colors (always, auto, never)")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for nfetch")
	rootCmd.Flags().BoolVarP(&cfg.All, "all", "a", false, "show all info lines")
	rootCmd.Flags().Bool("timing", false, "show time taken for each info line")
	rootCmd.Flags().StringP("logo", "l", "", "override platform specific logo")
	rootCmd.Flags().BoolP("version", "v", false, "version for nfetch")
	viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	viper.BindPFlag("timing", rootCmd.Flags().Lookup("timing"))
	viper.BindPFlag("logo", rootCmd.Flags().Lookup("logo"))

	rootCmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func initConfig() {
	var configDir string
	var configPath string

	if cfg.File != "" {
		viper.SetConfigFile(cfg.File)
	} else {
		viper.AddConfigPath(".")

		homeDir, err := homedir.Dir()
		if err == nil {
			configPath = filepath.Join(homeDir, ".config/nfetch/config.yaml")
			configDir = filepath.Dir(configPath)
			viper.AddConfigPath(configDir)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// using built-in default config
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
			if err != nil {
				panic(fmt.Errorf("failed reading default config: %s\n", err))
			}

			// write default config if none exists
			if _, err := os.Stat(configPath); err != nil {
				os.MkdirAll(configDir, os.FileMode(0755))
				fmt.Fprintln(os.Stderr, "writing default config to ~/.config/nfetch/config.yaml")
				err = os.WriteFile(configPath, defaultConfig, os.FileMode(0644))
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to write default config", err)
				}
			}
		} else {
			panic(fmt.Errorf("fatal error reading config file: %s\n", err))
		}
	}

	// handle global flags
	switch viper.GetString("color") {
	case "auto":
		color.InitColors(ioutils.IsTerminal)
	case "always":
		color.InitColors(true)
	case "never":
		color.InitColors(false)
	}
}
