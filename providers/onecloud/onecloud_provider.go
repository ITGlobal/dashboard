package onecloud

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"io/ioutil"

	"math"

	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
)

type oneCloudProvider struct {
	id              tile.ID
	name            string
	url             string
	token           string
	period          time.Duration
	manager         tile.Manager
	state           tile.State
	descriptionText string
	statusValue     int
}

const (
	goodBoundary    = 14
	warningBoundary = 7
	rubleSign       = " RUB"
)

// Gets provider unique ID
func (p *oneCloudProvider) ID() string {
	return string(p.id)
}

// Gets provider type key
func (p *oneCloudProvider) Type() string {
	return providerType
}

// Initializes a provider
func (p *oneCloudProvider) Init() error {
	p.state = tile.StateIndeterminate
	p.descriptionText = ""

	p.syncTile()

	go func() {
		for {
			p.updateOnce()
			p.syncTile()

			time.Sleep(p.period)
		}
	}()

	return nil
}

func (p *oneCloudProvider) syncTile() {
	u := p.manager.BeginUpdate(p)
	defer u.EndUpdate()

	t := u.AddOrUpdateTile(p.id)
	t.SetType(tile.TypeTextStatus2)
	t.SetSize(tile.Size1x)
	t.SetTitleText(p.name)
	t.SetDescriptionText(p.descriptionText)
	t.SetState(p.state)
	t.SetStatusValue(p.statusValue)
}

func (p *oneCloudProvider) updateOnce() {
	http.Get("")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account", p.url), nil)
	if err != nil {
		log.Errorf("[1cloud] Unable to create request: %s", err)
		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("http.NewRequest(): %s", err)
		p.statusValue = 0
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("[1cloud] %s %s failed: %s", req.Method, req.URL, err)
		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("http.Get(): %s", err)
		p.statusValue = 0
		return
	}

	if resp.StatusCode != 200 {
		log.Errorf("[1cloud] %s %s failed with %s", req.Method, req.URL, resp.Status)
		p.state = tile.StateError
		p.descriptionText = resp.Status
		p.statusValue = 0
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[1cloud] %s %s failed (read error): %s", req.Method, req.URL, err)
		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("ioutil.ReadAll(): %s", err)
		p.statusValue = 0
		return
	}

	var response oneCloudResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Errorf("[1cloud] %s %s failed (parse error): %s", req.Method, req.URL, err)
		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("json.Unmarshal(): %s", err)
		p.statusValue = 0
		return
	}

	balanceTill, err := time.Parse("2006-01-02", response.BalanceTill)
	if err != nil {
		log.Errorf("[1cloud] %s %s failed (time parse error): %s", req.Method, req.URL, err)
		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("time.Parse(): %s", err)
		p.statusValue = 0
		return
	}

	duration := balanceTill.Sub(time.Now().UTC())
	days := int(math.Floor(duration.Hours() / 24.0))
	if days < 0 {
		days = 0
	}

	p.statusValue = days
	p.descriptionText = fmt.Sprintf("Paid until %s (%0.0f%s)", balanceTill.Format("2 Jan 2006"), response.Balance, rubleSign)

	if days >= goodBoundary {
		p.state = tile.StateSuccess
		log.Debugf("[1cloud] paid until %s (%d days, %0.0f RUB) - that's good", balanceTill, days, response.Balance)
	} else {
		if days >= warningBoundary {
			p.state = tile.StateWarning
			log.Debugf("[1cloud] paid until %s (%d days, %0.0f RUB) - that's not very good", balanceTill, days, response.Balance)
		} else {
			p.state = tile.StateError
			log.Debugf("[1cloud] paid until %s (%d days, %0.0f RUB) - that's time to panic", balanceTill, days, response.Balance)
		}
	}
}

type oneCloudResponse struct {
	Balance     float32 `json:"Balance"`
	BalanceTill string  `json:"BalanceTillDateUtc"`
}
