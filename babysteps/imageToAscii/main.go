package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	file, err := os.Open("./richard.jpg")
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(img)
	fmt.Println(img.Bounds())
	fmt.Println(format)

	bounds := img.Bounds()
	asciiChars := "@%#*+=-:. " // Dark to light
	var result strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + 12 {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + 5 {
			r, g, b, _ := img.At(x, y).RGBA()

			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)

			chars := float64(len(asciiChars) - 1)
			greys := float64(gray) / 255
			charIndex := int(greys * chars)
			result.WriteByte(asciiChars[charIndex])

		}
		result.WriteByte('\n')
	}

	fmt.Print("\033[H\033[2J") // Clear terminal
	fmt.Print(result.String())
	os.WriteFile("richard.txt", []byte(result.String()), 06400)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "", 10)
	lines := strings.Split(result.String(), "\n")
	for _, line := range lines {
		pdf.CellFormat(0, 6, line, "", 1, "", false, 0, "")
	}
	pdf.OutputFileAndClose("richard.pdf")

}
