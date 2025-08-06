package matcha

import (
	"github.com/cchirag/matcha/style"
)

type Color16 style.Color16

func Color(hex24 string, color256 uint8, color16 style.Color16) style.Color {
	return style.NewColor(hex24, color256, color16)
}

// matcha:import color16Enum
// Code injected by Matcha generator. DO NOT EDIT.
// matcha:inject-begin color16Enum
const (
	ColorBlack Color16 = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// matcha:inject-end color16Enum
