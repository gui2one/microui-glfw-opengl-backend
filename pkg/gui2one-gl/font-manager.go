package gui2onegl

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"slices"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"

	"github.com/depp/skelly64/lib/rectpack"
)


func rasterizeGlyph(font *sfnt.Font, idx sfnt.GlyphIndex, fontSize int) *image.RGBA {
  buf := new(sfnt.Buffer)

	// idx, _ := font.GlyphIndex(buf,glyph)

	if idx == 0 {
		log.Println("ha !!!!!!!!!!!!!!!!!!!!!!")
		return nil
	}
    segs, _ := font.LoadGlyph(buf,idx, fixed.Int26_6(fontSize << 6), nil)

		minX := segs.Bounds().Min.X.Floor()
		maxX := segs.Bounds().Max.X.Ceil()

		minY := segs.Bounds().Min.Y.Floor()
		maxY := segs.Bounds().Max.Y.Ceil()

		glyphHeight := maxY - minY
		glyphWidth := maxX - minX

padding := 2
w := glyphWidth + padding*2
h := glyphHeight + padding*2

r := vector.NewRasterizer(w, h)
	
	// Simple transform values

	scale := float32(1.0/ 64.0)
	offsetX := float32(-minX + 2)
	offsetY := float32(-minY + 2) // baseline

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


img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
    r.Draw(img, img.Bounds(), image.White, image.Point{})

	return img

}

func GenerateAtlas(fontFilePath string) {
    fontFile, _ := os.ReadFile(fontFilePath)

    font, _ := sfnt.Parse(fontFile)


	images := []*image.RGBA{}
	for i := 66; i < 128; i++ {
	var buf sfnt.Buffer
	glyphIndex, err := font.GlyphIndex(&buf, rune(i))
	if err != nil || glyphIndex == 0 {
		continue
	}		
		img3 := rasterizeGlyph(font, glyphIndex, 64)
		if img3 == nil {
			continue
		}
		images = append(images, img3)
	}

	rectangles := []rectpack.Point{}

	for _, img := range images {
		rectangles = append(rectangles, rectpack.Point{X: int32((*img).Bounds().Dx()), Y: int32((*img).Bounds().Dy())})
	}

	fmt.Println("rectangles : ", rectangles)




	p := rectpack.New()
	
	bounds, pos, err := rectpack.AutoPackSingle(p, rectangles)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bounds : ", bounds)
	fmt.Println("positions : ", pos)

	// sorting rectangles to match rectpack heuristic ?!
	slices.SortFunc(rectangles, func(a, b rectpack.Point) int {
		return int(b.Y - a.Y)
	})
fmt.Println("rectangles : ", rectangles)	
	finalIMG := image.NewRGBA(image.Rect(0, 0, 640, 480))
	
	img := images[1]
    draw.Draw(finalIMG, img.Bounds(), img, image.Point{}, draw.Src)
	
    f, _ := os.Create("out.png")
    defer f.Close()
    png.Encode(f, finalIMG)
}