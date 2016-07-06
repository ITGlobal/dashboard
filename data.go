package main

import (
	"sync"
	"time"

	dash "github.com/itglobal/dashboard/api"
)

// Maps provider keys to list of provider's items
var itemsByProvider = make(map[string][]*dash.Item)

// A mutex for itemsByProvider
var itemsByProviderLock = sync.Mutex{}

var drawRequestChan = make(chan int)

// Issues a redraw request (non-blocking)
func requestDraw() {
	go func() {
		timerStateLock.Lock()
		defer timerStateLock.Unlock()

		if timerState == tsIdle {
			timerState = tsWait
		}
	}()
}

// Callback is a function to push updates from data providers into the UI loop
func Callback(provider dash.Provider, items []*dash.Item) {
	updateProviderData(provider, items)
	//	log.Printf("[dash] Provider updated: '%s'", provider.Key())
	requestDraw()
}

func updateProviderData(provider dash.Provider, items []*dash.Item) {
	itemsByProviderLock.Lock()
	defer itemsByProviderLock.Unlock()

	itemsByProvider[provider.Key()] = items
}

var timerState = tsIdle
var timerStateLock sync.Mutex
var timer = time.NewTicker(timerDuration)

type timerStateType int

const (
	tsIdle   timerStateType = 0
	tsWait   timerStateType = 1
	tsInProc timerStateType = 2
)

const timerDuration = 1 * time.Second

func timerHandler() {
	for {
		select {
		case <-timer.C:
			if checkForRedrawRequest() {
				go func() { redrawAndReset() }()
			}
		}
	}
}

func checkForRedrawRequest() bool {
	timerStateLock.Lock()
	defer timerStateLock.Unlock()

	if timerState == tsWait {
		timerState = tsInProc
	}

	return timerState != tsIdle
}

func redrawAndReset() {
	drawRequestChan <- 0

	timerStateLock.Lock()
	defer timerStateLock.Unlock()

	timerState = tsIdle
}
