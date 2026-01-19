package gui2onegl

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

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
	AtlasTexture       Texture
	NumFloatsPerVertex int
	SizeOfFloat32      int
}

func (a *App) Init() {
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

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.IndexBuffer)

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

	//fmt.Println("Pushed Rect", len(m.Vertices)/numFloatsPerVertex)
}
func (a *App) FlushRects() {
	m := a.MeshBuffer

	gl.BindVertexArray(m.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, a.SizeOfFloat32*len(m.Vertices), gl.Ptr(m.Vertices), gl.STREAM_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.IndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, a.SizeOfFloat32*len(m.Indices), gl.Ptr(m.Indices), gl.STREAM_DRAW)

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

// DrawMyStuff draws my stuff
func DrawMyStuff(app *App, w, h int) {
	app.FlushRects()
	proj := mgl.Ortho2D(0, float32(w)/float32(h), 0, 1.0)
	gl.ClearColor(0.1, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, app.AtlasTexture.ID)

	gl.BindVertexArray(app.MeshBuffer.VAO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, app.MeshBuffer.IndexBuffer)
	gl.UseProgram(app.MainShader)

	loc := gl.GetUniformLocation(app.MainShader, gl.Str("uTexture\x00"))
	gl.Uniform1i(loc, 0)

	loc = gl.GetUniformLocation(app.MainShader, gl.Str("uProj\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &proj[0])
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
