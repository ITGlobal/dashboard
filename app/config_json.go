package app

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/kpango/glg"
)

type Config struct {
	json map[string]interface{}
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("Unable to open config file \"%s\". %s", path, err)
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Unable to read config file \"%s\". %s", path, err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Errorf("Unable to parse config file \"%s\". %s", path, err)
		return nil, err
	}

	log.Debugf("Using config file \"%s\"", path)

	return &Config{data}, nil
}

func (c *Config) listProviders() []*providerConfig {
	value, exists := c.json["providers"]
	if !exists {
		log.Error("Config file doesn't have a required key \"providers\"")
		return nil
	}

	switch jsons := value.(type) {
	case []interface{}:
		var arr []*providerConfig
		for i, json := range jsons {
			config, ok := newProviderConfig(i, json)
			if ok {
				arr = append(arr, config)
			}
		}
		return arr

	default:
		log.Error("Bad config: \"providers\" is not an array")
		return nil
	}
}

type providerConfig struct {
	Type  string
	index int
	json  map[string]interface{}
}

func newProviderConfig(i int, json interface{}) (*providerConfig, bool) {
	jsonMap, ok := json.(map[string]interface{})
	if !ok {
		log.Errorf("Bad config: \"providers[%d]\" is not an object", i)
		return nil, false
	}

	value, exists := jsonMap["type"]
	if !exists {
		log.Errorf("Bad config: \"providers[%d].type\" is missing", i)
		return nil, false
	}

	var providerType string

	switch v := value.(type) {
	case string:
		if v == "" {
			log.Errorf("Bad config: \"providers[%d].type\" is empty", i)
			return nil, false
		}

		providerType = v

	default:
		log.Errorf("Bad config: \"providers[%d].type\" is not a string", i)
		return nil, false
	}

	enabledRaw, ok := jsonMap["enabled"]
	if ok {
		enabled, ok := enabledRaw.(bool)
		if ok && !enabled {
			log.Warnf("Provider \"providers[%d]\" is disabled", i)
			return nil, false
		}
	}

	return &providerConfig{providerType, i, jsonMap}, true
}

// Read a string config parameter
func (c *providerConfig) GetString(key string) (string, bool) {
	raw, exists := c.json[key]
	if !exists {
		log.Warnf("Config: \"providers[%d].%s\" is missing", c.index, key)
		return "", false
	}

	str, ok := raw.(string)
	if !ok {
		log.Warnf("Config: \"providers[%d].%s\" is not a string", c.index, key)
		return "", false
	}

	return str, true
}

// Readss a string array config parameter
func (c *providerConfig) GetStringArray(key string) ([]string, bool) {
	raw, exists := c.json[key]
	if !exists {
		log.Warnf("Config: \"providers[%d].%s\" is missing", c.index, key)
		return nil, false
	}

	arr, ok := raw.([]interface{})
	if !ok {
		log.Warnf("Config: \"providers[%d].%s\" is not an array", c.index, key)
		return nil, false
	}

	var result []string
	for i, item := range arr {
		str, ok := item.(string)
		if !ok {
			log.Warnf("Config: \"providers[%d].%s[%d]\" is not a string", c.index, key, i)
			continue
		}

		result = append(result, str)
	}

	return result, true
}

// Reads an int config parameter
func (c *providerConfig) GetInteger(key string) (int, bool) {
	raw, exists := c.json[key]
	if !exists {
		log.Warnf("Config: \"providers[%d].%s\" is missing", c.index, key)
		return 0, false
	}

	integer, ok := raw.(int)
	if !ok {
		log.Warnf("Config: \"providers[%d].%s\" is not a number", c.index, key)
		return 0, false
	}

	return integer, true
}
