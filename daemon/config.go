package main

// Common imports
import (
	"encoding/json"
	"errors"
	"io/ioutil"

	dash "github.com/itglobal/dashboard/api"

	_ "expvar"
)

func createProviders(callback dash.Callback) []dash.Provider {
	// Read config
	logger.Printf("reading config file '%s'", ConfigFileName)
	configBytes, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		logger.Printf("unable to read config file '%s': %s", ConfigFileName, err)
		panic(err)
	}

	// Parse config
	var providerConfigs []interface{}
	err = json.Unmarshal(configBytes, &providerConfigs)
	if err != nil {
		logger.Printf("unable to parse config file '%s': %s", ConfigFileName, err)
		panic(err)
	}

	// Create providers
	var providers []dash.Provider
	for i, providerConfig := range providerConfigs {
		provider := createProvider(callback, i, providerConfig)
		if provider == nil {
			continue
		}

		providers = append(providers, *provider)
	}

	return providers
}

func createProvider(callback dash.Callback, i int, providerConfig interface{}) *dash.Provider {
	config, err := newConfigReader(providerConfig)
	if err != nil {
		logger.Printf("item #%d is misconfigured: %s", i, err)
		return nil
	}

	t, err := config.GetString("type")
	if err != nil {
		logger.Printf("item #%d has no 'type' property: %s", i, err)
		return nil
	}

	factory := dash.GetFactory(t)
	if factory == nil {
		logger.Printf("item #%d - provider %s is unknown: %s", i, t, err)
		return nil
	}

	provider, err := factory(config, callback)
	if err != nil {
		logger.Printf("item #%d - unable to create a %s provider: %s", i, t, err)
		return nil
	}

	return &provider
}

type configReader map[string]interface{}

func newConfigReader(data interface{}) (dash.Config, error) {
	switch data := data.(type) {
	case map[string]interface{}:
		return configReader(data), nil
	default:
		return nil, errors.New("Bad JSON data")
	}
}

func (c configReader) GetString(key string) (string, error) {
	value, exists := c[key]
	if !exists {
		return "", dash.ErrNoSuchKey
	}

	return value.(string), nil
}

func (c configReader) GetStringOrDefault(key string, def string) string {
	value, err := c.GetString(key)
	if err != nil {
		return def
	}

	return value
}

func (c configReader) GetInt32(key string) (int32, error) {
	value, exists := c[key]
	if !exists {
		return 0, dash.ErrNoSuchKey
	}

	switch v := value.(type) {
	case int32:
		return v, nil
	case float64:
		return int32(v), nil
	default:
		return 0, dash.ErrNoSuchKey
	}
}

func (c configReader) GetInt32OrDefault(key string, def int32) int32 {
	value, err := c.GetInt32(key)
	if err != nil {
		return def
	}

	return value
}
