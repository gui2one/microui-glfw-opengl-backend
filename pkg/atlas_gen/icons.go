package atlas_gen

import (
	"image"
	"image/png"
	"os"

	xdraw "golang.org/x/image/draw"
)

func LoadIcon(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GetResizedIcon(src image.Image, width, height int) *image.RGBA {
	// 1. Create the destination canvas with your desired size
	dstRect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dstRect)

	// 2. Scale the source image into the destination
	// xdraw.BiLinear is a good balance between speed and quality
	xdraw.BiLinear.Scale(dst, dstRect, src, src.Bounds(), xdraw.Src, nil)

	return dst
}
