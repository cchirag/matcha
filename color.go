package matcha

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Color16 uint8

// matcha:export color16Enum
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

// matcha:end

type Color struct {
	hex24    string
	color256 uint8
	color16  Color16
}

func NewColor(hex24 string, color256 uint8, color16 Color16) Color {
	if !strings.HasPrefix(hex24, "#") || len(hex24) != 7 {
		panic("matcha: hex color must be in the form #RRGGBB")
	}
	if color256 > 255 {
		panic("matcha: 256-color value must be between 0–255")
	}
	if color16 > ColorBrightWhite {
		panic("matcha: invalid 16-color enum value")
	}
	return Color{hex24: hex24, color256: color256, color16: color16}
}

func (c Color) resolve(screen tcell.Screen) tcell.Color {
	var col tcell.Color

	switch {
	case screen.Colors() >= 1<<24:
		val64, _ := strconv.ParseUint(strings.TrimPrefix(c.hex24, "#"), 16, 32)
		col = tcell.NewHexColor(int32(val64))
	case screen.Colors() >= 256:
		col = tcell.Color(c.color256)
	default:
		col = mapColor16ToTCell(c.color16)
	}
	return col
}

func mapColor16ToTCell(c Color16) tcell.Color {
	switch c {
	case ColorBlack:
		return tcell.ColorBlack
	case ColorRed:
		return tcell.ColorMaroon
	case ColorGreen:
		return tcell.ColorGreen
	case ColorYellow:
		return tcell.ColorOlive
	case ColorBlue:
		return tcell.ColorNavy
	case ColorMagenta:
		return tcell.ColorPurple
	case ColorCyan:
		return tcell.ColorTeal
	case ColorWhite:
		return tcell.ColorSilver
	case ColorBrightBlack:
		return tcell.ColorGray
	case ColorBrightRed:
		return tcell.ColorRed
	case ColorBrightGreen:
		return tcell.ColorLime
	case ColorBrightYellow:
		return tcell.ColorYellow
	case ColorBrightBlue:
		return tcell.ColorBlue
	case ColorBrightMagenta:
		return tcell.ColorFuchsia
	case ColorBrightCyan:
		return tcell.ColorAqua
	case ColorBrightWhite:
		return tcell.ColorWhite
	default:
		panic("matcha: unknown Color16 value")
	}
}
