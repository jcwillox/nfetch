package pkg

import (
	"fmt"
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

func ColorIndexFg(index int) string {
	if index < 8 {
		return fmt.Sprintf("\x1b[%dm", index+30)
	}
	if index < 16 {
		return fmt.Sprintf("\x1b[%dm", index+82)
	}
	return fmt.Sprintf("\x1b[38;5;%dm", index)
}
