package main

import (
	"github.com/jcwillox/nfetch/nfetch/cmd"
	"github.com/jcwillox/nfetch/pkg/utils"
	"os"
)

func main() {
	utils.StripWSLPath()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
