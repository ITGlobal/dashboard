package app

import (
	"encoding/json"
	"net/http"
)

type restAPIHandler struct {
	manager *tileManager
}

func (h *restAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := h.manager.getJSON()

	str := r.URL.Query().Get("v")

	if str != "" && data.Version == str {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("ETag", data.Version)
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.SetIndent("", "    ")
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
}
