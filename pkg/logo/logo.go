package logo

import (
	"bytes"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/jcwillox/emerald"
	"github.com/jcwillox/nfetch/internal/color"
	"github.com/jcwillox/nfetch/pkg/sysinfo"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"text/template"
)

// GetLogo returns the logo and default colors for the current distro
func GetLogo() (string, []string) {
	logo := viper.GetString("logo")
	if logo == "" {
		logo = sysinfo.Distro()
	}
	// case-insensitive
	logo = strings.ToLower(logo)
	return GetDistroLogo(logo)
}

func RenderLogo(rawLogoTemplate string, colors []string) (logo []string, logoWidth int, logoHeight int) {
	if rawLogoTemplate == "" {
		return []string{}, 0, 0
	}

	logoTemplate, err := template.New("logo").Parse(rawLogoTemplate)
	if err != nil {
		panic(fmt.Errorf("Unable to parse logo template: %s \n", err))
	}

	colorsLen := len(colors)
	if viper.IsSet("logo_colors") {
		colors = viper.GetStringSlice("logo_colors")
		if len(colors) > colorsLen {
			colorsLen = len(colors)
		}
	}

	colorData := map[string]string{}
	for i := 0; i < colorsLen; i++ {
		colorData["C"+strconv.Itoa(i+1)] = ""
	}

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
	for i, c := range colors {
		colorData["C"+strconv.Itoa(i+1)] = color.ColorizerCode(c, true)
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
