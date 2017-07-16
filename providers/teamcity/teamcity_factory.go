package teamcity

import (
	"time"

	"github.com/itglobal/dashboard/tile"
	"github.com/kapitanov/go-teamcity"
	log "github.com/kpango/glg"
	"github.com/satori/go.uuid"
)

const (
	providerType         = "teamcity"
	defaultFetchInterval = "20s"
)

func init() {
	tile.RegisterFactory(&teamcityProviderFactory{})
}

type teamcityProviderFactory struct{}

// Gets provider type key
func (f *teamcityProviderFactory) Type() string {
	return providerType
}

// Create a provider
func (f *teamcityProviderFactory) Create(config tile.Config, manager tile.Manager) (tile.Provider, error) {
	url, exists := config.GetString("url")
	if !exists {
		log.Errorf("Missing \"url\" parameter")
		return nil, tile.ErrBadConfig
	}

	username := tile.GetStringDefault(config, "username", "")
	password := tile.GetStringDefault(config, "password", "")

	fetchIntervalRaw := tile.GetStringDefault(config, "timer", defaultFetchInterval)
	fetchInterval, err := time.ParseDuration(fetchIntervalRaw)
	if err != nil {
		log.Errorf("Parameter \"timer\" is not a valid duration: \"%s\"", fetchIntervalRaw)
		return nil, tile.ErrBadConfig
	}

	var auth teamcity.Authorizer
	if username != "" && password != "" {
		auth = teamcity.BasicAuth(username, password)
	} else {
		auth = teamcity.GuestAuth()
	}

	provider := &teamcityProvider{
		uuid.NewV4().String(),
		teamcity.NewClient(url, auth),
		manager,
		fetchInterval,
		make(map[tile.ID]*teamcityDataItem),
		make(map[string]*teamcityDataItem),
	}

	return provider, nil
}
