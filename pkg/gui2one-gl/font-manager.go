package gui2onegl

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)


func rasterizeGlyph(font *sfnt.Font, glyph rune, fontSize int) image.Image {
  buf := new(sfnt.Buffer)

	idx, _ := font.GlyphIndex(buf,glyph)
	// println("a" , idx)
    segs, _ := font.LoadGlyph(buf,idx, fixed.Int26_6(fontSize << 6), nil)

		minX := segs.Bounds().Min.X.Floor()
		maxX := segs.Bounds().Max.X.Ceil()

		minY := segs.Bounds().Min.Y.Floor()
		maxY := segs.Bounds().Max.Y.Ceil()
		
		fmt.Println(minX, maxX, minY, maxY)
    r := vector.NewRasterizer(fontSize, fontSize)
	
	// Simple transform values

	scale := float32(1/ 64.0)
	offsetX := float32(1)
	offsetY := float32(50) // baseline

	for _, seg := range segs {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			r.MoveTo(
				offsetX+scale * float32(seg.Args[0].X),
				offsetY+scale * float32(seg.Args[0].Y),
			)

		case sfnt.SegmentOpLineTo:
			r.LineTo(
				offsetX+scale * float32(seg.Args[0].X),
				offsetY+scale * float32(seg.Args[0].Y),
			)

		case sfnt.SegmentOpQuadTo:
			r.QuadTo(
				offsetX+scale * float32(seg.Args[0].X),
				offsetY+scale * float32(seg.Args[0].Y),
				offsetX+scale * float32(seg.Args[1].X),
				offsetY+scale * float32(seg.Args[1].Y),
			)

		case sfnt.SegmentOpCubeTo:
			r.CubeTo(
				offsetX+scale * float32(seg.Args[0].X),
				offsetY+scale * float32(seg.Args[0].Y),
				offsetX+scale * float32(seg.Args[1].X),
				offsetY+scale * float32(seg.Args[1].Y),
				offsetX+scale * float32(seg.Args[2].X),
				offsetY+scale * float32(seg.Args[2].Y),
			)
		}
	}

    img := image.NewRGBA(image.Rect(0, 0, fontSize, fontSize))

	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
    r.Draw(img, img.Bounds(), image.White, image.Point{})

	return img

}
func GenerateAtlas(fontFilePath string) {
    fontFile, _ := os.ReadFile(fontFilePath)

    font, _ := sfnt.Parse(fontFile)
	img := rasterizeGlyph(font, 'L', 64)

    finalIMG := image.NewRGBA(image.Rect(0, 0, 640, 480))
	
    draw.Draw(finalIMG, img.Bounds(), img, image.Point{}, draw.Src)
	
    f, _ := os.Create("out.png")
    defer f.Close()
    png.Encode(f, finalIMG)
}