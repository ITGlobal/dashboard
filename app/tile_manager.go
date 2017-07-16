package app

import (
	"sort"
	"sync"

	"github.com/itglobal/dashboard/tile"
	"github.com/satori/go.uuid"
)

type tileManagerSubscrition struct {
	id      int
	Channel chan interface{}
}

type tileManager struct {
	containers        map[tile.Provider]*tileContainer
	mutex             sync.Mutex
	versionLock       sync.Mutex
	version           string
	subscriptionsLock sync.Mutex
	subscriptionsId   int
	subscriptions     map[int]*tileManagerSubscrition
}

func newTileManager() *tileManager {
	return &tileManager{
		containers:    make(map[tile.Provider]*tileContainer),
		subscriptions: make(map[int]*tileManagerSubscrition),
		version:       uuid.Nil.String(),
	}
}

func (m *tileManager) subscribe() *tileManagerSubscrition {
	m.subscriptionsLock.Lock()
	defer m.subscriptionsLock.Unlock()
	m.subscriptionsId++
	id := m.subscriptionsId
	s := &tileManagerSubscrition{id, make(chan interface{})}
	m.subscriptions[id] = s
	return s
}

func (m *tileManager) unsubscribe(s *tileManagerSubscrition) {
	m.subscriptionsLock.Lock()
	defer m.subscriptionsLock.Unlock()
	delete(m.subscriptions, s.id)
}

func (m *tileManager) BeginUpdate(provider tile.Provider) tile.Updater {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	container, exists := m.containers[provider]
	if !exists {
		container = newTileContainer(m, provider)
		m.containers[provider] = container
	}

	container.BeginUpdate()
	return container
}

func (m *tileManager) notifyChanged(p tile.Provider) {
	m.versionLock.Lock()
	newVersion := uuid.NewV4().String()
	m.version = newVersion
	m.versionLock.Unlock()
	m.subscriptionsLock.Lock()
	for _, t := range m.subscriptions {
		ch := t.Channel
		go func() { ch <- nil }()
	}
	m.subscriptionsLock.Unlock()
}

func (m *tileManager) getJSON() *dataJSON {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var arr []tileJSON

	for _, container := range m.containers {
		arr = container.getJSON(arr)
	}

	var less = func(i, j int) bool {
		x := arr[i]
		y := arr[j]

		if x.Source < y.Source {
			return true
		}

		if x.Source > y.Source {
			return false
		}

		if x.ID < y.ID {
			return true
		}

		return false
	}

	sort.SliceStable(arr, less)

	m.versionLock.Lock()
	defer m.versionLock.Unlock()

	return &dataJSON{
		Version: m.version,
		Tiles:   arr,
	}
}
