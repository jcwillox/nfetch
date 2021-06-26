package lines

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/viper"
	"nfetch/internal/color"
	. "nfetch/pkg"
	"nfetch/pkg/ioutils"
	"strings"
	"time"
)

func RenderLines(offset int, lines []interface{}, logo []string) int {
	lineStrSlice := make([]string, 0, len(lines))

	type Result struct {
		Line    string
		Title   string
		Content string
	}

	lineResults := make(map[string]Result)
	results := make(chan Result)
	showTiming := viper.GetBool("timing")

	var diskTitles []string
	var diskContent []string

	start := time.Now()

	// start goroutines to render lines
	for _, entry := range lines {
		line, config := GetLineConfig(entry)
		lineStrSlice = append(lineStrSlice, line)

		switch line {
		case LineTitle, LineDashes, LineColorbar, LineBlank:
			continue
		case LineDisk:
			go func() {
				var err error
				diskTitles, diskContent, err = Disk(config)
				if err != nil {
					diskTitles = []string{"Disk"}
					diskContent = []string{color.ErrorMsg}
				}
			}()
			continue
		}

		// get title
		title := config.GetString("title")
		if title == "" {
			title = DefaultTitleMap[line]
			if title == "" {
				title = line
			}
		}

		contentFunc, ok := FuncMap[line]

		go func() {
			if !ok {
				results <- Result{
					Line:    line,
					Title:   title,
					Content: color.Error("(does not exist)").String(),
				}
				return
			}

			start := time.Now()
			content, err := contentFunc(config)
			if err != nil {
				content = color.ErrorMsg
			}

			if showTiming {
				content = fmt.Sprint(content, aurora.Yellow(" [took "), aurora.Red(time.Since(start)).String(), aurora.Yellow("]").String())
			}

			results <- Result{
				Line:    line,
				Title:   title,
				Content: content,
			}
		}()
	}

	writtenLines := 0
	var printLine func(a ...interface{})

	if logo == nil {
		// -1 to account for difference between counting columns and characters
		prefix := CursorRight(offset + 1)
		printLine = func(a ...interface{}) {
			ioutils.Print(prefix)
			ioutils.Print(a...)
			ioutils.Println()
			writtenLines += 1
		}
	} else {
		printLine = func(a ...interface{}) {
			if len(logo) > writtenLines {
				ioutils.Print(logo[writtenLines], strings.Repeat(" ", offset-len(logo[writtenLines])))
			} else {
				ioutils.Print(strings.Repeat(" ", offset))
			}
			ioutils.Print(a...)
			ioutils.Println()
			writtenLines += 1
		}
	}

	// receive results
	for i := 0; i < len(lineStrSlice); i++ {
		line := lineStrSlice[i]

		switch line {
		case LineTitle:
			printLine(Title())
		case LineDashes:
			printLine(Dashes())
		case LineBlank:
			printLine()
		case LineColorbar:
			if color.NoColor {
				break
			}
			for _, s := range Colorbar() {
				printLine(s)
			}
		case LineDisk:
			if diskTitles != nil && diskContent != nil {
				for i := range diskTitles {
					printLine(aurora.Colorize(diskTitles[i], color.Colors.C1), ": ", diskContent[i])
				}
				break
			}
			// we haven't go disks yet so process a result
			fallthrough
		default:
			// try get from map
			if res, present := lineResults[line]; present {
				printLine(aurora.Colorize(res.Title, color.Colors.C1), ": ", res.Content)
				break
			}

			// read a result
			result := <-results
			// check if current result
			if result.Line == line {
				printLine(aurora.Colorize(result.Title, color.Colors.C1), ": ", result.Content)
			} else {
				i--
			}

			// otherwise get next result
			lineResults[result.Line] = result

		}
	}

	if showTiming {
		printLine("total time: ", time.Since(start))
		writtenLines++
	}

	// print remaining lines when in no-color mode
	for i := writtenLines; i < len(logo); i++ {
		printLine()
	}

	return writtenLines
}
