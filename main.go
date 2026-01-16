package main

import (
	"fmt"
	gui2onegl "font-stuff/pkg/gui2one-gl"
	"os"

	"runtime"

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
    gui2onegl.GenerateAtlas("assets/fonts/ADRIP1.TTF")

    fmt.Println("OK")
    os.Exit(0) 

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
        // gl.Viewport(0,0,640, 480)
        gui2onegl.DrawMyStuff(&myApp)
        wnd.SwapBuffers()
        glfw.WaitEvents()
    }

    glfw.Terminate()


}