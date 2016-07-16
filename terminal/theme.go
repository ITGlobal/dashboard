package main

import "github.com/nsf/termbox-go"

type Chars struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune

	VerticalAndRight rune
	VerticalAndLeft  rune

	VLine       rune
	HLine       rune
	HLineSingle rune

	FullBlock  rune
	LightShade rune
	Ellipsis   rune
}

type Colors struct {
	DefaultFg termbox.Attribute
	DefaultBg termbox.Attribute

	GoodFg termbox.Attribute
	GoodBg termbox.Attribute

	BadFg termbox.Attribute
	BadBg termbox.Attribute

	PendingFg termbox.Attribute
	PendingBg termbox.Attribute
}

type ThemeData struct {
	Chars  Chars
	Colors Colors
}

var Theme ThemeData = ThemeData{singleThemeChars, lightThemeColors}

type ThemeKey int

const (
	ThemeDefault ThemeKey = 0

	ThemeSingle ThemeKey = 0x01
	ThemeDouble ThemeKey = 0x02

	ThemeLight ThemeKey = 0x04
	ThemeDark  ThemeKey = 0x08
)

func SetTheme(key ThemeKey) {
	if key&ThemeSingle == ThemeSingle {
		Theme.Chars = singleThemeChars
	}

	if key&ThemeDouble == ThemeDouble {
		Theme.Chars = doubleThemeChars
	}

	if key&ThemeLight == ThemeLight {
		Theme.Colors = lightThemeColors
	}

	if key&ThemeDark == ThemeDark {
		Theme.Colors = darkThemeColors
	}
}

var singleThemeChars = Chars{
	TopLeft:     rune(0x250C),
	TopRight:    rune(0x2510),
	BottomLeft:  rune(0x2514),
	BottomRight: rune(0x2518),

	VerticalAndRight: rune(0x251C),
	VerticalAndLeft:  rune(0x2524),

	VLine:       rune(0x2502),
	HLine:       rune(0x2500),
	HLineSingle: rune(0x2500),

	FullBlock:  rune(0x2588),
	LightShade: rune(0x2591),
	Ellipsis:   rune(0x2026),
}

var doubleThemeChars = Chars{
	TopLeft:     rune(0x2554),
	TopRight:    rune(0x2557),
	BottomLeft:  rune(0x255A),
	BottomRight: rune(0x255D),

	VerticalAndRight: rune(0x255F),
	VerticalAndLeft:  rune(0x2562),

	VLine:       rune(0x2551),
	HLine:       rune(0x2550),
	HLineSingle: rune(0x2500),

	FullBlock:  rune(0x2588),
	LightShade: rune(0x2591),
	Ellipsis:   rune(0x2026),
}

var lightThemeColors = Colors{
	DefaultFg: termbox.ColorBlack,
	DefaultBg: termbox.ColorWhite | termbox.AttrBold,

	GoodFg: termbox.ColorBlack,
	GoodBg: termbox.ColorGreen | termbox.AttrBold,

	BadFg: termbox.ColorBlack,
	BadBg: termbox.ColorRed | termbox.AttrBold,

	PendingFg: termbox.ColorBlack,
	PendingBg: termbox.ColorCyan | termbox.AttrBold,
}

var darkThemeColors = Colors{
	DefaultFg: termbox.ColorWhite | termbox.AttrBold,
	DefaultBg: termbox.ColorBlack,

	GoodFg: termbox.ColorGreen | termbox.AttrBold,
	GoodBg: termbox.ColorBlack,

	BadFg: termbox.ColorRed | termbox.AttrBold,
	BadBg: termbox.ColorBlack,

	PendingFg: termbox.ColorCyan | termbox.AttrBold,
	PendingBg: termbox.ColorBlack,
}
