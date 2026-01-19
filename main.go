package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"
	"math/rand"
	"path"

	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var myApp gui2onegl.App
var Width = 640
var Height = 480

func initMyStuff() {

	fmt.Println("Init App OpenGL Resources")

	myApp.Init()

	myApp.PushRect(0.1, 0.1, 0.3, 0.3,
		gui2onegl.Rect{
			P1: gui2onegl.Point{X: 0.0, Y: 0.0},
			P2: gui2onegl.Point{X: 1.0, Y: 1.0},
		},
		[3]float32{1.0, 1.0, 1.0},
	)

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
			float32(x)/float32(Width)*(float32(Width)/float32(Height)),
			(float32(Height)-float32(y))/float32(Height),
			0.1, 0.1,
			gui2onegl.Rect{
				P1: gui2onegl.Point{X: 0.0, Y: 0.0},
				P2: gui2onegl.Point{X: 1.0, Y: 1.0},
			},
			[3]float32{rand.Float32(), rand.Float32(), rand.Float32()},
		)
	}
}
func handleMouseButton(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		x, y := wnd.GetCursorPos()

		myApp.PushRect(
			float32(x)/float32(Width)*(float32(Width)/float32(Height)),
			(float32(Height)-float32(y))/float32(Height),
			0.1, 0.1,
			gui2onegl.Rect{
				P1: gui2onegl.Point{X: 0.0, Y: 0.0},
				P2: gui2onegl.Point{X: 1.0, Y: 1.0},
			},
			[3]float32{rand.Float32(), rand.Float32(), rand.Float32()},
		)
	}
}
func main() {

	runtime.LockOSThread()

	if glfw.Init() != nil {
		panic("Unable to initialize GLFW")
	}
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

	initMyStuff()
	gl.Viewport(0, 0, int32(Width), int32(Height))
	glfw.SwapInterval(0)
	for !wnd.ShouldClose() {
		glfw.WaitEvents()
		gui2onegl.DrawMyStuff(&myApp, Width, Height)
		wnd.SwapBuffers()

	}

	glfw.Terminate()

}
