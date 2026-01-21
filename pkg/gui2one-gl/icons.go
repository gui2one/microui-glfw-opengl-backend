package gui2onegl

import (
	"image"
	"image/png"
	"os"
)

func CloseIcon() (image.Image, error) {
	file, err := os.Open("assets/icons/close.png")
	if err != nil {
		return nil, err

	}

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
