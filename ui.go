package main

import (
	"log"
	"sort"
	"strings"

	dash "github.com/itglobal/dashboard/api"

	"github.com/nsf/termbox-go"
)

// Run is an entrypoint for UI subsystem
func Run() {
	// Initialize UI
	err := termbox.Init()
	if err != nil {
		log.Panicf("Unable to init terminal: %s\n", err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	// Push an invalidation request
	requestDraw()

	events := make(chan termbox.Event, 1)

	go func() {
		for {
			e := termbox.PollEvent()
			events <- e

			if e.Type == termbox.EventKey && e.Key == termbox.KeyEsc {
				return
			}
		}
	}()

	go func() {
		timerHandler()
	}()

	for {
		select {
		case <-drawRequestChan:
			redraw("REQUEST")
		case e := <-events:
			if e.Type == termbox.EventKey && e.Key == termbox.KeyEsc {
				return
			}
			if e.Type == termbox.EventResize {
				redraw("RESIZE")
			}
		}
	}
}

func redraw(reason string) {
	termbox.Sync()
	width, _ := termbox.Size()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	renderLayout(width)

	termbox.Flush()
}

func calcItemWidth(w int) int {
	c := w / 30

	if c >= 12 {
		return w / 12
	}
	if c >= 6 {
		return w / 6
	}
	if c >= 4 {
		return w / 4
	}
	if c >= 3 {
		return w / 3
	}
	if c >= 2 {
		return w / 2
	}

	return w
}

func renderLayout(width int) {
	itemsByProviderLock.Lock()
	defer itemsByProviderLock.Unlock()

	i := 0
	y := 0

	itemWidth := calcItemWidth(width)
	itemCount := width / itemWidth
	itemHeight := 7
	extraWidthForLast := width - itemCount*itemWidth
	sortedProviders := sortProviders(itemsByProvider)

	for _, items := range sortedProviders {
		sortedItems := sortItems(items)
		for _, item := range sortedItems {
			currentItemWidth := itemWidth
			if i == itemCount-1 {
				currentItemWidth += extraWidthForLast
			}

			renderer := NewItemRenderer(i*itemWidth, y, currentItemWidth, itemHeight)
			RenderCell(item, renderer)

			if i == itemCount-1 {
				i = 0
				y += itemHeight
				continue
			}

			i++
		}
	}
}

func sortProviders(m map[string][]*dash.Item) [][]*dash.Item {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	result := make([][]*dash.Item, len(keys))
	for i, key := range keys {
		result[i] = m[key]
	}
	return result
}

type sortedDashItems []*dash.Item

func (array sortedDashItems) Len() int {
	return len(array)
}

func (array sortedDashItems) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func (array sortedDashItems) Less(i, j int) bool {
	return strings.Compare(array[i].Key, array[j].Key) < 0
}

func sortItems(items []*dash.Item) []*dash.Item {
	sortable := sortedDashItems(items)
	sort.Sort(sortable)
	return sortable
}
