package app

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/kpango/glg"
)

var upgrader = websocket.Upgrader{} // use default options

type wsAPIHandler struct {
	manager *tileManager
}

type wsMsg struct {
	messageType int
	p           []byte
	err         error
}

func (h *wsAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		panic(err)
	}

	log.Debugf("WS: %s connected", ws.RemoteAddr())
	defer func() {
		log.Debugf("WS: %s disconnected", ws.RemoteAddr())
	}()

	data := h.manager.getJSON()
	err = ws.WriteJSON(data)
	if err != nil {
		ws.Close()
		panic(err)
	}

	subscription := h.manager.subscribe()
	defer func() {
		h.manager.unsubscribe(subscription)
	}()

	recv := make(chan wsMsg)
	update := subscription.Channel

	go func() {
		for {
			t, msg, err := ws.ReadMessage()
			recv <- wsMsg{t, msg, err}

			if err != nil {
				return
			}
		}
	}()

	for {
		select {
		case m := <-recv:
			if m.err != nil {
				ws.Close()
				return
			}
		case <-update:
			data := h.manager.getJSON()
			err = ws.WriteJSON(data)
			if err != nil {
				ws.Close()
				panic(err)
			}
		}
	}
}
