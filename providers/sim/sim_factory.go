package sim

import (
	"github.com/itglobal/dashboard/tile"

	"github.com/satori/go.uuid"
)

const (
	providerType = "sim"
	defaultCount = 10
)

func init() {
	tile.RegisterFactory(&simProviderFactory{})
}

type simProviderFactory struct{}

// Gets provider type key
func (f *simProviderFactory) Type() string {
	return providerType
}

// Create a provider
func (f *simProviderFactory) Create(config tile.Config, manager tile.Manager) (tile.Provider, error) {
	n := tile.GetIntegerDefault(config, "count", defaultCount)

	provider := &simProvider{
		newUID(),
		make([]tile.ID, n),
		manager,
	}

	return provider, nil
}

func newUID() string {
	return uuid.NewV4().String()
}
