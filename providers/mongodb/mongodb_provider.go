package mongodb

import (
	"fmt"
	"log"
	"time"

	dash "github.com/itglobal/dashboard/api"

	"gopkg.in/mgo.v2"
)

const key = "mongodb"

type mongodbProvider struct {
	callback dash.Callback
	name     string
	url      string
	interval time.Duration

	item *dash.Item
}

func (p *mongodbProvider) Key() string {
	return key
}

func (p *mongodbProvider) init() {
	item := new(dash.Item)
	item.Key = p.name
	item.ProviderKey = key
	item.Name = p.name
	item.Status = dash.StatusPending
	item.StatusText = "?"
	item.Progress = dash.NoProgress

	p.item = item

	p.update()
	go p.loop()
}

func (p *mongodbProvider) loop() {
	for {
		p.update()
		p.callback(p, []*dash.Item{p.item})

		time.Sleep(p.interval)
	}
}

func (p *mongodbProvider) update() {
	session, err := mgo.Dial(p.url)
	if err != nil {
		log.Printf("[mongodb] %s: dial() failed: %s", p.name, err)

		p.item.Status = dash.StatusBad
		p.item.StatusText = fmt.Sprintf("dial(): %s", err)
		p.item.Progress = dash.NoProgress
		return
	}

	defer session.Close()

	result := replStatus{}
	err = session.DB("admin").Run("replSetGetStatus", &result)
	if err != nil {
		log.Printf("[mongodb] %s: r.status() failed: %s", p.name, err)

		p.item.Status = dash.StatusBad
		p.item.StatusText = fmt.Sprintf("rs.status(): %s", err)
		p.item.Progress = dash.NoProgress
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
		p.item.Status = dash.StatusGood
		p.item.StatusText = "All nodes are healthy"
		p.item.Progress = 100

		return
	}

	health := float64(100) * float64(alive) / float64(len(result.Members))

	p.item.Status = dash.StatusBad
	p.item.StatusText = fmt.Sprintf("%d/%d nodes are healthy", alive, len(result.Members))
	p.item.Progress = int(health)
}

type replStatus struct {
	Set     string       `bson:"set"`
	Members []replMember `bson:"members"`
}

type replMember struct {
	Name   string `bson:"name"`
	Health int    `bson:"health"`
}

func factory(config dash.Config, callback dash.Callback) (dash.Provider, error) {
	url, err := config.GetString("url")
	if err != nil {
		return nil, err
	}

	name, err := config.GetString("name")
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(config.GetStringOrDefault("timer", "30s"))
	if err != nil {
		return nil, err
	}

	p := new(mongodbProvider)
	p.callback = callback
	p.name = name
	p.url = url
	p.interval = interval

	p.init()

	return p, nil
}

func init() {
	dash.RegisterFactory(key, factory)
}
