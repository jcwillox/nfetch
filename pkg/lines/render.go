package lines

import (
	"fmt"
	"github.com/jcwillox/emerald"
	"github.com/spf13/viper"
	"nfetch/internal/color"
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
	showNone := viper.GetBool("show_none")
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
					Content: color.Error("(does not exist)"),
				}
				return
			}

			start := time.Now()
			content, err := contentFunc(config)
			if err != nil {
				content = color.ErrorMsg
			}

			curShowNone := showNone
			if config.Has("show_none") {
				curShowNone = config.GetBool("show_none")
			}
			if content == "" && (curShowNone || showTiming) {
				content = "(none)"
			}

			if showTiming {
				content = fmt.Sprint(content, emerald.Yellow, " [took ", emerald.Red, time.Since(start), emerald.Yellow, "]", emerald.Reset)
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
		// +1 to account for difference between column and character count
		prefix := emerald.CursorRightVar(offset + 1)
		printLine = func(a ...interface{}) {
			emerald.Print(prefix)
			emerald.Print(a...)
			emerald.Println()
			writtenLines += 1
		}
	} else {
		// write logo progressively
		paddingAmt := 0
		if len(logo) > 0 {
			paddingAmt = viper.GetInt("padding")
		}
		padding := strings.Repeat(" ", paddingAmt)
		printLine = func(a ...interface{}) {
			if len(logo) > writtenLines {
				emerald.Print(padding, logo[writtenLines], strings.Repeat(" ", offset-len(logo[writtenLines])))
			} else {
				emerald.Print(padding, strings.Repeat(" ", offset))
			}
			emerald.Print(a...)
			emerald.Println()
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
			if emerald.ColorEnabled {
				for _, s := range Colorbar() {
					printLine(s)
				}
			}
		case LineDisk:
			diskResult := <-diskResult
			for i, title := range diskResult.Titles {
				printLine(color.Subtitle(title), color.Separator(": "), color.Info(diskResult.Content[i]))
			}
		default:
			// try get from map
			if res, present := lineResults[line]; present {
				if res.Content != "" {
					printLine(color.Subtitle(res.Title), color.Separator(": "), color.Info(res.Content))
				}
				break
			}

			// read a result
			result := <-results
			// check if current result
			if result.Line == line {
				if result.Content != "" {
					printLine(color.Subtitle(result.Title), color.Separator(": "), color.Info(result.Content))
				}
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
