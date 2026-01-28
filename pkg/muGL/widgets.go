package muGL

import (
	"fmt"

	mu "github.com/zeozeozeo/microui-go"
)

func SliderWithLabel(ctx *mu.Context, labelText string, val *float32, min float32, max float32) {
	ctx.LayoutRow(2, []int{100, -1}, 0)
	ctx.Label(labelText)
	ctx.Slider(val, min, max)
}

func (app *App) MenuBar() {
	ctx := app.CTX
	opts := mu.MU_OPT_NOCLOSE | mu.MU_OPT_NOTITLE | mu.MU_OPT_NORESIZE | mu.MU_OPT_NOSCROLL
	if ctx.BeginWindowEx("menu_bar", mu.NewRect(0, 0, app.Width, 30), opts) != 0 {

		ctx.LayoutRow(len(app.Windows), []int{100}, 30)
		for i := range app.Windows {
			container := ctx.GetContainer(app.Windows[i].Name)

			var opened = false
			ctx.PushID([]byte(fmt.Sprintf("%v__checkbox", app.Windows[i].Name)))
			if ctx.Checkbox(app.Windows[i].Name, &opened) != 0 {
				container.Open = !container.Open
				fmt.Println(container.Open)
			}
			ctx.PopID()
		}
		// ctx.Checkbox("0", &opened)
		ctx.EndWindow()
	}

}
