package main

import (
	"nfetch/cmd"
	"nfetch/pkg"
	"os"
)

func main() {
	pkg.StripWSLPath()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
