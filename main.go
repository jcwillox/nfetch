package main

import (
	"nfetch/cmd"
	"nfetch/pkg/utils"
	"os"
)

func main() {
	utils.StripWSLPath()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
