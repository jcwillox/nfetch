package ioutils

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
)

var (
	Stdout     io.Writer
	IsTerminal bool
)

func init() {
	IsTerminal = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	Stdout = colorable.NewColorableStdout()
}

func Print(a ...interface{}) (n int) {
	n, _ = fmt.Fprint(Stdout, a...)
	return n
}

func Printf(format string, a ...interface{}) (n int) {
	n, _ = fmt.Fprintf(Stdout, format, a...)
	return n
}

func Println(a ...interface{}) (n int) {
	n, _ = fmt.Fprintln(Stdout, a...)
	return n
}
