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

	_font "golang.org/x/image/font" // This contains the Hinting constants
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
	Black         Rect
	White         Rect
	CloseIcon     Rect
	CheckedIcon   Rect
	CollapsedIcon Rect
	ExpandedIcon  Rect
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
			fmt.Printf("    GlyphMetrics -- %c\n", g.UnicodeID)
			fmt.Println("      IDX :", g.IDX)

			fmt.Printf("      UnicodeID : 0x%04x : %c\n", g.UnicodeID, g.UnicodeID)
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

	// 1. Calculate scale (assuming 72 DPI for simplicity, or pass DPI as a param)
	// fixed.I(x) is shorthand for x << 6
	scale := fixed.I(fontSize)

	// 2. Get Advance
	adv, _ := font.GlyphAdvance(buf, glyphIndex, scale, _font.HintingNone)

	// 3. Get Bounds once
	bounds := segs.Bounds()

	// 4. Calculate dimensions using Ceil/Floor appropriately
	// Use Floor for Min and Ceil for Max to ensure the glyph fits
	// entirely within the integer pixel boundaries.
	minX, maxX := bounds.Min.X.Floor(), bounds.Max.X.Ceil()
	minY, maxY := bounds.Min.Y.Floor(), bounds.Max.Y.Ceil()

	return &GlyphMetrics{
		IDX:      glyphIndex,
		AdvanceX: adv.Floor(),
		BearingX: minX,
		BearingY: maxY, // Usually the offset from baseline to top
		Width:    maxX - minX,
		Height:   maxY - minY,
	}
}

func rasterizeGlyph(font *sfnt.Font, idx sfnt.GlyphIndex, fontSize int) (*image.RGBA, *GlyphMetrics) {

	buf := new(sfnt.Buffer)
	// Load glyph scaled to fontSize (this returns units in 26.6 fixed point)
	segs, _ := font.LoadGlyph(buf, idx, fixed.Int26_6(fontSize<<6), nil)
	metrics, err := font.Metrics(buf, fixed.Int26_6(fontSize<<6), 0)
	if err != nil {
		panic(err)
	}
	// 1. Get the metrics in floating point
	fAscent := float32(metrics.Ascent) / 64.0
	fDescent := float32(metrics.Descent) / 64.0 // Distance below baseline

	// 2. The height of the font 'line' is actually:
	totalLineHeight := fAscent + fDescent

	offsetX := float32(0.0)
	offsetY := fAscent - fDescent
	canvasWidth := fontSize
	canvasHeight := int(math.Ceil(float64(totalLineHeight))) + 1
	r := vector.NewRasterizer(canvasWidth, canvasHeight)

	for _, seg := range segs {
		f := func(v fixed.Point26_6) (float32, float32) {
			x := offsetX + float32(v.X)/64.0
			y := offsetY + float32(v.Y)/64.0
			return x, y
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

	glypMetrics := getGlyphMetrics(font, idx, segs, fontSize)

	img := image.NewRGBA(image.Rect(0, 0, fontSize, int(totalLineHeight)))
	r.Draw(img, img.Bounds(), image.White, image.Point{})

	return img, glypMetrics
}

func addIconToAtlas(iconPath string, finalImage *image.RGBA, fontSize int, cellStepX float32, cellStepY float32, col int, row int) Rect {
	closeImage, err := LoadIcon(iconPath)
	if err != nil {
		log.Println(err)
		return Rect{}
	}
	iconRGBA := GetResizedIcon(closeImage, fontSize, fontSize)
	draw.Draw(finalImage, image.Rect(fontSize*col, 0, fontSize*(col+1), fontSize),
		iconRGBA,
		image.Point{},
		draw.Src,
	)

	result := Rect{
		P1: Point{
			X: cellStepX * float32(col),
			Y: float32(row + 1),
		},
		P2: Point{
			X: cellStepX * float32(col+1),
			Y: (float32(row+1) - cellStepY),
		},
	}

	return result
}
func GenerateAtlas(fontFilePath string, glyphsRange [2]int, fontSize int) *AtlasData {

	fontFile, err := os.ReadFile(fontFilePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, _ := sfnt.Parse(fontFile)
	fontMetrics := getFontMetrics(font, fontSize)
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
	stepX := int(fontSize)
	stepY := int(fontMetrics.Ascent + fontMetrics.Descent)
	finalWidth := numCols * stepX
	finalHeight := numCols * stepY
	finalIMG := image.NewRGBA(image.Rect(0, 0, int(finalWidth), int(finalHeight)))

	result := &AtlasData{}
	startX := bufferNum % numCols * stepX
	startY := bufferNum / numCols * stepY

	cellStepX := float32(fontSize) / float32(finalWidth)
	cellStepY := float32(fontSize) / float32(finalHeight)
	// draw black cell into Atlas
	draw.Draw(finalIMG, image.Rect(0, 0, stepX, stepY), image.Black, image.Point{}, draw.Src)
	result.Black = Rect{
		P1: Point{
			X: 0 + 0.02,
			Y: (1 - cellStepY) + 0.02, /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStepX*1.0 - 0.02,
			Y: 0.98, /* should be 1 but border issue with shading */
		},
	}
	// draw white cell into Atlas
	draw.Draw(finalIMG, image.Rect(stepX*1, 0, stepX*2, stepY), image.White, image.Point{}, draw.Src)
	result.White = Rect{
		P1: Point{
			X: cellStepX + 0.02,
			Y: (1 - cellStepY) + 0.02, /* +0.01 = border issue with shading */
		},
		P2: Point{
			X: cellStepX*2.0 - 0.02,
			Y: 0.98, /* should be 1 but border issue with shading */
		},
	}

	result.CloseIcon = addIconToAtlas("assets/icons/close.png", finalIMG, fontSize, cellStepX, cellStepY, 2, 0)
	result.CheckedIcon = addIconToAtlas("assets/icons/checked.png", finalIMG, fontSize, cellStepX, cellStepY, 3, 0)
	result.CollapsedIcon = addIconToAtlas("assets/icons/collapsed.png", finalIMG, fontSize, cellStepX, cellStepY, 4, 0)
	result.ExpandedIcon = addIconToAtlas("assets/icons/expanded.png", finalIMG, fontSize, cellStepX, cellStepY, 5, 0)

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
		startX += stepX
		if startX >= int(finalWidth) {
			startX = 0
			startY += stepY
		}
	}

	// write file on disk ... for now
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, finalIMG)

	result.FontSize = fontSize
	result.FontName = path.Base(fontFilePath)
	result.Width = finalWidth
	result.Height = finalHeight
	result.Atlas = finalIMG
	result.FontMetrics = fontMetrics
	result.GlyphsRange = glyphsRange
	result.Glyphs = glyphs_metrics

	return result

}
