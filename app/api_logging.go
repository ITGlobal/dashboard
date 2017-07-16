package app

import (
	"bufio"
	"net"
	"net/http"
	"time"

	log "github.com/kpango/glg"
)

func withLogging(handler http.Handler) http.Handler {
	return &logHTTPHandler{handler}
}

type logHTTPHandler struct {
	handler http.Handler
}

func (h *logHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	wr := &httpResponseInterceptor{0, w, nil}
	h.serveImpl(wr, req)

	durationMS := time.Now().Sub(startTime).Nanoseconds() / 1000
	url := req.URL.RequestURI()

	if wr.err != nil {
		log.Errorf("%s %s -> %d: %s", req.Method, url, wr.status, wr.err)
	} else {
		log.Debugf("%s %s -> %d in %dms", req.Method, url, wr.status, durationMS)
	}
}

func (h *logHTTPHandler) serveImpl(w *httpResponseInterceptor, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.err = err.(error)
		}
	}()

	w.Header().Set("X-Server", "dashboard")
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

func (i *httpResponseInterceptor) Hijack() (c net.Conn, rw *bufio.ReadWriter, err error) {
	hijacker, ok := i.w.(http.Hijacker)
	if !ok {
		panic(http.ErrHijacked)
	}

	return hijacker.Hijack()
}
