package atlas_gen

import (
	_ "embed"
)

//go:embed assets/fonts/CONSOLAB.TTF
var defaultFont []byte

func GetFont() []byte {
	// You now have the font in memory, no matter where the app is running
	return defaultFont
}
