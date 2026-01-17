package gui2onegl

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Texture struct. Mainly holds an ID
type Texture struct {
	ID     uint32
	Width  uint32
	Height uint32
}

// NewTexture creates a new Texture*
func NewTexture() *Texture {
	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE, 0)
	return &Texture{
		ID:     id,
		Width:  1,
		Height: 1,
	}
}

// Bind binds the texture for rendering
func (t Texture) Bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

// BindWithActiveUnit binds the texture for rendering and specify the texture unit
func (t Texture) BindWithActiveUnit(num uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + num)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

// GenerateSolidTexture generates a texture with a solid color
func GenerateSolidTexture(width, height int) *Texture {

	pixels := make([]byte, width*height*4)
	for i := 0; i < len(pixels); i += 4 {
		pixels[i+0] = 0
		pixels[i+1] = 0
		pixels[i+2] = 255
		pixels[i+3] = 255
	}
	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &Texture{
		ID: id,
	}
}

// GenerateTexture generates a texture with a solid color
func GenerateTexture(width, height int, pixels []byte) *Texture {

	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &Texture{
		ID:     id,
		Width:  uint32(width),
		Height: uint32(height),
	}
}

// GenerateTypedTexture generates a texture
func GenerateTypedTexture(width, height int, internalFormat int32, format uint32) *Texture {
	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, int32(width), int32(height), 0, format, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &Texture{
		ID:     id,
		Width:  uint32(width),
		Height: uint32(height),
	}
}

// SetTextureParams sets the texture parameters
func SetTextureParams(t *Texture, width int, height int, internalFormat int32, format uint32) {
	t.Width = uint32(width)
	t.Height = uint32(height)
	t.Bind()
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, int32(width), int32(height), 0, format, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	// gl.BindTexture(gl.TEXTURE_2D, 0)
}
func SetTextureData(t *Texture, width int, height int, internalFormat int32, format uint32, data []byte) {
	t.Width = uint32(width)
	t.Height = uint32(height)
	t.Bind()
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, int32(width), int32(height), 0, format, gl.UNSIGNED_BYTE, gl.Ptr(data))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}

// LoadJPEG loads a texture from a file
func LoadImageFile(imgFilePath string) (*Texture, error) {

	ext := path.Ext(imgFilePath)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".JPG" && ext != ".JPEG" && ext != ".png" && ext != ".PNG" {
		return &Texture{}, fmt.Errorf("unsupported file type: %s", ext)
	}
	file, err := os.Open(imgFilePath)
	if err != nil {
		fmt.Println(err)
		return &Texture{}, err
	}
	var img image.Image
	switch ext {
	case ".jpg", ".jpeg", ".JPG", ".JPEG":
		img, err = jpeg.Decode(file)
		if err != nil {
			fmt.Println(err)
			return &Texture{}, err
		}
	case ".png", ".PNG":
		img, err = png.Decode(file)
		if err != nil {
			fmt.Println(err)
			return &Texture{}, err
		}

	}

	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	bytes := make([]byte, w*h*4)

	for i := 0; i < len(bytes); i += 4 {
		clr := img.At((i/4)%w, h-((i/4)/w))

		bytes[i+0], bytes[i+1], bytes[i+2], bytes[i+3] = rgbaToPixel(clr.RGBA())

	}
	// fmt.Printf("dimensions : %v\n", img.Bounds().Max)

	texture := GenerateTexture(w, h, bytes)
	return texture, err
}
func FromImage(img image.Image) *Texture {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	bytes := make([]byte, w*h*4)

	for i := 0; i < len(bytes); i += 4 {
		clr := img.At((i/4)%w, h-((i/4)/w))

		bytes[i+0], bytes[i+1], bytes[i+2], bytes[i+3] = rgbaToPixel(clr.RGBA())
	}

	return GenerateTexture(w, h, bytes)
}
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) (byte, byte, byte, byte) {
	return byte(r / 256), byte(g / 256), byte(b / 256), byte(a / 256)
}
