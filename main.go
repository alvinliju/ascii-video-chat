package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"time"

	"github.com/pion/webrtc/v3"
	"gocv.io/x/gocv"
)

const asciiChars = "@%#*+=-:. "

func main() {
	//config
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}

	//peer coneection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Print("\033[H\033[2J") // Clear terminal
			fmt.Print(string(msg.Data))
		})
	})

	dataChannel, err := peerConnection.CreateDataChannel("video", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		fmt.Printf("ICE Candidate: %s\n", c.ToJSON().Candidate)
	})

	offer, _ := peerConnection.CreateOffer(nil)
	peerConnection.SetLocalDescription(offer)
	fmt.Printf("SDP Offer:\n%s\n", offer.SDP)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your SDP address")
	remoteSDP, _ := reader.ReadString('\n')

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		log.Println(err)
		return
	}

	peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  strings.TrimSpace(remoteSDP),
	})

	webcam, _ := gocv.OpenVideoCapture(0)
	img := gocv.NewMat()

	dataChannel.OnOpen(func() {
		fmt.Print("Data channel open")
		for {
			webcam.Read(&img)
			goImg, _ := img.ToImage()

			asciiArt := toASCIIArt(goImg, asciiChars)
			dataChannel.SendText(asciiArt)
			time.Sleep(100 * time.Millisecond)
			// fmt.Print("\033[H\033[2J")
			// fmt.Println(asciiArt)
		}
	})

	select {}

}

// func gocvInit(webcam gocv.VideoCapture, img gocv.Mat){
// 	for {
// 		webcam.Read(&img)
// 		goImg, _ := img.ToImage()

// 		asciiArt := toASCIIArt(goImg, asciiChars)
// 		datachannel.SendText(asciiArt)
// 		time.Sleep(100 * time.Millisecond)
// 		// fmt.Print("\033[H\033[2J")
// 		// fmt.Println(asciiArt)
// 	}
// }

func toASCIIArt(img image.Image, charMap string) string {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var sb strings.Builder

	for y := 0; y < height; y = y + 15 {
		for x := 0; x < width; x = x + 10 {
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
