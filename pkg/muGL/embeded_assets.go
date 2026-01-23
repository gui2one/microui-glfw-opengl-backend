package muGL

import (
	"embed"
	_ "embed"
)

//go:embed assets/shaders/*
var STATIC_SHADERS embed.FS
