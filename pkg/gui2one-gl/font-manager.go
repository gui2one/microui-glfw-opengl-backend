package gui2onegl

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

type GlyphMetrics struct {
	IDX       sfnt.GlyphIndex
	UnicodeID uint16
	AdvanceX  int
	BearingX  int
	BearingY  int

	X, Y   int
	Width  int
	Height int
}
type FontMetrics struct {
	Ascent     int
	Descent    int
	LineHeight int
}

func (m *FontMetrics) Print() {
	fmt.Println("Font Metrics ---->")
	fmt.Println("  Ascent :", m.Ascent)
	fmt.Println("  Descent :", m.Descent)
	fmt.Println("  LineHeight :", m.LineHeight)
}

type AtlasData struct {
	FontName    string
	Width       int
	Height      int
	Atlas       *image.RGBA
	FontMetrics *FontMetrics
	Glyphs      []*GlyphMetrics
}

func (a *AtlasData) Print(showGlyphs bool) {
	fmt.Println("AtlasData ---->")
	fmt.Println("  FontName :", a.FontName)
	fmt.Println("  Width :", a.Width)
	fmt.Println("  Height :", a.Height)

	a.FontMetrics.Print()
	fmt.Println(len(a.Glyphs), "  Glyphs ---->")
	if showGlyphs {

		for _, g := range a.Glyphs {
			fmt.Println("    GlyphMetrics --")
			fmt.Println("      IDX :", g.IDX)

			fmt.Printf("      UnicodeID : 0x%04x --> %c\n", g.UnicodeID, g.UnicodeID)
			fmt.Println("      Coords :", g.X, g.Y)

			fmt.Println("      AdvanceX :", g.AdvanceX)
			fmt.Println("      BearingX :", g.BearingX)
			fmt.Println("      BearingY :", g.BearingY)
			fmt.Println("      Width :", g.Width)
			fmt.Println("      Height :", g.Height)
		}
	}
}

func getFontMetrics(font *sfnt.Font, fontSize int) *FontMetrics {
	buf := new(sfnt.Buffer)
	metrics, _ := font.Metrics(buf, fixed.Int26_6(fontSize<<6), 0)
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Floor()
	lineHeight := metrics.Height.Ceil()

	fontMetrics := &FontMetrics{
		Ascent:     ascent,
		Descent:    descent,
		LineHeight: lineHeight,
	}

	return fontMetrics
}

func getGlyphMetrics(font *sfnt.Font, glyphIndex sfnt.GlyphIndex, segs sfnt.Segments, fontSize int) *GlyphMetrics {
	buf := new(sfnt.Buffer)
	adv, _ := font.GlyphAdvance(buf, glyphIndex, fixed.Int26_6(fontSize<<6), 0)

	bounds := segs.Bounds()

	bearingX := bounds.Min.X.Floor()
	bearingY := bounds.Max.Y.Floor()

	minX := segs.Bounds().Min.X.Floor()
	maxX := segs.Bounds().Max.X.Ceil()

	minY := segs.Bounds().Min.Y.Floor()
	maxY := segs.Bounds().Max.Y.Ceil()

	glyphHeight := maxY - minY
	glyphWidth := maxX - minX
	return &GlyphMetrics{
		IDX:      glyphIndex,
		AdvanceX: adv.Ceil(),
		BearingX: bearingX,
		BearingY: bearingY,
		Width:    glyphWidth,
		Height:   glyphHeight,
	}

}

