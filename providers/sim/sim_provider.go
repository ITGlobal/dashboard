package sim

import (
	"fmt"
	"math/rand"
	"time"

	dash "github.com/itglobal/dashboard/api"
)

const key = "sim"

const animStepDuration = 1 * time.Second
const stateStepDuration = 5 * time.Second

type simProvider struct {
	items    []*dash.Item
	callback dash.Callback
}

func (p *simProvider) Key() string {
	return key
}

func (p *simProvider) init() {
	for _, item := range p.items {
		j := rand.Int()
		item.ProviderKey = key
		item.Key = fmt.Sprintf("%s_%d", key, j)
		item.Name = fmt.Sprintf("generated item %d", j)
		item.Progress = dash.NoProgress
		item.Status = dash.StatusBad
		item.StatusText = ""
	}

	go p.loop()
}

func (p *simProvider) loop() {
	for index, item := range p.items {
		go p.animate(item, index)
	}
}

type animation func(item *dash.Item, callback func())

func pendingAnimation(item *dash.Item, callback func()) {
	for progress := 0; progress <= 100; progress += 10 {
		item.Status = dash.StatusPending
		item.Progress = progress
		item.StatusText = fmt.Sprintf("Running (%d%%)", progress)

		callback()
		time.Sleep(animStepDuration)
	}
}

func goodAnimation(item *dash.Item, callback func()) {
	item.Status = dash.StatusGood
	item.Progress = dash.NoProgress
	item.StatusText = "Successful result text content"

	callback()
	time.Sleep(stateStepDuration)
}

func badAnimation(item *dash.Item, callback func()) {
	item.Status = dash.StatusBad
	item.Progress = dash.NoProgress
	item.StatusText = "Unsuccessful result text content"

	callback()
	time.Sleep(stateStepDuration)
}

var animationSteps = []animation{pendingAnimation, goodAnimation, pendingAnimation, badAnimation}

func (p *simProvider) animate(item *dash.Item, index int) {
	index = index % len(animationSteps)
	callback := func() { p.callback(p, p.items) }

	for {
		animationSteps[index](item, callback)

		index++
		if index >= len(animationSteps) {
			index = 0
		}
	}
}

func factory(config dash.Config, callback dash.Callback) (dash.Provider, error) {
	p := new(simProvider)
	p.callback = callback
	count := 10
	p.items = make([]*dash.Item, count)

	for i := 0; i < count; i++ {
		p.items[i] = new(dash.Item)
	}

	p.init()

	return p, nil
}

func init() {
	dash.RegisterFactory(key, factory)
}
