package color

import (
	"github.com/mgutz/ansi"
	"strconv"
)

var (
	Error    func(string) string
	ErrorMsg string
	NoColor  bool
)

var (
	Title     func(string) string
	At        func(string) string
	Dashes    func(string) string
	Subtitle  func(string) string
	Separator func(string) string
	Info      func(string) string
)

func InitColors(enableColors bool) {
	NoColor = !enableColors
	if NoColor {
		ansi.DisableColors(true)
	}
	Error = ansi.ColorFunc("red")
	ErrorMsg = Error("(error)")
}

// SetFromLogoColors sets the text colors based on the logo colors
func SetFromLogoColors(colors []int) {
	if NoColor {
		Title = func(s string) string { return s }
		At = Title
		Dashes = Title
		Subtitle = Title
		Separator = Title
		Info = Title
	} else {
		if len(colors) > 0 {
			Title = ansi.ColorFunc(strconv.Itoa(colors[0]) + "+b")
		}
		if len(colors) < 2 || colors[1] == 7 || colors[1] == 15 {
			Subtitle = Title
		} else {
			Subtitle = ansi.ColorFunc(strconv.Itoa(colors[1]) + "+b")
		}
		At = ansi.ColorFunc("white")
		Dashes = At
		Separator = At
		Info = At
	}
}
