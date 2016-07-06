package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	dash "github.com/itglobal/dashboard/api"
)

type configJSON struct {
	Theme     themeJSON     `json:"theme"`
	Providers []interface{} `json:"providers"`
}

type themeJSON struct {
	Style  string `json:"style"`
	Colors string `json:"colors"`
}

type configReader struct {
	data map[string]interface{}
}

var configFileName string = defaultConfigFileName

const defaultConfigFileName = "dash.json"

func readConfig() (*configJSON, error) {
	configBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var config configJSON
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func newConfigReader(data interface{}) (dash.Config, error) {
	switch data := data.(type) {
	case map[string]interface{}:
		return &configReader{data: data}, nil
	default:
		return nil, errors.New("Bad JSON data")
	}
}

func (c *configReader) GetString(key string) (string, error) {
	value, exists := c.data[key]
	if !exists {
		return "", dash.ErrNoSuchKey
	}

	return value.(string), nil
}

func (c *configReader) GetStringOrDefault(key string, def string) string {
	value, err := c.GetString(key)
	if err != nil {
		return def
	}

	return value
}

func (c *configReader) GetInt32(key string) (int32, error) {
	value, exists := c.data[key]
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

func (c *configReader) GetInt32OrDefault(key string, def int32) int32 {
	value, err := c.GetInt32(key)
	if err != nil {
		return def
	}

	return value
}
