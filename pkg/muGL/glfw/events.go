package muGL

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	mu "github.com/zeozeozeo/microui-go"
)

func SetScrollCallback(ctx *mu.Context, x, y float64) {
	ctx.InputScroll(int(x)*10, int(-y*10))
}
func SetCharCallBack(ctx *mu.Context, c rune) {
	ctx.InputText([]rune{c})
}

func SetKeyCallback(ctx *mu.Context, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press, glfw.Repeat:
		switch key {
		case glfw.KeyBackspace:
			ctx.InputKeyDown(mu.MU_KEY_BACKSPACE)
		case glfw.KeyEnter:
			ctx.InputKeyDown(mu.MU_KEY_RETURN)
			// Add other functional keys as needed
		}
	case glfw.Release:
		switch key {
		case glfw.KeyBackspace:
			ctx.InputKeyUp(mu.MU_KEY_BACKSPACE)
		case glfw.KeyEnter:
			ctx.InputKeyUp(mu.MU_KEY_RETURN)
		}
	}
}

func SetMouseButtonCallback(ctx *mu.Context, wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	// Map GLFW buttons to MicroUI buttons
	var muBtn int
	switch button {
	case glfw.MouseButtonLeft:
		muBtn = mu.MU_MOUSE_LEFT
	case glfw.MouseButtonRight:
		muBtn = mu.MU_MOUSE_RIGHT
	case glfw.MouseButtonMiddle:
		muBtn = mu.MU_MOUSE_MIDDLE
	default:
		return
	}
	switch action {
	case glfw.Release:
		x, y := wnd.GetCursorPos()
		ctx.InputMouseUp(int(x), int(y), muBtn)
	case glfw.Press:
		x, y := wnd.GetCursorPos()
		ctx.InputMouseDown(int(x), int(y), muBtn)

	}

}

func SetCursorPosCallback(ctx *mu.Context, x, y float64) {
	ctx.InputMouseMove(int(x), int(y))
}
