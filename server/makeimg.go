package server

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"tcp/test/constants"
)

func SaveAsJPG(img image.Image, length int, clientAddr string) error {
	_, err := os.Stat(constants.Dest)
	if err != nil {
		err := os.Mkdir(constants.Dest, 0755)
		if err != nil {
			return err
		}
	}

	size := img.Bounds().Size()
	name := fmt.Sprintf(clientAddr+"%v", size)
	filepath := fmt.Sprintf(constants.Dest+"/%v.jpg", name)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{
		Quality: 70,
	})
}
