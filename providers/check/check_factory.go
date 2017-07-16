package check

import (
	"net/url"
	"time"

	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
	"github.com/satori/go.uuid"
)

const (
	providerType  = "check"
	defaultPeriod = "1m"
)

func init() {
	tile.RegisterFactory(&checkProviderFactory{})
}

type checkProviderFactory struct{}

// Gets provider type key
func (f *checkProviderFactory) Type() string {
	return providerType
}

// Create a provider
func (f *checkProviderFactory) Create(config tile.Config, manager tile.Manager) (tile.Provider, error) {
	addr, ok := config.GetString("url")
	if !ok {
		log.Errorf("Missing \"url\" parameter")
		return nil, tile.ErrBadConfig
	}

	u, err := url.Parse(addr)
	if err != nil {
		log.Errorf("Parameter \"url\" is not a valid URL: \"%s\"", addr)
		return nil, tile.ErrBadConfig
	}

	periodRaw := tile.GetStringDefault(config, "timer", defaultPeriod)
	period, err := time.ParseDuration(periodRaw)
	if err != nil {
		log.Errorf("Parameter \"timer\" is not a valid duration: \"%s\"", periodRaw)
		return nil, tile.ErrBadConfig
	}

	name := tile.GetStringDefault(config, "name", u.Host)

	provider := &pingProvider{
		tile.ID(uuid.NewV4().String()),
		addr,
		name,
		manager,
		period,
	}

	return provider, nil
}
