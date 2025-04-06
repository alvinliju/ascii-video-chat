package main

import (
	"fmt"
	"image"
	"net"
	"strings"

	"gocv.io/x/gocv"
)

const asciiChars = "@%#*+=-:. "
const addr = "192.168.20.11:8080"

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	defer webcam.Close()
	img := gocv.NewMat()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to UDP:", err)
		return
	}
	defer conn.Close()

	for {
		webcam.Read(&img)
		goImg, _ := img.ToImage()
		asciiArt := toASCIIArt(goImg, asciiChars)
		fmt.Print("\033[H\033[2J") // Clear terminal
		fmt.Println(asciiArt)

		_, err := fmt.Fprint(conn, asciiArt)
		if err != nil {
			fmt.Println("Error sending TCP data:", err)
		}
	}

}

func toASCIIArt(img image.Image, charMap string) string {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var sb strings.Builder

	for y := 0; y < height; y = y + 20 {
		for x := 0; x < width; x = x + 15 {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
			chars := float64(len(charMap) - 1)
			normalizedGray := float64(gray) / 255.0
			charIndex := int(normalizedGray * chars)
			sb.WriteByte(charMap[charIndex])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
