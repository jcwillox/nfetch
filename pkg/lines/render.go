package lines

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/viper"
	"nfetch/internal/color"
	"time"
)

func RenderLines(offset string, lines []interface{}) int {
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

	// receive results
	writtenLines := 0
	for i := 0; i < len(lineStrSlice); i++ {
		line := lineStrSlice[i]

		switch line {
		case LineTitle:
			fmt.Print(offset, Title(), "\n")
			writtenLines += 1
		case LineDashes:
			fmt.Print(offset, Dashes(), "\n")
			writtenLines += 1
		case LineBlank:
			fmt.Println()
			writtenLines += 1
		case LineColorbar:
			for _, s := range Colorbar() {
				fmt.Print(offset, s, "\n")
				writtenLines += 1
			}

		case LineDisk:
			if diskTitles != nil && diskContent != nil {
				for i := range diskTitles {
					fmt.Print(offset, aurora.Colorize(diskTitles[i], color.Colors.C1), ": ")
					fmt.Println(diskContent[i])
					writtenLines += 1
				}
				break
			}
			// we haven't go disks yet so process a result
			fallthrough
		default:
			// try get from map
			if res, present := lineResults[line]; present {
				fmt.Print(offset, aurora.Colorize(res.Title, color.Colors.C1), ": ")
				fmt.Println(res.Content)
				writtenLines += 1
				break
			}

			// read a result
			result := <-results
			// check if current result
			if result.Line == line {
				fmt.Print(offset, aurora.Colorize(result.Title, color.Colors.C1), ": ")
				fmt.Println(result.Content)
				writtenLines += 1
			} else {
				i--
			}

			// otherwise get next result
			lineResults[result.Line] = result

		}
	}

	if showTiming {
		fmt.Print(offset, "total time: ", time.Since(start), "\n")
		writtenLines++
	}

	return writtenLines
}
