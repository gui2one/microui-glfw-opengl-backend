package main

import (
	"fmt"
	AG "font-stuff/pkg/atlas_gen"
	muGL "font-stuff/pkg/muGL"
	"path"

	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/zeozeozeo/microui-go"
	mu "github.com/zeozeozeo/microui-go"
)

var myApp muGL.App
var Width = 1280
var Height = 720
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
	MuCtx.InputMouseMove(int(x), int(y))

	action := wnd.GetMouseButton(glfw.MouseButton1)
	if action == glfw.Press {

	}
}
func handleGLFWMouseButton(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {

	// Map GLFW buttons to MicroUI buttons
	var muBtn int
	switch button {
	case glfw.MouseButtonLeft:
		muBtn = microui.MU_MOUSE_LEFT
	case glfw.MouseButtonRight:
		muBtn = microui.MU_MOUSE_RIGHT
	case glfw.MouseButtonMiddle:
		muBtn = microui.MU_MOUSE_MIDDLE
	default:
		return
	}
	switch action {
	case glfw.Release:
		x, y := wnd.GetCursorPos()
		MuCtx.InputMouseUp(int(x), int(y), muBtn)
	case glfw.Press:
		x, y := wnd.GetCursorPos()
		MuCtx.InputMouseDown(int(x), int(y), muBtn)

	}

}
func handleKeyDown(key int) {
	fmt.Println(key)
}
func handleGLFWKey(wnd *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press, glfw.Repeat:
		switch key {
		case glfw.KeyBackspace:
			MuCtx.InputKeyDown(microui.MU_KEY_BACKSPACE)
		case glfw.KeyEnter:
			MuCtx.InputKeyDown(microui.MU_KEY_RETURN)
			// Add other functional keys as needed
		}
	case glfw.Release:
		switch key {
		case glfw.KeyBackspace:
			MuCtx.InputKeyUp(microui.MU_KEY_BACKSPACE)
		case glfw.KeyEnter:
			MuCtx.InputKeyUp(microui.MU_KEY_RETURN)
		}
	}

}
func handleGLFWChar(wnd *glfw.Window, char rune) {
	MuCtx.InputText([]rune{char})
}
func handleGLFWScroll(_ *glfw.Window, x, y float64) {
	MuCtx.InputScroll(int(x*10), int(-y*10))
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
	myApp.InitGL()

	myApp.InitMuContext(MuCtx)

	myApp.Width = Width
	myApp.Height = Height
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
