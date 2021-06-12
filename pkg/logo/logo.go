package logo

import (
	"bytes"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"nfetch/internal/color"
	"os"
	"text/template"
)

func PrintLogo(logo string) (int, int) {
	logoTemplate, err := template.New("logo").Parse(logo)
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

	logoHeight := 0
	logoWidth := 0

	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		logoHeight++
		length := len(line)
		if length > logoWidth {
			logoWidth = length
		}
	}

	// TODO: only parse template once when no colors

	// parse with colours
	colorData.C1 = "\x1b[" + color.Colors.C1.Nos(true) + "m"
	colorData.C2 = "\x1b[" + color.Colors.C2.Nos(true) + "m"
	colorData.C3 = "\x1b[" + color.Colors.C3.Nos(true) + "m"
	colorData.C4 = "\x1b[" + color.Colors.C4.Nos(true) + "m"

	err = logoTemplate.Execute(os.Stdout, colorData)
	if err != nil {
		panic(fmt.Errorf("Unable to render logo: %s \n", err))
	}

	// remove newline from width count
	return logoWidth - 1, logoHeight
}

func PrintAsciiImage(path string) {
	img, err := imgio.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	fmt.Println(width, height)

	newWidth := 35
	ratio := float32(height) / float32(width)
	newHeight := int(ratio*float32(newWidth) + 0.5)

	fmt.Println(newWidth, newHeight, ratio)

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
			fmt.Print(foreVT + backVT + block + "\x1b[0m")
		}
		fmt.Println()
	}
}
