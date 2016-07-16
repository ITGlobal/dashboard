package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type rootHttpHandler struct {
	logger  *log.Logger
	handler http.Handler
}

func CreateHandler(logger *log.Logger, handler http.Handler) http.Handler {
	return rootHttpHandler{logger, handler}
}

func (h rootHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	wr := &httpResponseInterceptor{0, w, nil}
	h.serveImpl(wr, req)

	durationMS := time.Now().Sub(startTime).Nanoseconds() / 1000
	url := req.URL.RequestURI()

	if wr.err != nil {
		h.logger.Printf("%s %s -> %d: %s", req.Method, url, wr.status, wr.err)
	} else {
		h.logger.Printf("%s %s -> %d in %dms", req.Method, url, wr.status, durationMS)
	}
}

func (h rootHttpHandler) serveImpl(w *httpResponseInterceptor, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.err = err.(error)
		}
	}()

	w.Header().Set("X-Server", fmt.Sprintf("dashd/%s", Version))
	h.handler.ServeHTTP(w, req)
}

type httpResponseInterceptor struct {
	status int
	w      http.ResponseWriter
	err    error
}

func (i *httpResponseInterceptor) Header() http.Header {
	return i.w.Header()
}

func (i *httpResponseInterceptor) Write(buffer []byte) (int, error) {
	return i.w.Write(buffer)
}

func (i *httpResponseInterceptor) WriteHeader(status int) {
	i.w.WriteHeader(status)
	i.status = status
}
