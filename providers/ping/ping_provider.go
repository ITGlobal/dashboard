package ping

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	dash "github.com/itglobal/dashboard/api"
)

const key = "ping"

var providerIndex = 0

type pingProvider struct {
	key      string
	url      string
	item     *dash.Item
	items    []*dash.Item
	callback dash.Callback
}

func (p *pingProvider) Key() string {
	return p.item.Key
}

func (p *pingProvider) init(period time.Duration) {
	go p.loop(period)
}

func (p *pingProvider) loop(period time.Duration) {
	for {
		resp, err := http.Get(p.url)
		if err != nil {
			log.Printf("[ping] GET %s failed: %s", p.url, err)
			p.item.Status = dash.StatusBad
			p.item.StatusText = err.Error()
		} else {
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				log.Printf("[ping] GET %s -> %s", p.url, resp.Status)
				p.item.Status = dash.StatusGood
			} else {
				log.Printf("[ping] GET %s -> %s", p.url, resp.Status)
				p.item.Status = dash.StatusBad
			}

			p.item.StatusText = resp.Status
		}

		now := time.Now()
		p.item.StatusText = fmt.Sprintf("[%s] %s", now.Format("15:04:05"), p.item.StatusText)

		go p.callback(p, p.items)
		time.Sleep(period)
	}
}

func factory(config dash.Config, callback dash.Callback) (dash.Provider, error) {
	interval, err := time.ParseDuration(config.GetStringOrDefault("timer", "1m"))
	if err != nil {
		return nil, err
	}

	addr, err := config.GetString("url")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	item := new(dash.Item)
	item.Name = u.Host
	item.Status = dash.StatusPending
	item.StatusText = "?"
	item.Progress = dash.NoProgress

	providerIndex++

	p := new(pingProvider)
	p.key = fmt.Sprintf("%s-%d", key, providerIndex)
	p.item = item
	p.url = addr
	p.callback = callback
	p.items = make([]*dash.Item, 1)
	p.items[0] = item

	item.ProviderKey = key
	item.Key = p.key

	p.init(interval)

	return p, nil
}

func init() {
	dash.RegisterFactory(key, factory)
}
