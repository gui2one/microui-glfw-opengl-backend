package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"

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
	myApp.Square = &gui2onegl.Square
	myApp.Square.Init()

}

func handleDrop(wnd *glfw.Window, paths []string) {
	fmt.Println("Dropped", len(paths), "files")
	fmt.Println(paths)
	fmt.Println(myApp.AtlasTexture.Width)
}
func handleResize(wnd *glfw.Window, width, height int) {
	Width = width
	Height = height
	gl.Viewport(0, 0, int32(Width), int32(Height))
}
func main() {

	runtime.LockOSThread()
	fmt.Println("Starting App...")
	atlas := gui2onegl.GenerateAtlas("assets/fonts/CONSOLAB.TTF", [2]int{0x0020, 0x007E})
	// atlas := gui2onegl.GenerateAtlas("assets/fonts/ARIAL.TTF", [2]int{0x0020, 0x0023})

	atlas.Print(false)

	if glfw.Init() != nil {
		panic("Unable to initialize GLFW")
	}
	wnd, err := glfw.CreateWindow(Width, Height, "Hello World", nil, nil)
	wnd.SetDropCallback(handleDrop)
	wnd.SetFramebufferSizeCallback(handleResize)
	if err != nil {
		panic("Unable to create GLFW window")
	}

	wnd.MakeContextCurrent()

	gui2onegl.InitGL()
	myApp.AtlasTexture = *gui2onegl.FromImage(atlas.Atlas)
	initMyStuff()
	gl.Viewport(0, 0, int32(Width), int32(Height))
	for !wnd.ShouldClose() {

		gui2onegl.DrawMyStuff(&myApp, Width, Height)
		wnd.SwapBuffers()
		glfw.WaitEvents()
	}

	glfw.Terminate()

}
