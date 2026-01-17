package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"

	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var myApp gui2onegl.App

func initMyStuff() {
	texture, err := gui2onegl.LoadImageFile("out.png")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Init App OpenGL Resources")
	myApp.Init()
	myApp.Square = &gui2onegl.Square
	myApp.Square.Init()
	myApp.AtlasTexture = *texture
}

func main() {

	runtime.LockOSThread()
	fmt.Println("Starting App...")
	atlas := gui2onegl.GenerateAtlas("assets/fonts/ARIAL.TTF", [2]int{0x0020, 0x007E})
	// atlas := gui2onegl.GenerateAtlas("assets/fonts/ARIAL.TTF", [2]int{0x0020, 0x0023})

	atlas.Print(false)

	if glfw.Init() != nil {
		panic("Unable to initialize GLFW")
	}

	wnd, err := glfw.CreateWindow(640, 480, "Hello World", nil, nil)
	if err != nil {
		panic("Unable to create GLFW window")
	}

	wnd.MakeContextCurrent()

	gui2onegl.InitGL()

	initMyStuff()

	for !wnd.ShouldClose() {

		gui2onegl.DrawMyStuff(&myApp)
		wnd.SwapBuffers()
		glfw.WaitEvents()
	}

	glfw.Terminate()

}
