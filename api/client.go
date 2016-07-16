package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Connection is an interface for connection to dash
type Connection interface {
	// Close closes connection to dashd
	Close()
}

// ClientCallback is a function that is triggered when any new data is received for dash
type ClientCallback func([]Item)

// NewConnection creates new connection to dashd
func NewConnection(logger *log.Logger, url string, callback ClientCallback) Connection {
	c := &connection{logger, url, callback, 0, make(chan int), make(chan int)}
	go c.run()
	return c
}

const fetchInterval = time.Second

type connection struct {
	logger   *log.Logger
	url      string
	callback ClientCallback
	version  uint
	stop     chan int
	stopped  chan int
}

func (c *connection) run() {
	c.fetchOnce()

	sleep := time.After(fetchInterval)

	for {
		select {
		case <-sleep:
			c.fetchOnce()
			sleep = time.After(fetchInterval)
		case <-c.stop:
			go func() { c.stopped <- 0 }()
			return
		}
	}
}

func (c *connection) Close() {
	go func() { c.stop <- 0 }()
	<-c.stopped
}

func (c *connection) fetchOnce() {
	url := fmt.Sprintf("%s?v=%d", c.url, c.version)
	resp, err := http.Get(url)
	if err != nil {
		c.logger.Printf("GET %s failed with %s", url, err)
		return
	}

	if resp.StatusCode == http.StatusNotModified {
		return
	}

	defer resp.Body.Close()

	var result ItemList
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Printf("Unable to parse result of GET %s: %s", url, err)
		return
	}

	c.version = result.Version
	c.callback(result.Items)
}
