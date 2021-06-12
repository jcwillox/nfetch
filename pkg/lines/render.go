package lines

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"nfetch/internal/color"
	. "nfetch/pkg"
)

// TODO: simplify / breakup function

func RenderLines(offset string, lines []interface{}) {
	for _, entry := range lines {
		line, config := GetLineConfig(entry)
		switch line {
		case LineTitle:
			title, _ := Title()
			fmt.Println(offset + title)
		case LineDashes:
			dashes, _ := Dashes()
			fmt.Println(offset + dashes)
		case LineDisk:
			titles, contents, err := Disk(config)
			if err != nil {
				fmt.Print(offset, aurora.Colorize("Disk", color.Colors.C1))
				fmt.Println(": " + color.ErrorMsg)
			} else {
				for i := range titles {
					fmt.Print(offset, aurora.Colorize(titles[i], color.Colors.C1))
					fmt.Println(": " + contents[i])
				}
			}
		case LineBlank:
			fmt.Println()
		case LineColorbar:
			fmt.Println(offset + "\x1b[0;40m   \x1b[0;41m   \x1b[0;42m   \x1b[0;43m   \x1b[0;44m   \x1b[0;45m   \x1b[0;46m   \x1b[0;47m   \x1b[0m")
			fmt.Println(offset + "\x1b[0;100m   \x1b[0;101m   \x1b[0;102m   \x1b[0;103m   \x1b[0;104m   \x1b[0;105m   \x1b[0;106m   \x1b[0;107m   \x1b[0m")
		default:
			title := config.GetString("title")
			if title == "" {
				title = defaultTitleMap[line]
				if title == "" {
					title = line
				}
			}
			fmt.Print(offset, aurora.Colorize(title, color.Colors.C1))
			contentFunc, ok := funcMap[line]

			if !ok {
				fmt.Println(":", color.Error("(does not exist)"))
				continue
			}

			content, err := contentFunc(config)
			if err != nil {
				fmt.Println(":", color.ErrorMsg)
			} else {
				fmt.Println(": " + content)
			}
		}
	}
}
