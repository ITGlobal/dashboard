package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dash "github.com/itglobal/dashboard/api"
)

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	minVersion := parseMinVersion(r)

	items, version := getData()
	if version <= minVersion {
		// 304 Not Modified
		http304(w, version)
	} else {
		// 200 OK
		http200(w, items, version)
	}
}

func parseMinVersion(r *http.Request) uint {
	str := r.URL.Query().Get("v")
	if str == "" {
		return 0
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	if value < 0 {
		return 0
	}

	return uint(value)
}

func http304(w http.ResponseWriter, version uint) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("ETag", fmt.Sprintf("%d", version))
	w.WriteHeader(http.StatusNotModified)

	writeJson(w, nil)
}

func http200(w http.ResponseWriter, items []dash.Item, version uint) {
	result := dash.ItemList{version, items}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("ETag", fmt.Sprintf("%d", version))
	w.WriteHeader(http.StatusOK)

	writeJson(w, result)
}

func writeJson(w http.ResponseWriter, data interface{}) {
	buffer, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	_, err = w.Write(buffer)
	if err != nil {
		panic(err)
	}
}
