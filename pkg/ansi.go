package pkg

import (
	"nfetch/pkg/ioutils"
	"strconv"
)

func CursorUp(lines int) {
	ioutils.Printf("\x1b[%dA", lines)
}

func CursorDown(lines int) {
	ioutils.Printf("\x1b[%dB", lines)
}

func CursorRight(columns int) string {
	return "\x1b[" + strconv.Itoa(columns) + "G"
}

func ShowCursor() {
	ioutils.Print("\x1b[?25h")
}

func HideCursor() {
	ioutils.Print("\x1b[?25l")
}

func EnableLineWrap() {
	ioutils.Print("\x1b[?7h")
}

func DisableLineWrap() {
	ioutils.Print("\x1b[?7l")
}
