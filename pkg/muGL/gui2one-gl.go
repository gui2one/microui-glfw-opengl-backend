package muGL

import (
	"fmt"

	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/zeozeozeo/microui-go"
	mu "github.com/zeozeozeo/microui-go"
)

var GLYPHS_RANGE = [2]int{0x0020, 0x007E}

func InitGL() {

	if gl.Init() != nil {
		panic("Unable to initialize OpenGL")
	}
	SetupGLDebug()
}

// App base app structure
type App struct {
	MeshBuffer         *GlMeshData
	MainShader         uint32
	AtlasData          AtlasData
	AtlasTexture       Texture
	NumFloatsPerVertex int
	SizeOfFloat32      int

	Width  int
	Height int
}

func (a *App) InitGL() {
	a.NumFloatsPerVertex = 2 + 2 + 3
	a.SizeOfFloat32 = 4

	a.MainShader = generateShader("assets/shaders/main_vertex.glsl", "assets/shaders/main_fragment.glsl")

	a.MeshBuffer = NewGlMeshData()
	m := a.MeshBuffer
	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.IndexBuffer)

	stride := int32(a.SizeOfFloat32 * a.NumFloatsPerVertex)
	gl.BindVertexArray(m.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)

	/* position 2d */
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, stride, 0)
	/* uvs */
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(2*a.SizeOfFloat32))
	/* color RGB */
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, stride, uintptr((2+2)*a.SizeOfFloat32))

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Texture ATLAS
	atlasData := GenerateAtlas("assets/fonts/ConsolaMono-Bold.TTF", GLYPHS_RANGE, 36)
	a.AtlasData = *atlasData
	a.AtlasTexture = *FromImage(atlasData.Atlas)
	atlasData.Print(true)

	// panic("OK stop here !!!!")

	a.AtlasTexture.Bind()

}

func (app *App) InitMuContext(ctx *microui.Context) {
	myFontHandle := app.AtlasData

	ctx.Style.Font = myFontHandle
	ctx.TextHeight = app.TextHeight
	ctx.TextWidth = app.TextWidth
}
func (a *App) TextWidth(font microui.Font, text string) int {
	w := a.ComputeTextWidth(text)
	return w
}
func (a *App) TextHeight(font microui.Font) int {
	return a.AtlasData.FontMetrics.LineHeight
}
func (a *App) PushRect(x, y, w, h float32, uvs Rect, color [3]float32) {
	m := a.MeshBuffer

	numVertices := len(m.Vertices) / a.NumFloatsPerVertex
	rect := Rect{
		P1: Point{X: x, Y: y},
		P2: Point{X: x + w, Y: y + h},
	}
	vertices := []float32{
		/*pos */ rect.P1.X, rect.P1.Y /*uvs */, uvs.P1.X, uvs.P1.Y /* color */, color[0], color[1], color[2],
		/*pos */ rect.P2.X, rect.P1.Y /*uvs */, uvs.P2.X, uvs.P1.Y /* color */, color[0], color[1], color[2],
		/*pos */ rect.P2.X, rect.P2.Y /*uvs */, uvs.P2.X, uvs.P2.Y /* color */, color[0], color[1], color[2],
		/*pos */ rect.P1.X, rect.P2.Y /*uvs */, uvs.P1.X, uvs.P2.Y /* color */, color[0], color[1], color[2],
	}
	m.Vertices = append(m.Vertices, vertices...)

	indices := []uint32{
		uint32(numVertices) + 0, uint32(numVertices) + 1, uint32(numVertices) + 2,
		uint32(numVertices) + 2, uint32(numVertices) + 3, uint32(numVertices) + 0,
	}
	m.Indices = append(m.Indices, indices...)
}
func (a *App) PushText(x, y float32, text string, color [3]float32) {
	penX := x
	penY := y

	glyphsRange := a.AtlasData.GlyphsRange
	for _, c := range text {
		if c == '\n' {
			penX = x
			penY += float32(a.AtlasData.FontMetrics.LineHeight)
		}
		if c >= rune(glyphsRange[0]) && c <= rune(glyphsRange[1]) {
			glyph := a.AtlasData.Glyphs[c-rune(glyphsRange[0])]
			// aaa := (float32(a.AtlasData.FontSize) - float32(glyph.Height) + float32(glyph.BearingY)) / float32(a.AtlasData.Width)

			uvStartX := float32(glyph.X) / float32(a.AtlasData.Width)
			uvStartY := float32(glyph.Y) / float32(a.AtlasData.Height)
			uvW := float32(glyph.Width) / float32(a.AtlasData.Width)
			uvH := float32(glyph.Height) / float32(a.AtlasData.Height)
			uvsRect := Rect{
				P1: Point{X: uvStartX, Y: 1.0 - uvStartY},
				P2: Point{X: uvStartX + uvW, Y: 1.0 - uvStartY - uvH},
			}

			drawY := penY
			a.PushRect(penX+float32(glyph.BearingX), drawY, float32(glyph.Width), float32(a.AtlasData.FontSize), uvsRect, color)
			penX += float32(glyph.AdvanceX)

		}
	}
}
func (a *App) ComputeTextWidth(text string) int {
	penX := 0
	maxWidth := 0
	curWidth := 0
	glyphsRange := a.AtlasData.GlyphsRange
	startRune := rune(glyphsRange[0])
	endRune := rune(glyphsRange[1])
	for _, c := range text {
		if c == '\n' {
			penX = 0
			curWidth = 0
		}
		if c >= startRune && c <= endRune {
			glyph := a.AtlasData.Glyphs[c-startRune]
			penX += int(glyph.AdvanceX)
		}
		curWidth = int(penX)
		if curWidth > maxWidth {
			maxWidth = curWidth
		}
	}

	// fmt.Println("text : ", text, "\nwidth : ", maxWidth)
	return maxWidth
}
func (a *App) FlushRects() {
	m := a.MeshBuffer

	if len(m.Vertices) == 0 {
		return
	}
	gl.BindVertexArray(m.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, a.SizeOfFloat32*len(m.Vertices), gl.Ptr(m.Vertices), gl.STREAM_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.IndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, a.SizeOfFloat32*len(m.Indices), gl.Ptr(m.Indices), gl.STREAM_DRAW)

}

