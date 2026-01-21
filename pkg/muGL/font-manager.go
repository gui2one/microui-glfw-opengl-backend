package muGL

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

func (g *GlyphMetrics) Print() {
	fmt.Println("Glyph Metrics ---->")
	fmt.Println("    GlyphMetrics --")
	fmt.Println("      IDX :", g.IDX)

	fmt.Printf("      UnicodeID : 0x%04x --> %c\n", g.UnicodeID, g.UnicodeID)
	fmt.Println("      Coords :", g.X, g.Y)

	fmt.Println("      AdvanceX :", g.AdvanceX)
	fmt.Println("      BearingX :", g.BearingX)
	fmt.Println("      BearingY :", g.BearingY)
	fmt.Println("      Width :", g.Width)
	fmt.Println("      Height :", g.Height)
	fmt.Println("-----------------------------")
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
	FontSize    int
	FontName    string
	Width       int
	Height      int
	Atlas       *image.RGBA
	FontMetrics *FontMetrics
	Glyphs      []*GlyphMetrics
	GlyphsRange [2]int
	// colors and icons
	Black       Rect
	White       Rect
	CloseIcon   Rect
	CheckedIcon Rect
}

func (a *AtlasData) Print(showGlyphs bool) {
	fmt.Println("AtlasData ---->")
	fmt.Println("  FontSize :", a.FontSize)
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
	maxX := segs.Bounds().Max.X.Ceil() + 1

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
	// Load glyph scaled to fontSize (this returns units in 26.6 fixed point)
	segs, _ := font.LoadGlyph(buf, idx, fixed.Int26_6(fontSize<<6), nil)

	// Calculate bounds in pixels (divide by 64)
	b := segs.Bounds()
	minX := float32(b.Min.X) / 64.0
	minY := float32(b.Min.Y) / 64.0

	r := vector.NewRasterizer(fontSize, fontSize)

	// We want to translate the glyph so it fits in our [fontSize x fontSize] box.
	// We subtract minX/minY to move the glyph's top-left to (0,0).
	// Adding a small padding (like 2) is fine, but be careful not to exceed fontSize.
	offsetX := -minX
	offsetY := -minY

	for _, seg := range segs {
		// Helper to convert 26.6 Fixed to Float pixels
		f := func(v fixed.Point26_6) (float32, float32) {
			return offsetX + (float32(v.X) / 64.0), offsetY + (float32(v.Y) / 64.0)
		}

		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			x, y := f(seg.Args[0])
			r.MoveTo(x, y)
		case sfnt.SegmentOpLineTo:
			x, y := f(seg.Args[0])
			r.LineTo(x, y)
		case sfnt.SegmentOpQuadTo:
			x1, y1 := f(seg.Args[0])
			x2, y2 := f(seg.Args[1])
			r.QuadTo(x1, y1, x2, y2)
		case sfnt.SegmentOpCubeTo:
			x1, y1 := f(seg.Args[0])
			x2, y2 := f(seg.Args[1])
			x3, y3 := f(seg.Args[2])
			r.CubeTo(x1, y1, x2, y2, x3, y3)
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, fontSize, fontSize))
	r.Draw(img, img.Bounds(), image.White, image.Point{})

	glypMetrics := getGlyphMetrics(font, idx, segs, fontSize)
	return img, glypMetrics
}

func GenerateAtlas(fontFilePath string, glyphsRange [2]int, fontSize int) *AtlasData {

	fontFile, err := os.ReadFile(fontFilePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, _ := sfnt.Parse(fontFile)

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

	result := &AtlasData{}
	step := int(fontSize)
	startX := bufferNum % numCols * step
	startY := bufferNum / numCols * step

	// draw black cell into Atlas
	draw.Draw(finalIMG, image.Rect(0, 0, fontSize, fontSize), image.Black, image.Point{}, draw.Src)
	// draw white cell into Atlas
	draw.Draw(finalIMG, image.Rect(fontSize*1, 0, fontSize*2, fontSize), image.White, image.Point{}, draw.Src)

	cellStep := float32(fontSize) / float32(finalDIM)
	result.Black = Rect{
		P1: Point{
			X: 0 + 0.02,
			Y: (1 - cellStep) + 0.02, /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStep*1.0 - 0.02,
			Y: 0.98, /* should be 1 but border issue with shading */
		},
	}
	result.White = Rect{
		P1: Point{
			X: cellStep + 0.02,
			Y: (1 - cellStep) + 0.02, /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStep*2.0 - 0.02,
			Y: 0.98, /* should be 1 but border issue with shading */
		},
	}
	result.CloseIcon = Rect{
		P1: Point{
			X: cellStep * 2,
			Y: (1 - cellStep), /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStep * 3.0,
			Y: 1.0, /* should be 1 but border issue with shading */
		},
	}
	result.CheckedIcon = Rect{
		P1: Point{
			X: cellStep * 3,
			Y: 1.0, /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStep * 4,
			Y: (1 - cellStep), /* should be 1 but border issue with shading */
		},
	}

	// add icons to atlas
	closeImage, err := LoadIcon("assets/icons/close.png")
	if err != nil {
		log.Println(err)
		return nil
	}
	iconRGBA := GetResizedIcon(closeImage, fontSize, fontSize)
	draw.Draw(finalIMG, image.Rect(fontSize*2, 0, fontSize*3, fontSize),
		iconRGBA,
		image.Point{},
		draw.Src,
	)

	checkedImage, err := LoadIcon("assets/icons/checked.png")
	if err != nil {
		log.Println(err)
		return nil
	}
	iconRGBA = GetResizedIcon(checkedImage, fontSize, fontSize)
	draw.Draw(finalIMG, image.Rect(fontSize*3, 0, fontSize*4, fontSize),
		iconRGBA,
		image.Point{},
		draw.Src,
	)
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

	result.FontSize = fontSize
	result.FontName = path.Base(fontFilePath)
	result.Width = finalDIM
	result.Height = finalDIM
	result.Atlas = finalIMG
	result.FontMetrics = fontMetrics
	result.GlyphsRange = glyphsRange
	result.Glyphs = glyphs_metrics

	return result

}
