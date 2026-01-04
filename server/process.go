package server

import (
	"bytes"
	"image"
	"log"
)

func Process(imgBytes []byte, clientAddr string) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Print(err)
		return
	}
	err = SaveAsJPG(img, len(imgBytes), clientAddr)
	if err != nil {
		log.Print(err)
	}
}
