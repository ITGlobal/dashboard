package mongodb

import (
	"fmt"
	"time"

	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
	"gopkg.in/mgo.v2"
)

type mongodbProvider struct {
	id              tile.ID
	manager         tile.Manager
	interval        time.Duration
	url             string
	name            string
	state           tile.State
	descriptionText string
	statusValue     int
}

// Gets provider unique ID
func (p *mongodbProvider) ID() string {
	return string(p.id)
}

// Gets provider type key
func (p *mongodbProvider) Type() string {
	return providerType
}

// Initializes a provider
func (p *mongodbProvider) Init() error {
	p.state = tile.StateIndeterminate
	p.descriptionText = ""
	p.statusValue = 100

	p.syncTile()

	go func() {
		for {
			p.update()
			p.syncTile()

			time.Sleep(p.interval)
		}
	}()

	return nil
}

func (p *mongodbProvider) syncTile() {
	u := p.manager.BeginUpdate(p)
	defer u.EndUpdate()

	t := u.AddOrUpdateTile(p.id)

	t.SetType(tile.TypeTextStatusProgress)
	t.SetSize(tile.Size1x)
	t.SetTitleText(p.name)
	t.SetState(p.state)
	t.SetDescriptionText(p.descriptionText)
	t.SetStatusValue(p.statusValue)
}

func (p *mongodbProvider) update() {
	session, err := mgo.Dial(p.url)
	if err != nil {
		log.Printf("[mongodb] %s: dial() failed: %s", p.name, err)

		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("dial(): %s", err)
		p.statusValue = 0

		return
	}

	defer session.Close()

	result := replStatus{}
	err = session.DB("admin").Run("replSetGetStatus", &result)
	if err != nil {
		log.Printf("[mongodb] %s: r.status() failed: %s", p.name, err)

		p.state = tile.StateError
		p.descriptionText = fmt.Sprintf("rs.status(): %s", err)
		p.statusValue = 0
		return
	}

	alive := 0
	for _, member := range result.Members {
		if member.Health == 1 {
			alive++
		}
	}

	log.Printf("[mongodb] %s: %d/%d nodes are healthy", p.name, alive, len(result.Members))

	if alive == len(result.Members) {
		p.state = tile.StateSuccess
		p.descriptionText = "All nodes are healthy"
		p.statusValue = 100

		return
	}

	health := float64(100) * float64(alive) / float64(len(result.Members))

	p.state = tile.StateWarning
	p.descriptionText = fmt.Sprintf("%d/%d nodes are healthy", alive, len(result.Members))
	p.statusValue = int(health)
}

type replStatus struct {
	Set     string       `bson:"set"`
	Members []replMember `bson:"members"`
}

type replMember struct {
	Name   string `bson:"name"`
	Health int    `bson:"health"`
}
