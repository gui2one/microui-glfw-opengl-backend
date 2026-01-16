package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)


var myApp gui2onegl.App



func initMyStuff(){


        fmt.Println("hey")
        
        myApp.Init()
        myApp.Square =  &gui2onegl.Square
        myApp.Square.Init()

}


func main() {

    runtime.LockOSThread()
    if glfw.Init() != nil {
        panic("Unable to initialize GLFW")
    }

    wnd, err := glfw.CreateWindow(640, 480, "Hello World", nil, nil)
    if err != nil {
        panic("Unable to create GLFW window")
    }

    wnd.MakeContextCurrent()
    if gl.Init() != nil {
        panic("Unable to initialize OpenGL")
    }
    gui2onegl.SetupGLDebug()

    initMyStuff()
    
    for !wnd.ShouldClose() {
        // gl.Viewport(0,0,640, 480)
        gui2onegl.DrawMyStuff(&myApp)
        wnd.SwapBuffers()
        glfw.PollEvents()
    }

    glfw.Terminate()


}