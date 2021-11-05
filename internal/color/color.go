package color

import (
	"github.com/jcwillox/emerald"
	"strconv"
)

var (
	Error    func(string) string
	ErrorMsg string
)

var (
	Title     func(string) string
	At        func(string) string
	Dashes    func(string) string
	Subtitle  func(string) string
	Separator func(string) string
	Info      func(string) string
)

// SetFromLogoColors sets the text colors based on the logo colors
func SetFromLogoColors(colors []int) {
	if emerald.ColorEnabled {
		if len(colors) > 0 {
			Title = emerald.ColorFunc(strconv.Itoa(colors[0]) + "+b")
		}
		if len(colors) < 2 || colors[1] == 7 || colors[1] == 15 {
			Subtitle = Title
		} else {
			Subtitle = emerald.ColorFunc(strconv.Itoa(colors[1]) + "+b")
		}
		At = emerald.ColorFunc("white")
		Dashes = At
		Separator = At
		Info = At
		Error = emerald.ColorFunc("red")
		ErrorMsg = Error("(error)")
	} else {
		Title = func(s string) string { return s }
		At = Title
		Dashes = Title
		Subtitle = Title
		Separator = Title
		Info = Title
		Error = Title
		ErrorMsg = "(error)"
	}
}
