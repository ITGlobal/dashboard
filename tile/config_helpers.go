package tile

// Config defines methods to read tile config parameters
type Config interface {
	// Read a string config parameter
	GetString(key string) (string, bool)

	// Readss a string array config parameter
	GetStringArray(key string) ([]string, bool)

	// Reads an integer config parameter
	GetInteger(key string) (int, bool)
}

// GetStringDefault reads a string config parameter or returns default value
func GetStringDefault(c Config, key string, def string) string {
	value, ok := c.GetString(key)
	if !ok {
		return def
	}

	return value
}

// GetIntegerDefault reads an integer config parameter or returns default value
func GetIntegerDefault(c Config, key string, def int) int {
	value, ok := c.GetInteger(key)
	if !ok {
		return def
	}

	return value
}
