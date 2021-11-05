package color

import (
	"github.com/jcwillox/emerald"
	"strings"
)

var (
	Title     emerald.Colorizer
	At        emerald.Colorizer
	Dashes    emerald.Colorizer
	Subtitle  emerald.Colorizer
	Separator emerald.Colorizer
	Info      emerald.Colorizer
	Error     emerald.Colorizer
	ErrorMsg  string
)

func ColorizerCode(style string, bold bool) string {
	colorCode := emerald.ColorCode(style)
	if bold && !strings.HasPrefix(colorCode, "\x1b[0;1") {
		return strings.Replace(colorCode, "[0;", "[0;1;", 1)
	}
	return colorCode
}

func ColorizerFunc(style string, bold bool) emerald.Colorizer {
	if !bold || style == "" {
		return emerald.ColorFunc(style)
	}
	colorCode := ColorizerCode(style, bold)
	return func(s string) string {
		if !emerald.ColorEnabled || s == "" {
			return s
		}
		return colorCode + s + emerald.Reset
	}
}

func SetColors(colors []string) {
	if emerald.ColorEnabled {
		if len(colors) > 0 {
			Title = ColorizerFunc(colors[0], true)
		}
		if len(colors) < 2 || colors[1] == "7" || colors[1] == "15" {
			Subtitle = Title
		} else {
			Subtitle = ColorizerFunc(colors[1], true)
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
