package logo

import (
	"bytes"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/jcwillox/emerald"
	"github.com/spf13/viper"
	"nfetch/pkg/sysinfo"
	"strings"
	"text/template"
)

func GetLogo() (string, []int) {
	logo := viper.GetString("logo")
	if logo == "" {
		logo = sysinfo.Distro()
	}
	// case-insensitive
	logo = strings.ToLower(logo)
	return getLogo(logo)
}

func RenderLogo(rawLogoTemplate string, colors []int) (logo []string, logoWidth int, logoHeight int) {
	if rawLogoTemplate == "" {
		return []string{}, 0, 0
	}

	logoTemplate, err := template.New("logo").Parse(rawLogoTemplate)
	if err != nil {
		panic(fmt.Errorf("Unable to parse logo template: %s \n", err))
	}

	colorData := struct {
		C1 string
		C2 string
		C3 string
		C4 string
	}{}

	// parse without colours
	var buff bytes.Buffer
	err = logoTemplate.Execute(&buff, colorData)
	if err != nil {
		panic(fmt.Errorf("Unable to render plain logo: %s \n", err))
	}

	// TODO: count newlines to get size?
	lines := make([]string, 0, 10)

	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		logoHeight++
		line = strings.TrimRight(line, "\r\n")
		length := len(line)
		if length > logoWidth {
			logoWidth = length
		}

		if !emerald.ColorEnabled {
			lines = append(lines, line)
		}
	}

	if !emerald.ColorEnabled {
		return lines, logoWidth, logoHeight
	}

	buff.Reset()

	// parse with colours
	if len(colors) > 0 {
		colorData.C1 = emerald.ColorIndexFg(colors[0])
	}
	if len(colors) > 1 {
		colorData.C2 = emerald.ColorIndexFg(colors[1])
	}
	if len(colors) > 2 {
		colorData.C3 = emerald.ColorIndexFg(colors[2])
	}
	if len(colors) > 3 {
		colorData.C4 = emerald.ColorIndexFg(colors[3])
	}

	err = logoTemplate.Execute(&buff, colorData)
	if err != nil {
		panic(fmt.Errorf("Unable to render logo: %s \n", err))
	}

	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}

	return lines, logoWidth, logoHeight
}

func PrintAsciiImage(path string) {
	img, err := imgio.Open(path)
	if err != nil {
		emerald.Println(err)
		return
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	emerald.Println(width, height)

	newWidth := 35
	ratio := float32(height) / float32(width)
	newHeight := int(ratio*float32(newWidth) + 0.5)

	emerald.Println(newWidth, newHeight, ratio)

	resized := transform.Resize(img, newWidth, newHeight, transform.Linear)

	block := "\u2580"

	for y := 0; y < newHeight; y += 2 {
		for x := 0; x < newWidth; x++ {
			backVT := ""
			if y < newHeight-1 {
				backRGBA := resized.RGBAAt(x, y+1)
				backVT = fmt.Sprintf("\x1b[48;2;%d;%d;%dm", backRGBA.R, backRGBA.G, backRGBA.B)
			}
			foreRGBA := resized.RGBAAt(x, y)
			foreVT := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", foreRGBA.R, foreRGBA.G, foreRGBA.B)
			emerald.Print(foreVT + backVT + block + "\x1b[0m")
		}
		emerald.Println()
	}
}