func (a *App) ClearRects() {

	a.MeshBuffer.Vertices = a.MeshBuffer.Vertices[:0]
	a.MeshBuffer.Indices = a.MeshBuffer.Indices[:0]
}
func (app *App) SetScissor(r microui.Rect) {
	// Convert top-left Y to bottom-left Y
	x := r.X
	y := app.Height - (r.Y + r.H)
	w := r.W
	h := r.H
	gl.Scissor(
		int32(x),
		int32(y),
		int32(w),
		int32(h),
	)
	gl.Enable(gl.SCISSOR_TEST)
}

func (app *App) Render(ctx *microui.Context) {
	PrepareGLobalState(app)
	app.ClearRects()
	gl.Disable(gl.SCISSOR_TEST) // Start with no scissor
	for _, cmd := range ctx.CommandList {
		switch cmd.Type {
		case microui.MU_COMMAND_CLIP:
			DrawMyStuff(app)
			app.ClearRects()
			app.SetScissor(cmd.Clip.Rect)

		case microui.MU_COMMAND_RECT:

			rgba := cmd.Rect.Color.ToRGBA()
			app.PushRect(float32(cmd.Rect.Rect.X), float32(cmd.Rect.Rect.Y), float32(cmd.Rect.Rect.W), float32(cmd.Rect.Rect.H),
				app.AtlasData.White,
				[3]float32{float32(rgba.R) / 255.0, float32(rgba.G) / 255.0, float32(rgba.B) / 255.0},
			)

		case microui.MU_COMMAND_TEXT:

			clr := cmd.Text.Color.ToRGBA()
			app.PushText(
				float32(cmd.Text.Pos.X),
				float32(cmd.Text.Pos.Y),
				cmd.Text.Str,
				[3]float32{
					float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})

		case mu.MU_COMMAND_ICON:
			switch cmd.Icon.Id {
			case mu.MU_ICON_CLOSE:
				clr := cmd.Icon.Color.ToRGBA()
				app.PushRect(
					float32(cmd.Icon.Rect.X),
					float32(cmd.Icon.Rect.Y),
					float32(cmd.Icon.Rect.W),
					float32(cmd.Icon.Rect.H),
					app.AtlasData.CloseIcon,
					[3]float32{
						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})

			case mu.MU_ICON_CHECK:
				clr := cmd.Icon.Color.ToRGBA()
				app.PushRect(
					float32(cmd.Icon.Rect.X),
					float32(cmd.Icon.Rect.Y),
					float32(cmd.Icon.Rect.W),
					float32(cmd.Icon.Rect.H),
					app.AtlasData.CheckedIcon,
					[3]float32{
						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
			case mu.MU_ICON_COLLAPSED:
				clr := cmd.Icon.Color.ToRGBA()
				app.PushRect(
					float32(cmd.Icon.Rect.X),
					float32(cmd.Icon.Rect.Y),
					float32(cmd.Icon.Rect.W),
					float32(cmd.Icon.Rect.H),
					app.AtlasData.CollapsedIcon,
					[3]float32{
						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
			case mu.MU_ICON_EXPANDED:
				clr := cmd.Icon.Color.ToRGBA()
				app.PushRect(
					float32(cmd.Icon.Rect.X),
					float32(cmd.Icon.Rect.Y),
					float32(cmd.Icon.Rect.W),
					float32(cmd.Icon.Rect.H),
					app.AtlasData.ExpandedIcon,
					[3]float32{
						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
			case mu.MU_ICON_MAX:
				clr := cmd.Icon.Color.ToRGBA()
				app.PushRect(
					float32(cmd.Icon.Rect.X),
					float32(cmd.Icon.Rect.Y),
					float32(cmd.Icon.Rect.W),
					float32(cmd.Icon.Rect.H),
					app.AtlasData.Black,
					[3]float32{
						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
			}
		}

	}

	DrawMyStuff(app)
}

type GlMeshData struct {
	Vertices    []float32
	Indices     []uint32
	VAO         uint32
	VBO         uint32
	IndexBuffer uint32
}

func NewGlMeshData() *GlMeshData {
	return &GlMeshData{}
}

type glTexture struct {
	TextureID uint32
	Width     int32
	Height    int32
}

type Point struct {
	X, Y float32
}

type Rect struct {
	P1, P2 Point
}

func PrepareGLobalState(app *App) {
	proj := mgl.Ortho2D(0, float32(app.Width), float32(app.Height), 0)
	gl.ClearColor(0.65, 0.65, 0.65, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, app.AtlasTexture.ID)

	gl.BindVertexArray(app.MeshBuffer.VAO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, app.MeshBuffer.IndexBuffer)
	gl.UseProgram(app.MainShader)

	loc := gl.GetUniformLocation(app.MainShader, gl.Str("uTexture\x00"))
	gl.Uniform1i(loc, 0)

	loc = gl.GetUniformLocation(app.MainShader, gl.Str("uProj\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &proj[0])
}

// DrawMyStuff draws my stuff
func DrawMyStuff(app *App) {
	app.FlushRects()
	gl.DrawElements(gl.TRIANGLES, int32(len(app.MeshBuffer.Indices)), gl.UNSIGNED_INT, nil)
}

func loadShaderSource(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	source := string(data)
	source = source + "\x00"
	return source, err
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func generateShader(vertexSRCPath, fragmentSRCPath string) uint32 {
	vertexSRC, err := loadShaderSource(vertexSRCPath)
	if err != nil {
		log.Println(err)
	}

	vertexShader, err := compileShader(vertexSRC, gl.VERTEX_SHADER)
	if err != nil {
		log.Println(err)
	}

	fragmentSRC, err := loadShaderSource(fragmentSRCPath)
	if err != nil {
		log.Println(err)
	}
	fragmentShader, err := compileShader(fragmentSRC, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Println(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))
		panic("LINK ERROR: " + log)
	}

	return prog
}
