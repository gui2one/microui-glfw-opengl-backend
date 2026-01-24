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
	mu "github.com/zeozeozeo/microui-go"
)

var myApp muGL.App
var Width = 1280
var Height = 600
var MuCtx *microui.Context
var Val1 float32 = 5
var Text1 string = "text variable"
var Bool1 bool = true

type AppWindow struct {
	Name string
	Draw func()
}

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
	muEvents.SetCursorPosCallback(MuCtx, x, y)

	action := wnd.GetMouseButton(glfw.MouseButton1)
	if action == glfw.Press {

	}
}
func handleGLFWMouseButton(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	muEvents.SetMouseButtonCallback(MuCtx, wnd, button, action, mods)
}
func handleKeyDown(key int) {
	fmt.Println(key)
}
func handleGLFWKey(wnd *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	muEvents.SetKeyCallback(MuCtx, key, scancode, action, mods)
}
func handleGLFWChar(wnd *glfw.Window, char rune) {
	muEvents.SetCharCallBack(MuCtx, char)
}
func handleGLFWScroll(_ *glfw.Window, x, y float64) {
	muEvents.SetScrollCallback(MuCtx, x, y)
}

func MainWindow() {

	muGL.SliderWithLabel(MuCtx, "Slider", &Val1, 0.0, 10.0)

	MuCtx.LayoutRow(1, []int{-1}, 0)
	MuCtx.Text("Ici ... du texte Ici ... du texte Ici ... du texte Ici ... du texte Ici ... du texte")
	MuCtx.TextBox(&Text1)
	MuCtx.Checkbox("Bool Value", &Bool1)

	if MuCtx.Header("Collapsible Header") {
		MuCtx.Text("encore du texte")
		MuCtx.PushID([]byte("id1"))
		MuCtx.TextBox(&Text1)
		MuCtx.PopID()
	}
}
func OptionsWindow() {
	MuCtx.LayoutRow(1, []int{-1}, 0)
	for i := 0; i < 10; i++ {
		MuCtx.PushID([]byte{byte(i)})
		MuCtx.Slider(&Val1, 0.0, 10.0)
		MuCtx.PopID()
	}

}

// app window utils
func moveToFront(name string, windows []AppWindow) []AppWindow {
	for i, w := range windows {
		if w.Name == name {
			// Remove from current position
			windows = append(windows[:i], windows[i+1:]...)
			// Add to the end (top)
			return append(windows, w)
		}
	}
	return windows
}
func main() {
	Windows := []AppWindow{
		{
			Name: "Main",
			Draw: MainWindow,
		},
		{
			Name: "Options",
			Draw: OptionsWindow,
		},
		{
			Name: "Options2",
			Draw: OptionsWindow,
		},
	}
	runtime.LockOSThread()

	if glfw.Init() != nil {
		panic("Unable to initialize GLFW")
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	wnd, err := glfw.CreateWindow(Width, Height, "gui2one | GL engine | another one ? ... he should stop ", nil, nil)
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

	MuCtx = microui.NewContext()

	muGL.InitGL()
	myApp.InitGL(Width, Height)

	myApp.InitMuContext(MuCtx)

	gl.Viewport(0, 0, int32(myApp.Width), int32(myApp.Height))
	glfw.SwapInterval(0)

	fmt.Println("MicroUI GL Backend")
	for !wnd.ShouldClose() {
		windowToMove := ""
		glfw.WaitEvents()

		MuCtx.Begin()

		for i, w := range Windows {
			if MuCtx.BeginWindow(w.Name, mu.NewRect((i+1)*150, (i+1)*150, 300, 650)) {
				container := MuCtx.GetCurrentContainer()

				if MuCtx.MousePressed == microui.MU_MOUSE_LEFT && MuCtx.HoverRoot == container {
					windowToMove = w.Name
				}
				w.Draw()
				MuCtx.EndWindow()
			}
		}

		MuCtx.End()

		gl.ClearColor(0.5, 0.1, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		myApp.Render(MuCtx)

		wnd.SwapBuffers()

		if windowToMove != "" {
			Windows = moveToFront(windowToMove, Windows)
			windowToMove = ""
			glfw.PostEmptyEvent()
		}
	}

	glfw.Terminate()

}
