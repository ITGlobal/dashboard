package main

import (
	dash "github.com/itglobal/dashboard/api"
	"github.com/nsf/termbox-go"
)

func RenderCell(item *dash.Item, r ItemRenderer) {
	defer r.Render()

	w := r.Width()
	h := r.Height()

	fg, bg := getColors(item)
	r.Clear(fg, bg)

	// Render frame
	r.SetCh(0, 0, Theme.Chars.TopLeft)
	r.SetCh(w-1, 0, Theme.Chars.TopRight)
	r.SetCh(0, h-1, Theme.Chars.BottomLeft)
	r.SetCh(w-1, h-1, Theme.Chars.BottomRight)

	r.HLine(1, 0, w-2)
	r.HLine(1, h-1, w-2)

	r.VLine(0, 1, h-2)
	r.VLine(w-1, 1, h-2)

	// Render provider key
	r.SetCh(2, 0, '[')
	len := r.WriteText(3, 0, item.ProviderKey, w-10)
	r.SetCh(len+3, 0, ']')

	// Render name
	r.WriteText(1, 1, item.Name, w-2)
	r.HLineSingle(1, 2, w-2)
	r.SetCh(0, 2, Theme.Chars.VerticalAndRight)
	r.SetCh(w-1, 2, Theme.Chars.VerticalAndLeft)

	// Render status text
	i := r.WriteTextEx(1, 3, item.StatusText, w-2, 0, false)
	r.WriteTextEx(1, 4, item.StatusText, w-2, i, true)

	// Render progress bar
	if item.Progress != dash.NoProgress {
		r.ProgressBar(1, 5, w-2, item.Progress, fg, bg)
	}
}

func getColors(item *dash.Item) (termbox.Attribute, termbox.Attribute) {
	switch item.Status {
	case dash.StatusUnknown:
		return Theme.Colors.DefaultFg, Theme.Colors.DefaultBg
	case dash.StatusGood:
		return Theme.Colors.GoodFg, Theme.Colors.GoodBg
	case dash.StatusBad:
		return Theme.Colors.BadFg, Theme.Colors.BadBg
	case dash.StatusPending:
		return Theme.Colors.PendingFg, Theme.Colors.PendingBg
	}

	return Theme.Colors.DefaultFg, Theme.Colors.DefaultBg
}
