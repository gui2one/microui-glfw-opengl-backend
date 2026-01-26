package main

import (
	"fmt"
	"path"

	AG "github.com/gui2one/microui-glfw-opengl-backend/pkg/atlas_gen"
	muGL "github.com/gui2one/microui-glfw-opengl-backend/pkg/muGL"
	muEvents "github.com/gui2one/microui-glfw-opengl-backend/pkg/muGL/glfw"

	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/zeozeozeo/microui-go"
)

var myApp muGL.App

var ShowMainWindow = true
var Width = 1280
var Height = 720
var muCtx *microui.Context
var val1 float32 = 5
var text1 string = "text variable"
var bool1 bool = true

func handleGLFWDrop(wnd *glfw.Window, paths []string) {
	fmt.Println("Dropped", len(paths), "files")
	fmt.Println(paths)
	fmt.Println(myApp.AtlasTexture.Width)
	first := paths[0]
	if path.Ext(first) == ".ttf" || path.Ext(first) == ".TTF" {
		atlas := AG.GenerateAtlas(first, muGL.GLYPHS_RANGE, 18)
		gl.DeleteTextures(1, &myApp.AtlasTexture.ID)
		myApp.AtlasTexture = *muGL.FromImage(atlas.Atlas)
	}
}
func handleGLFWResize(wnd *glfw.Window, width, height int) {
	myApp.Width = width
	myApp.Height = height
	gl.Viewport(0, 0, int32(myApp.Width), int32(myApp.Height))
}
func handleGLFWCursorPos(wnd *glfw.Window, x, y float64) {
	muEvents.SetCursorPosCallback(muCtx, x, y)

	action := wnd.GetMouseButton(glfw.MouseButton1)
	if action == glfw.Press {

	}
}
func handleGLFWMouseButton(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	muEvents.SetMouseButtonCallback(muCtx, wnd, button, action, mods)
}
func handleKeyDown(key int) {
	fmt.Println(key)
}
func handleGLFWKey(wnd *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	muEvents.SetKeyCallback(muCtx, key, scancode, action, mods)
	if action == glfw.Press {
		switch key {
		case glfw.KeySpace:

			for i := range myApp.Windows {
				w := &myApp.Windows[i]
				w.Closed = false
				fmt.Println(w.Closed)
			}
		}
	}
}
func handleGLFWChar(wnd *glfw.Window, char rune) {
	muEvents.SetCharCallBack(muCtx, char)
}
func handleGLFWScroll(_ *glfw.Window, x, y float64) {
	muEvents.SetScrollCallback(muCtx, x, y)
}

func InitGL() {

	if gl.Init() != nil {
		panic("Unable to initialize OpenGL")
	}
	muGL.SetupGLDebug()
}

func MainWindow() {

	muGL.SliderWithLabel(muCtx, "Slider", &val1, 0.0, 10.0)

	muCtx.LayoutRow(1, []int{-1}, 0)
	muCtx.Text("Ici ... du texte Ici ... du texte Ici ... du texte Ici ... du texte Ici ... du texte")
	muCtx.TextBox(&text1)
	muCtx.Checkbox("Bool Value", &bool1)

	if muCtx.Header("Collapsible Header") {
		muCtx.Text("encore du texte")
		muCtx.PushID([]byte("id1"))
		muCtx.TextBox(&text1)
		muCtx.PopID()
	}
}
func OptionsWindow() {
	muCtx.LayoutRow(1, []int{-1}, 0)
	for i := 0; i < 10; i++ {
		muCtx.PushID([]byte{byte(i)})
		muCtx.Slider(&val1, 0.0, 10.0)
		muCtx.PopID()
	}

}

func main() {
	myApp.Windows = []muGL.Window{
		{
			Name:   "Main",
			Draw:   MainWindow,
			X:      0,
			Y:      0,
			Width:  300,
			Height: 400,
		},
		{
			Name:   "Options",
			Draw:   OptionsWindow,
			X:      100,
			Y:      50,
			Width:  300,
			Height: 400,
		},
		{
			Name:   "Options2",
			Draw:   OptionsWindow,
			X:      200,
			Y:      100,
			Width:  300,
			Height: 400,
		},
	}
	runtime.LockOSThread()

	if glfw.Init() != nil {
		panic("Unable to initialize GLFW")
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	wnd, err := glfw.CreateWindow(Width, Height, "muGL | a microui opengl backend", nil, nil)
	if err != nil {
		panic("Unable to create GLFW window")
	}
	wnd.SetDropCallback(handleGLFWDrop)
	wnd.SetFramebufferSizeCallback(handleGLFWResize)
	wnd.SetCursorPosCallback(handleGLFWCursorPos)
	wnd.SetMouseButtonCallback(handleGLFWMouseButton)
	wnd.SetKeyCallback(handleGLFWKey)
	wnd.SetCharCallback(handleGLFWChar)
	wnd.SetScrollCallback(handleGLFWScroll)

	// OpenGL Starts here !!
	wnd.MakeContextCurrent()

	muCtx = microui.NewContext()

	InitGL()
	myApp.InitGL(Width, Height)
	myApp.InitMuContext(muCtx)

	gl.Viewport(0, 0, int32(myApp.Width), int32(myApp.Height))
	glfw.SwapInterval(0)

	fmt.Println("MicroUI GL Backend")
	for !wnd.ShouldClose() {

		glfw.WaitEvents()
		myApp.CTX.Begin()

		myApp.PutWindows()

		myApp.CTX.End()

		if myApp.WindowToMove != "" {
			myApp.Windows = muGL.MoveToFront(myApp.WindowToMove, myApp.Windows)

			myApp.WindowToMove = ""
			glfw.PostEmptyEvent()
		}

		gl.ClearColor(0.5, 0.1, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		myApp.Render()

		wnd.SwapBuffers()

	}

	glfw.Terminate()

}
