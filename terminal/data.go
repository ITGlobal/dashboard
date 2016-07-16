package main

import (
	"sync"
	"time"

	dash "github.com/itglobal/dashboard/api"
)

var items []dash.Item
var itemsLock = sync.Mutex{}
var drawRequestChan = make(chan int)

func requestDraw() {
	go func() {
		timerStateLock.Lock()
		defer timerStateLock.Unlock()

		if timerState == tsIdle {
			timerState = tsWait
		}
	}()
}

func Callback(is []dash.Item) {
	updateData(is)
	requestDraw()
}

func updateData(is []dash.Item) {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	items = is
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
