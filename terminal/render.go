package main

import (
	"math"

	"github.com/nsf/termbox-go"
)

type ItemRenderer interface {
	Width() int
	Height() int

	Clear(fg, bg termbox.Attribute)

	SetFg(x, y int, attr termbox.Attribute)
	SetBg(x, y int, attr termbox.Attribute)
	SetCh(x, y int, ch rune)

	HLine(x, y, len int)
	HLineSingle(x, y, len int)
	VLine(x, y, len int)

	WriteText(x, y int, text string, maxWidth int) int
	WriteTextEx(x, y int, text string, maxWidth int, offset int, enableEllipsis bool) int

	ProgressBar(x, y, width, value int, fg, bg termbox.Attribute)

	Render()
}

type itemRenderer struct {
	x, y, w, h int
	cells      []termbox.Cell
}

func NewItemRenderer(x, y, w, h int) ItemRenderer {
	cells := make([]termbox.Cell, w*h)
	for i := range cells {
		cells[i] = termbox.Cell{Bg: termbox.ColorBlack, Fg: termbox.ColorWhite, Ch: ' '}
	}

	return &itemRenderer{x, y, w, h, cells}
}

func (r *itemRenderer) Width() int  { return r.w }
func (r *itemRenderer) Height() int { return r.h }

func (r *itemRenderer) Render() {
	for i := 0; i < r.w; i++ {
		x := i + r.x
		for j := 0; j < r.h; j++ {
			y := j + r.y

			idx := r.calcIndex(i, j)
			cell := r.cells[idx]
			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func (r *itemRenderer) Clear(fg, bg termbox.Attribute) {
	for i := range r.cells {
		r.cells[i].Fg = fg
		r.cells[i].Bg = bg
	}
}

func (r *itemRenderer) calcIndex(x, y int) int {
	i := r.w*y + x
	return i
}

func (r *itemRenderer) SetFg(x, y int, attr termbox.Attribute) {
	i := r.calcIndex(x, y)
	r.cells[i].Fg = attr
}

func (r *itemRenderer) SetBg(x, y int, attr termbox.Attribute) {
	i := r.calcIndex(x, y)
	r.cells[i].Bg = attr
}

func (r *itemRenderer) SetCh(x, y int, ch rune) {
	i := r.calcIndex(x, y)
	r.cells[i].Ch = ch
}

func (r *itemRenderer) HLine(x, y, len int) {
	for i := x; i < x+len; i++ {
		r.SetCh(i, y, Theme.Chars.HLine)
		r.SetCh(i, y, Theme.Chars.HLine)
	}
}

func (r *itemRenderer) HLineSingle(x, y, len int) {
	for i := x; i < x+len; i++ {
		r.SetCh(i, y, Theme.Chars.HLineSingle)
		r.SetCh(i, y, Theme.Chars.HLineSingle)
	}
}

func (r *itemRenderer) VLine(x, y, len int) {
	for i := y; i < y+len; i++ {
		r.SetCh(x, i, Theme.Chars.VLine)
		r.SetCh(x, i, Theme.Chars.VLine)
	}
}

func (r *itemRenderer) WriteText(x, y int, text string, maxWidth int) int {
	return r.WriteTextEx(x, y, text, maxWidth, 0, true)
}

func (r *itemRenderer) WriteTextEx(x, y int, text string, maxWidth int, offset int, enableEllipsis bool) int {
	var i int
	j := 0
	for i = offset; i < len(text); i++ {
		if j == maxWidth-1 {
			if enableEllipsis && i < len(text)-1 {
				r.SetCh(x+j, y, Theme.Chars.Ellipsis)
			} else {
				r.SetCh(x+j, y, rune(text[i]))
			}
			break
		}

		r.SetCh(x+j, y, rune(text[i]))
		j++
	}

	return i
}

func (r *itemRenderer) ProgressBar(x, y, width, value int, fg, bg termbox.Attribute) {
	widthToFill := int(math.Ceil(float64(width) * float64(value) / 100.0))
	i := 0

	for ; i < widthToFill; i++ {
		r.SetCh(x+i, y, Theme.Chars.FullBlock)
		r.SetFg(x+i, y, fg)
		r.SetBg(x+i, y, bg)
	}

	for ; i < width; i++ {
		r.SetCh(x+i, y, Theme.Chars.LightShade)
		r.SetFg(x+i, y, fg)
		r.SetBg(x+i, y, bg)
	}
}
