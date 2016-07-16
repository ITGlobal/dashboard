package main

import (
	"sort"
	"strings"
	"sync"

	dash "github.com/itglobal/dashboard/api"
)

var Providers []dash.Provider

var dataByProvider = make(map[string][]*dash.Item)
var dataLock = sync.Mutex{}
var dataVersion uint = 0

func updateData(provider dash.Provider, items []*dash.Item) {
	dataLock.Lock()
	defer dataLock.Unlock()

	dataByProvider[provider.Key()] = items
	dataVersion++
}

func getData() ([]dash.Item, uint) {
	dataLock.Lock()
	defer dataLock.Unlock()

	var result []dash.Item
	sortedProviders := sortProviders(dataByProvider)
	for _, items := range sortedProviders {
		sortedItems := sortItems(items)
		for _, item := range sortedItems {
			result = append(result, *item)
		}
	}

	return result, dataVersion
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
