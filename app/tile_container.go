package app

import (
	"sync"

	"github.com/itglobal/dashboard/tile"
)

type tileContainer struct {
	manager    *tileManager
	provider   tile.Provider
	tiles      map[tile.ID]*tileObject
	mutex      sync.Mutex
	hasChanges bool
}

func newTileContainer(manager *tileManager, provider tile.Provider) *tileContainer {
	return &tileContainer{
		manager:  manager,
		provider: provider,
		tiles:    make(map[tile.ID]*tileObject),
	}
}

func (u *tileContainer) BeginUpdate() {
	u.mutex.Lock()
}

func (u *tileContainer) GetTiles() []tile.Tile {
	var tiles = make([]tile.Tile, len(u.tiles))

	i := 0
	for _, t := range u.tiles {
		tiles[i] = t
		i++
	}

	return tiles
}

func (u *tileContainer) AddOrUpdateTile(id tile.ID) tile.Tile {
	tile, exists := u.tiles[id]
	if !exists {
		tile = newTileObject(u, id)
		u.tiles[id] = tile
		u.hasChanges = true
	}

	return tile
}

func (u *tileContainer) RemoveTile(id tile.ID) {
	if _, exists := u.tiles[id]; !exists {
		return
	}

	delete(u.tiles, id)
	u.hasChanges = true
}

func (u *tileContainer) EndUpdate() {
	hasChanges := u.hasChanges
	u.mutex.Unlock()

	if hasChanges {
		u.manager.notifyChanged(u.provider)
	}
}

func (u *tileContainer) Invalidate() {
	u.hasChanges = true
}

func (u *tileContainer) getJSON(arr []tileJSON) []tileJSON {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, tile := range u.tiles {
		jsonObject := tile.getJSON()
		arr = append(arr, jsonObject)
	}

	return arr
}
