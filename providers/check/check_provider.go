package check

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/kpango/glg"

	"github.com/itglobal/dashboard/tile"
)

type pingProvider struct {
	id      tile.ID
	url     string
	name    string
	manager tile.Manager
	period  time.Duration
}

// Gets provider unique ID
func (p *pingProvider) ID() string {
	return string(p.id)
}

// Gets provider type key
func (p *pingProvider) Type() string {
	return providerType
}

// Initializes a provider
func (p *pingProvider) Init() error {
	u := p.manager.BeginUpdate(p)

	t := u.AddOrUpdateTile(p.id)
	t.SetType(tile.TypeTextStatus)
	t.SetSize(tile.Size1x)
	t.SetTitleText(p.name)
	t.SetState(tile.StateIndeterminate)
	u.EndUpdate()

	p.checkOnce()
	go func() {
		for {
			p.checkOnce()

			time.Sleep(p.period)
		}
	}()

	return nil
}

func (p *pingProvider) checkOnce() {
	var state tile.State
	var text string

	resp, err := http.Get(p.url)
	if err != nil {
		log.Errorf("[check] GET %s failed: %s", p.url, err)
		state = tile.StateError
		text = err.Error()
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			log.Successf("[check] GET %s -> %s", p.url, resp.Status)
			state = tile.StateSuccess
		} else {
			log.Errorf("[check] GET %s -> %s", p.url, resp.Status)
			state = tile.StateError
		}

		text = resp.Status
	}

	now := time.Now()
	text = fmt.Sprintf("[%s] %s", now.Format("15:04:05"), text)

	u := p.manager.BeginUpdate(p)
	defer u.EndUpdate()

	t := u.AddOrUpdateTile(p.id)

	t.SetState(state)
	t.SetDescriptionText(text)
}
