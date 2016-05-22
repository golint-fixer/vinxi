package manager

import (
	"encoding/json"
)

// Config encapsulates a map of polymorfic values and provides
// useful method to lookup and retrieve fields by key.
//
// Config is designed to be read-only and thread-safe.
type Config map[string]interface{}

// GetString lookups a config value by key and returns it as string.
func (c Config) GetString(key string) string {
	return c[key].(string)
}

// GetBool lookups a config value by key and returns it as bool.
func (c Config) GetBool(key string) bool {
	return c[key].(bool)
}

// GetInt lookups a config value by key and returns it as int.
func (c Config) GetInt(key string) int {
	return c[key].(int)
}

// GetInt64 lookups a config value by key and returns it as int64.
func (c Config) GetInt64(key string) int64 {
	return c[key].(int64)
}

// GetFloat lookups a config value by key and returns it as float64.
func (c Config) GetFloat(key string) float64 {
	return c[key].(float64)
}

// Get lookups a config value by key and returns it.
func (c Config) Get(key string) interface{} {
	return c[key]
}

// Exists checks is config field is present.
func (c Config) Exists(key string) bool {
	return c[key] != nil
}

// JSON serializes the config fields a JSON.
func (c Config) JSON() ([]byte, error) {
	return json.Marshal(c)
}
