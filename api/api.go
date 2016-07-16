package api

import "errors"

// ItemStatus is a status of a dashboard item
type ItemStatus int

const (
	// StatusUnknown indicates an unknown state of a dashboard item
	StatusUnknown ItemStatus = iota
	// StatusGood indicates GOOD state of a dashboard item
	StatusGood
	// StatusPending indicates a pending state of a dashboard item
	StatusPending
	// StatusBad indicates BAD state of a dashboard item
	StatusBad
)

// Item is a dashboard item
type Item struct {
	// A provider key
	ProviderKey string `json:"provider"`

	// A item key
	Key string `json:"key"`

	// Item title
	Name string `json:"name"`

	// Item status
	Status ItemStatus `json:"status"`

	// Item status text
	StatusText string `json:"text"`

	// Progress value. This field is applicable only to items with status of StatusPending.
	// If this field has value of NoProgress then progress value is not used
	Progress int `json:"progress"`
}

// ItemList is a list of items with version number
type ItemList struct {
	// Version is a number that is incremented with every update
	Version uint `json:"version"`
	// Dashboard items
	Items []Item `json:"items"`
}

// NoProgress is a special value for field Progress of type Item
// that indicates no progress value
const NoProgress int = -1

// Callback is a callback function for Provider
type Callback func(Provider, []*Item)

// Provider defines methods for sources that provide dashboard items
type Provider interface {
	// Gets a provider key (see field ProviderKey of Item type)
	Key() string
}

// ErrNoSuchKey means that requested key is not found
var ErrNoSuchKey = errors.New("No such key in config")

// Config is a callback to read config parameters
type Config interface {
	// Read a string config parameter
	GetString(key string) (string, error)

	// Read a string config parameter or use default value
	GetStringOrDefault(key string, def string) string

	// Read an int32 config parameter
	GetInt32(key string) (int32, error)

	// Read an int32 config parameter or use default value
	GetInt32OrDefault(key string, def int32) int32
}

// Factory is a factory function for dash providers
type Factory func(Config, Callback) (Provider, error)

// RegisterFactory registers a provider factory
func RegisterFactory(providerType string, factory Factory) {
	providers[providerType] = factory
}

// GetFactory selects a provider factory by its type
func GetFactory(providerType string) Factory {
	factory, exists := providers[providerType]
	if !exists {
		return nil
	}
	return factory
}

var providers = make(map[string]Factory)
