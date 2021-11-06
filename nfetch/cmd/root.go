package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jcwillox/emerald"
	"github.com/jcwillox/nfetch/internal/color"
	"github.com/jcwillox/nfetch/pkg/lines"
	"github.com/jcwillox/nfetch/pkg/logo"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
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
		if emerald.ColorEnabled {
			emerald.HideCursor()
			emerald.DisableLineWrap()
			defer emerald.ShowCursor()
			defer emerald.EnableLineWrap()
		}

		logoString, colors := logo.GetLogo()
		if viper.IsSet("colors") {
			colors = viper.GetStringSlice("colors")
		}
		color.SetColors(colors)

		if viper.GetBool("no_image") {
			logoString = ""
		}

		renderedLogo, logoWidth, logoHeight := logo.RenderLogo(logoString, colors)

		if logoWidth > 0 {
			logoWidth += viper.GetInt("gap")
		}

		if emerald.ColorEnabled && logoString != "" {
			paddingAmt := viper.GetInt("padding")
			padding := strings.Repeat(" ", paddingAmt)
			logoWidth += paddingAmt
			for _, line := range renderedLogo {
				emerald.Print(padding, line)
			}
			emerald.CursorUp(logoHeight)
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

		if emerald.ColorEnabled {
			writtenLines := lines.RenderLines(logoWidth, showLines, nil)

			// move cursor back to the bottom
			diff := logoHeight - writtenLines
			if diff > 0 {
				emerald.CursorDown(diff)
			}
		} else {
			lines.RenderLines(logoWidth, showLines, renderedLogo)
		}

		// print a final blank line
		emerald.Println()
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
	rootCmd.Flags().Bool("show-none", false, "show info lines that have no information")
	rootCmd.Flags().Bool("no-image", false, "hide image or logo")
	rootCmd.Flags().Int("padding", 1, "space before the image")
	rootCmd.Flags().IntP("gap", "g", 3, "gap between the image and text")
	rootCmd.Flags().StringSlice("colors", nil, "override the distro colors this applies to the text and logo colors")
	rootCmd.Flags().StringSlice("logo-colors", nil, "override the logo colors specifically")

	viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	viper.BindPFlag("timing", rootCmd.Flags().Lookup("timing"))
	viper.BindPFlag("logo", rootCmd.Flags().Lookup("logo"))
	viper.BindPFlag("show_none", rootCmd.Flags().Lookup("show-none"))
	viper.BindPFlag("no_image", rootCmd.Flags().Lookup("no-image"))
	viper.BindPFlag("padding", rootCmd.Flags().Lookup("padding"))
	viper.BindPFlag("gap", rootCmd.Flags().Lookup("gap"))
	viper.BindPFlag("colors", rootCmd.Flags().Lookup("colors"))
	viper.BindPFlag("logo_colors", rootCmd.Flags().Lookup("logo-colors"))

	rootCmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
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
		emerald.AutoSetColorState()
	case "always":
		emerald.SetColorState(true)
	case "never":
		emerald.SetColorState(false)
	}
}
