package sim

import (
	"fmt"
	"time"

	"github.com/itglobal/dashboard/tile"
)

const animStepDuration = 1 * time.Second
const stateStepDuration = 5 * time.Second

type simProvider struct {
	uid     string
	ids     []tile.ID
	manager tile.Manager
}

// Gets provider unique ID
func (p *simProvider) ID() string {
	return p.uid
}

// Gets provider type key
func (p *simProvider) Type() string {
	return providerType
}

// Initializes a provider
func (p *simProvider) Init() error {
	op := p.manager.BeginUpdate(p)

	for i := 0; i < len(p.ids); i++ {
		p.ids[i] = tile.ID(newUID())

		t := op.AddOrUpdateTile(p.ids[i])

		t.SetType(tile.TypeTextStatusProgress)
		t.SetSize(tile.Size2x)
		t.SetState(tile.StateError)
		t.Clear()
		t.SetTitleText(fmt.Sprintf("generated item %d", i))
	}

	op.EndUpdate()

	for i, id := range p.ids {
		go p.animate(i, id)
	}

	return nil
}

type animationFunc func(p *simProvider, id tile.ID)

func pendingAnimation(p *simProvider, id tile.ID) {
	for progress := 0; progress <= 100; progress += 10 {

		op := p.manager.BeginUpdate(p)

		t := op.AddOrUpdateTile(id)
		t.Clear()
		t.SetState(tile.StateIndeterminate)
		t.SetDescriptionText(fmt.Sprintf("Running (%d%%)", progress))
		t.SetStatusValue(progress)

		op.EndUpdate()

		time.Sleep(animStepDuration)
	}
}

func goodAnimation(p *simProvider, id tile.ID) {
	op := p.manager.BeginUpdate(p)

	t := op.AddOrUpdateTile(id)
	t.Clear()
	t.SetState(tile.StateSuccess)
	t.SetDescriptionText("Successful result text content")

	op.EndUpdate()

	time.Sleep(stateStepDuration)
}

func badAnimation(p *simProvider, id tile.ID) {
	op := p.manager.BeginUpdate(p)

	t := op.AddOrUpdateTile(id)
	t.Clear()
	t.SetState(tile.StateError)
	t.SetDescriptionText("Unsuccessful result text content")

	op.EndUpdate()

	time.Sleep(stateStepDuration)
}

var animationSteps = []animationFunc{pendingAnimation, goodAnimation, pendingAnimation, badAnimation}

func (p *simProvider) animate(index int, id tile.ID) {
	index = index % len(animationSteps)

	for {
		animationSteps[index](p, id)

		index++
		if index >= len(animationSteps) {
			index = 0
		}
	}
}