func rasterizeGlyph(font *sfnt.Font, idx sfnt.GlyphIndex, fontSize int) (*image.RGBA, *GlyphMetrics) {
	buf := new(sfnt.Buffer)
	if idx == 0 {
		log.Println("glyph not found")
		return nil, nil
	}
	segs, _ := font.LoadGlyph(buf, idx, fixed.Int26_6(fontSize<<6), nil)

	minX := segs.Bounds().Min.X.Floor()
	minY := segs.Bounds().Min.Y.Floor()

	r := vector.NewRasterizer(fontSize, fontSize)

	// Simple transform values
	scale := float32(1.0 / float32(fontSize))
	offsetX := float32(-minX + 2)
	offsetY := float32(-minY + 2) // baseline

	// "apply" segmentOps
	for _, seg := range segs {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			r.MoveTo(
				offsetX+scale*float32(seg.Args[0].X),
				offsetY+scale*float32(seg.Args[0].Y),
			)

		case sfnt.SegmentOpLineTo:
			r.LineTo(
				offsetX+scale*float32(seg.Args[0].X),
				offsetY+scale*float32(seg.Args[0].Y),
			)

		case sfnt.SegmentOpQuadTo:
			r.QuadTo(
				offsetX+scale*float32(seg.Args[0].X),
				offsetY+scale*float32(seg.Args[0].Y),
				offsetX+scale*float32(seg.Args[1].X),
				offsetY+scale*float32(seg.Args[1].Y),
			)

		case sfnt.SegmentOpCubeTo:
			r.CubeTo(
				offsetX+scale*float32(seg.Args[0].X),
				offsetY+scale*float32(seg.Args[0].Y),
				offsetX+scale*float32(seg.Args[1].X),
				offsetY+scale*float32(seg.Args[1].Y),
				offsetX+scale*float32(seg.Args[2].X),
				offsetY+scale*float32(seg.Args[2].Y),
			)
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, fontSize, fontSize))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
	r.Draw(img, img.Bounds(), image.White, image.Point{})

	glypMetrics := getGlyphMetrics(font, idx, segs, fontSize)

	return img, glypMetrics

}
func GenerateAtlas(fontFilePath string, glyphsRange [2]int) *AtlasData {
	fontFile, err := os.ReadFile(fontFilePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, _ := sfnt.Parse(fontFile)
	fontSize := int(64)

	images := []*image.RGBA{}
	glyphs_metrics := []*GlyphMetrics{}
	if glyphsRange[0] > glyphsRange[1] {
		log.Println("bad glyphs range given...")
		return nil
	}
	for i := glyphsRange[0]; i <= glyphsRange[1]; i++ {
		var buf sfnt.Buffer
		glyphIndex, err := font.GlyphIndex(&buf, rune(i))
		if err != nil || glyphIndex == 0 {
			continue
		}
		img3, metrics := rasterizeGlyph(font, glyphIndex, fontSize)
		if img3 == nil {
			continue
		}
		images = append(images, img3)
		metrics.UnicodeID = uint16(i)
		glyphs_metrics = append(glyphs_metrics, metrics)
	}
	bufferNum := 30 /* number of empty spaces for White, icons and stuff */
	numCols := int(math.Ceil(math.Sqrt(float64(len(images) + bufferNum))))
	finalDIM := numCols * fontSize
	finalIMG := image.NewRGBA(image.Rect(0, 0, int(finalDIM), int(finalDIM)))

	step := int(fontSize)
	startX := bufferNum % numCols * step
	startY := bufferNum / numCols * step

	// draw blacl cell into Atlas
	draw.Draw(finalIMG, image.Rect(0, 0, fontSize, fontSize), image.Black, image.Point{}, draw.Src)
	// draw white cell into Atlas
	draw.Draw(finalIMG, image.Rect(fontSize*1, 0, fontSize*2, fontSize), image.White, image.Point{}, draw.Src)

	for i, img := range images {
		dstRect := image.Rect(
			int(startX),
			int(startY),
			int(startX)+img.Bounds().Dx(),
			int(startY)+img.Bounds().Dy(),
		)
		draw.Draw(finalIMG, dstRect, img, image.Point{}, draw.Src)
		glyphs_metrics[i].X = startX
		glyphs_metrics[i].Y = startY

		// prepare newt step
		startX += step
		if startX >= int(finalDIM) {
			startX = 0
			startY += step
		}
	}

	// write file on disk ... for now
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, finalIMG)

	fontMetrics := getFontMetrics(font, fontSize)

	atlasData := &AtlasData{
		FontName:    path.Base(fontFilePath),
		Width:       finalDIM,
		Height:      finalDIM,
		Atlas:       finalIMG,
		FontMetrics: fontMetrics,
		Glyphs:      glyphs_metrics,
	}

	return atlasData

}
