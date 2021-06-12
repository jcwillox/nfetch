package color

import "github.com/logrusorgru/aurora/v3"

const (
	flagFg  aurora.Color = 1 << 14
	shiftFg aurora.Color = 16
)

var (
	AU       aurora.Aurora
	Error    func(arg interface{}) aurora.Value
	ErrorMsg string
	NoColor  bool
)

type ColorsType struct {
	C1 aurora.Color
	C2 aurora.Color
	C3 aurora.Color
	C4 aurora.Color
}

func InitColors(enableColors bool) {
	NoColor = !enableColors
	AU = aurora.NewAurora(enableColors)
	Error = AU.Red
	ErrorMsg = Error("(error)").String()
}

var Colors = ColorsType{}

func SetColors(colors ...int) {
	length := len(colors)

	if length > 0 {
		Colors.C1 = (aurora.Color(colors[0]) << shiftFg) | flagFg
	}
	if length > 1 {
		Colors.C2 = (aurora.Color(colors[1]) << shiftFg) | flagFg
	}
	if length > 2 {
		Colors.C3 = (aurora.Color(colors[2]) << shiftFg) | flagFg
	}
	if length > 3 {
		Colors.C4 = (aurora.Color(colors[3]) << shiftFg) | flagFg
	}
}
