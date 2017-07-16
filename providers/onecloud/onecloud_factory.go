package onecloud

import (
	"time"

	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
	"github.com/satori/go.uuid"
)

const (
	providerType  = "1cloud"
	defaultURL    = "https://api.1cloud.ru"
	defaultPeriod = "60m"
)

func init() {
	tile.RegisterFactory(&oneCloudProviderFactory{})
}

type oneCloudProviderFactory struct{}

// Gets provider type key
func (f *oneCloudProviderFactory) Type() string {
	return providerType
}

// Create a provider
func (f *oneCloudProviderFactory) Create(config tile.Config, manager tile.Manager) (tile.Provider, error) {
	addr := tile.GetStringDefault(config, "url", defaultURL)
	if addr == "" {
		log.Errorf("Missing \"url\" parameter")
		return nil, tile.ErrBadConfig
	}

	name, ok := config.GetString("name")
	if !ok {
		log.Errorf("Missing \"name\" parameter")
		return nil, tile.ErrBadConfig
	}

	token, ok := config.GetString("token")
	if !ok {
		log.Errorf("Missing \"token\" parameter")
		return nil, tile.ErrBadConfig
	}

	periodRaw := tile.GetStringDefault(config, "timer", defaultPeriod)
	period, err := time.ParseDuration(periodRaw)
	if err != nil {
		log.Errorf("Parameter \"timer\" is not a valid duration: \"%s\"", periodRaw)
		return nil, tile.ErrBadConfig
	}

	provider := &oneCloudProvider{
		id:      tile.ID(uuid.NewV4().String()),
		url:     addr,
		name:    name,
		token:   token,
		period:  period,
		manager: manager,
	}

	return provider, nil
}
