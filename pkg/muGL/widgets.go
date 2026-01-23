package muGL

import (
	mu "github.com/zeozeozeo/microui-go"
)

func SliderWithLabel(ctx *mu.Context, labelText string, val *float32, min float32, max float32) {
	ctx.LayoutRow(2, []int{100, -1}, 0)
	ctx.Label(labelText)
	ctx.Slider(val, min, max)
}
