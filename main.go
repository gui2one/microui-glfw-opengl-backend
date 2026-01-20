package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"
	"math/rand"
	"path"

	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/zeozeozeo/microui-go"
)

var myApp gui2onegl.App
var Width = 512
var Height = 512
var MuCtx *microui.Context

func Render(ctx *microui.Context) {
	myApp.ClearRects()
	for _, cmd := range ctx.CommandList {
		switch cmd.Type {
		case microui.MU_COMMAND_RECT:
			// fmt.Println(cmd.Rect.Rect.X)
			myApp.PushRect(float32(cmd.Rect.Rect.X), float32(cmd.Rect.Rect.Y), float32(cmd.Rect.Rect.W), float32(cmd.Rect.Rect.H),
				myApp.AtlasData.White,
				[3]float32{0.5, 0.5, 1.0},
			)

		case microui.MU_COMMAND_TEXT:
			// fmt.Println(cmd.Rect.Rect.X)
			myApp.PushText(float32(cmd.Text.Pos.X), float32(cmd.Text.Pos.Y), cmd.Text.Str, [3]float32{1.0, 1.0, 1.0})

		case microui.MU_COMMAND_CLIP:
			// fmt.Println("clip", cmd.Clip.Rect.X)

		}
		// cmd = ctx.NextCommand(cmd)
	}

	gui2onegl.DrawMyStuff(&myApp, Width, Height)
}
func TextWidth(font microui.Font, text string) int {
	return 150
}
func TextHeight(font microui.Font) int {
	return 40
}
func initMyStuff() {

	fmt.Println("Init App OpenGL Resources")

	myApp.Init()

	// myApp.PushRect(10, 10, 512-20, 256-10,
	// 	myApp.AtlasData.White,
	// 	[3]float32{0.5, 0.2, 0.0},
	// )

	myApp.PushText(100, 256, "0123456789\ndefghijklmnopqrstuvwxyz/\n/;;\n;)", [3]float32{1.0, 1.0, 1.0})

}
func handleDrop(wnd *glfw.Window, paths []string) {
	fmt.Println("Dropped", len(paths), "files")
	fmt.Println(paths)
	fmt.Println(myApp.AtlasTexture.Width)
	first := paths[0]
	if path.Ext(first) == ".ttf" || path.Ext(first) == ".TTF" {
		atlas := gui2onegl.GenerateAtlas(first, [2]int{0x0020, 0x007E})
		gl.DeleteTextures(1, &myApp.AtlasTexture.ID)
		myApp.AtlasTexture = *gui2onegl.FromImage(atlas.Atlas)
	}
}
func handleResize(wnd *glfw.Window, width, height int) {
	Width = width
	Height = height
	gl.Viewport(0, 0, int32(Width), int32(Height))
}
func handleCursorPos(wnd *glfw.Window, x, y float64) {
	action := wnd.GetMouseButton(glfw.MouseButton1)
	if action == glfw.Press {
		myApp.PushRect(
			float32(x),
			float32(float64(Height)-y),
			rand.Float32()*100, rand.Float32()*100,
			myApp.AtlasData.White,
			[3]float32{rand.Float32(), rand.Float32(), rand.Float32()},
		)
	}
}
func handleMouseButton(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		x, y := wnd.GetCursorPos()

		myApp.PushRect(
			float32(x),
			float32(float64(Height)-y),
			rand.Float32()*100, rand.Float32()*100,
			myApp.AtlasData.White,
			[3]float32{rand.Float32(), rand.Float32(), rand.Float32()},
		)
	}
}
func main() {

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
	wnd.SetDropCallback(handleDrop)
	wnd.SetFramebufferSizeCallback(handleResize)
	wnd.SetCursorPosCallback(handleCursorPos)
	wnd.SetMouseButtonCallback(handleMouseButton)

	wnd.MakeContextCurrent()

	gui2onegl.InitGL()
	MuCtx = microui.NewContext()
	MuCtx.TextHeight = TextHeight
	MuCtx.TextWidth = TextWidth
	initMyStuff()
	gl.Viewport(0, 0, int32(Width), int32(Height))
	glfw.SwapInterval(0)
	for !wnd.ShouldClose() {

		glfw.WaitEvents()

		MuCtx.Begin()
		MuCtx.BeginWindow("window 1", microui.NewRect(100, 100, 256, 30))
		// MuCtx.Label("hello there !")
		MuCtx.EndWindow()
		MuCtx.End()

		Render(MuCtx)

		wnd.SwapBuffers()

	}

	glfw.Terminate()

}
