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

	type DiskResult struct {
		Titles  []string
		Content []string
	}

	lineResults := make(map[string]Result)

	results := make(chan Result)
	diskResult := make(chan DiskResult)

	showTiming := viper.GetBool("timing")
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
				titles, content, err := Disk(config)
				if err != nil {
					titles = []string{"Disk"}
					content = []string{color.ErrorMsg}
				}
				diskResult <- DiskResult{
					Titles:  titles,
					Content: content,
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
				content = fmt.Sprint(content, color.AU.Yellow(" [took "), color.AU.Red(time.Since(start)).String(), color.AU.Yellow("]").String())
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
			diskResult := <-diskResult
			for i, title := range diskResult.Titles {
				printLine(aurora.Colorize(title, color.Colors.C1), ": ", diskResult.Content[i])
			}
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
