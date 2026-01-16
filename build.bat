@echo off
@REM set CGO_ENABLED=1
set CGO_LDFLAGS=-LD:\dev\libs\glfw_install\lib -lglfw3 -lopengl32 -lgdi32
set CGO_CFLAGS=-ID:\dev\libs\glfw_install\include
go build -v