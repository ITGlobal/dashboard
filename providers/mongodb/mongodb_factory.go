package mongodb

import (
	"time"

	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
	"github.com/satori/go.uuid"
)

const (
	providerType         = "mongodb"
	defaultFetchInterval = "20s"
)

func init() {
	tile.RegisterFactory(&mongodbProviderFactory{})
}

type mongodbProviderFactory struct{}

// Gets provider type key
func (f *mongodbProviderFactory) Type() string {
	return providerType
}

// Create a provider
func (f *mongodbProviderFactory) Create(config tile.Config, manager tile.Manager) (tile.Provider, error) {
	url, exists := config.GetString("url")
	if !exists {
		log.Errorf("Missing \"url\" parameter")
		return nil, tile.ErrBadConfig
	}

	name, exists := config.GetString("name")
	if !exists {
		log.Errorf("Missing \"name\" parameter")
		return nil, tile.ErrBadConfig
	}

	fetchIntervalRaw := tile.GetStringDefault(config, "timer", defaultFetchInterval)
	fetchInterval, err := time.ParseDuration(fetchIntervalRaw)
	if err != nil {
		log.Errorf("Parameter \"timer\" is not a valid duration: \"%s\"", fetchIntervalRaw)
		return nil, tile.ErrBadConfig
	}

	provider := &mongodbProvider{
		id:       tile.ID(uuid.NewV4().String()),
		manager:  manager,
		interval: fetchInterval,
		url:      url,
		name:     name,
	}

	return provider, nil
}
