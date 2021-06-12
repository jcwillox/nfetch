package pkg

import (
	"fmt"
	"strconv"
)

func CursorUp(lines int) {
	fmt.Printf("\x1b[%dA", lines)
}

func CursorDown(lines int) {
	fmt.Printf("\x1b[%dB", lines)
}

func CursorRight(columns int) string {
	return "\x1b[" + strconv.Itoa(columns) + "G"
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func EnableLineWrap() {
	fmt.Print("\x1b[?7h")
}

func DisableLineWrap() {
	fmt.Print("\x1b[?7l")
}
