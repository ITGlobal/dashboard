package tile

import "errors"

// ErrBadConfig is an error that refers to wrong provider config parameters
var ErrBadConfig = errors.New("bad_config")

// Provider defines methods for sources that provide dashboard items
type Provider interface {
	// Gets provider unique ID
	ID() string
	// Gets provider type key
	Type() string
	// Initializes a provider
	Init() error
}

// Factory is a factory for tile providers
type Factory interface {
	// Gets provider type key
	Type() string
	// Creates a provider
	Create(config Config, manager Manager) (Provider, error)
}

// RegisterFactory registers a provider factory
func RegisterFactory(factory Factory) {
	factories[factory.Type()] = factory
}

// GetFactory selects a provider factory by its key
func GetFactory(key string) Factory {
	factory, exists := factories[key]
	if !exists {
		return nil
	}
	return factory
}

var factories = make(map[string]Factory)
