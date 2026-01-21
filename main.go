package main

import (
	"fmt"
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
var Text1 string = "0123456789"
var Bool1 bool = true

/* MicrUI "implementation" */
// func Render(ctx *microui.Context, app *muGL.App) {
// 	muGL.PrepareGLobalState(app)
// 	app.ClearRects()
// 	gl.Disable(gl.SCISSOR_TEST) // Start with no scissor
// 	for _, cmd := range ctx.CommandList {
// 		switch cmd.Type {
// 		case microui.MU_COMMAND_CLIP:
// 			muGL.DrawMyStuff(app)
// 			app.ClearRects()
// 			app.SetScissor(cmd.Clip.Rect)

// 		case microui.MU_COMMAND_RECT:

// 			rgba := cmd.Rect.Color.ToRGBA()
// 			app.PushRect(float32(cmd.Rect.Rect.X), float32(cmd.Rect.Rect.Y), float32(cmd.Rect.Rect.W), float32(cmd.Rect.Rect.H),
// 				app.AtlasData.White,
// 				[3]float32{float32(rgba.R) / 255.0, float32(rgba.G) / 255.0, float32(rgba.B) / 255.0},
// 			)

// 		case microui.MU_COMMAND_TEXT:

// 			clr := cmd.Text.Color.ToRGBA()
// 			app.PushText(
// 				float32(cmd.Text.Pos.X),
// 				float32(cmd.Text.Pos.Y),
// 				cmd.Text.Str,
// 				[3]float32{
// 					float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})

// 		case mu.MU_COMMAND_ICON:
// 			switch cmd.Icon.Id {
// 			case mu.MU_ICON_CLOSE:
// 				clr := cmd.Icon.Color.ToRGBA()
// 				app.PushRect(
// 					float32(cmd.Icon.Rect.X),
// 					float32(cmd.Icon.Rect.Y),
// 					float32(cmd.Icon.Rect.W),
// 					float32(cmd.Icon.Rect.H),
// 					app.AtlasData.CloseIcon,
// 					[3]float32{
// 						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})

// 			case mu.MU_ICON_CHECK:
// 				clr := cmd.Icon.Color.ToRGBA()
// 				app.PushRect(
// 					float32(cmd.Icon.Rect.X),
// 					float32(cmd.Icon.Rect.Y),
// 					float32(cmd.Icon.Rect.W),
// 					float32(cmd.Icon.Rect.H),
// 					app.AtlasData.CheckedIcon,
// 					[3]float32{
// 						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
// 			case mu.MU_ICON_COLLAPSED:
// 				clr := cmd.Icon.Color.ToRGBA()
// 				app.PushRect(
// 					float32(cmd.Icon.Rect.X),
// 					float32(cmd.Icon.Rect.Y),
// 					float32(cmd.Icon.Rect.W),
// 					float32(cmd.Icon.Rect.H),
// 					app.AtlasData.White,
// 					[3]float32{
// 						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
// 			case mu.MU_ICON_EXPANDED:
// 				clr := cmd.Icon.Color.ToRGBA()
// 				app.PushRect(
// 					float32(cmd.Icon.Rect.X),
// 					float32(cmd.Icon.Rect.Y),
// 					float32(cmd.Icon.Rect.W),
// 					float32(cmd.Icon.Rect.H),
// 					app.AtlasData.White,
// 					[3]float32{
// 						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
// 			case mu.MU_ICON_MAX:
// 				clr := cmd.Icon.Color.ToRGBA()
// 				app.PushRect(
// 					float32(cmd.Icon.Rect.X),
// 					float32(cmd.Icon.Rect.Y),
// 					float32(cmd.Icon.Rect.W),
// 					float32(cmd.Icon.Rect.H),
// 					app.AtlasData.Black,
// 					[3]float32{
// 						float32(clr.R) / 255.0, float32(clr.G) / 255.0, float32(clr.B) / 255.0})
// 			}
// 		}

// 	}

// 	muGL.DrawMyStuff(app)
// }

// func TextWidth(font microui.Font, text string) int {
// 	w := myApp.ComputeTextWidth(text)
// 	// fmt.Println("Width of ", text, " \nis ", w)
// 	return w
// }
// func TextHeight(font microui.Font) int {
// 	return myApp.AtlasData.FontMetrics.LineHeight
// }

func handleGLFWDrop(wnd *glfw.Window, paths []string) {
	fmt.Println("Dropped", len(paths), "files")
	fmt.Println(paths)
	fmt.Println(myApp.AtlasTexture.Width)
	first := paths[0]
	if path.Ext(first) == ".ttf" || path.Ext(first) == ".TTF" {
		atlas := muGL.GenerateAtlas(first, muGL.GLYPHS_RANGE, 18)
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
func handleGLFWScroll(wnd *glfw.Window, x, y float64) {
	MuCtx.InputScroll(int(x), int(y))
}
func MainWindow() {
	MuCtx.LayoutRow(1, []int{-1}, 0)
	MuCtx.Label("&&hello there!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	MuCtx.Slider(&Val1, 0.0, 10.0)
	MuCtx.Text("Ici ... du texte")
	MuCtx.TextBox(&Text1)
	MuCtx.Checkbox("Bool Value", &Bool1)
}
func OptionsWindow() {
	MuCtx.LayoutRow(1, []int{-1}, 0)
	MuCtx.Label("options !!")

	MuCtx.Text("Ici ... du texte")
	MuCtx.TextBox(&Text1)
}

type AppWindow struct {
	Name string
	Draw func()
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

	muGL.InitGL()
	MuCtx = microui.NewContext()

	myFontHandle := &myApp.AtlasData

	MuCtx.Style.Font = myFontHandle
	MuCtx.TextHeight = myApp.TextHeight
	MuCtx.TextWidth = myApp.TextWidth

	myApp.Init()
	myApp.Width = Width
	myApp.Height = Height
	gl.Viewport(0, 0, int32(myApp.Width), int32(myApp.Height))
	glfw.SwapInterval(0)

	for !wnd.ShouldClose() {
		windowToMove := ""
		glfw.WaitEvents()

		MuCtx.Begin()

		for i, w := range Windows {

			if MuCtx.BeginWindow(w.Name, mu.NewRect((i+1)*50, (i+1)*50, 200, 300)) {
				container := MuCtx.GetCurrentContainer()

				if MuCtx.MousePressed == microui.MU_MOUSE_LEFT && MuCtx.HoverRoot == container {
					windowToMove = w.Name
				}
				w.Draw()

				MuCtx.EndWindow()
			}
		}

		MuCtx.End()

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
