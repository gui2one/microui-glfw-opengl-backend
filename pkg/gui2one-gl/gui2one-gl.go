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
	Square       *GlMeshData
	MainShader   uint32
	AtlasTexture Texture
}

func (a *App) Init() {

	a.MainShader = generateShader("assets/shaders/main_vertex.glsl", "assets/shaders/main_fragment.glsl")
}

var Square = GlMeshData{
	Vertices: []float32{
		/*pos */ 0.0, 0.0 /*uvs */, 0.0, 0.0 /* color */, 1.0, 1.0, 1.0,
		/*pos */ 1.0, 0.0 /*uvs */, 1.0, 0.0 /* color */, 1.0, 1.0, 1.0,
		/*pos */ 1.0, 1.0 /*uvs */, 1.0, 1.0 /* color */, 1.0, 1.0, 1.0,
		/*pos */ 0.0, 1.0 /*uvs */, 0.0, 1.0 /* color */, 1.0, 1.0, 1.0,
	},
	Indices: []uint32{
		0, 1, 2,
		2, 3, 0,
	},
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

// Init GlMeshData gl resources
func (m *GlMeshData) Init() {
	sizeOfFloat32 := 4
	numFloatsPerVertex := 2 + 2 + 3
	stride := int32(sizeOfFloat32 * numFloatsPerVertex)

	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.IndexBuffer)
	gl.BindVertexArray(m.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, numFloatsPerVertex*len(m.Vertices), gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.IndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, numFloatsPerVertex*len(m.Indices), gl.Ptr(m.Indices), gl.STATIC_DRAW)

	/* position 2d */
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, stride, 0)
	/* uvs */
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(2*sizeOfFloat32))
	/* color RGB */
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, stride, uintptr((2+2)*sizeOfFloat32))

}

type glTexture struct {
	TextureID uint32
	Width     int32
	Height    int32
}

// DrawMyStuff draws my stuff
func DrawMyStuff(app *App, w, h int) {

	proj := mgl.Ortho2D(0, float32(w)/float32(h), 0, 1.0)
	gl.ClearColor(0.1, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, app.AtlasTexture.ID)

	gl.BindVertexArray(app.Square.VAO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, app.Square.IndexBuffer)
	gl.UseProgram(app.MainShader)

	loc := gl.GetUniformLocation(app.MainShader, gl.Str("uTexture\x00"))
	gl.Uniform1i(loc, 0)

	loc = gl.GetUniformLocation(app.MainShader, gl.Str("uProj\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &proj[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(app.Square.Indices)), gl.UNSIGNED_INT, nil)

	// gl.UseProgram(0)

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
