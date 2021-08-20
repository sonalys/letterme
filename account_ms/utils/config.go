package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFromEnv loads config from env into dst.
// dst must be a pointer.
func LoadFromEnv(key string, dst interface{}) error {
	if val, ok := os.LookupEnv(key); ok {
		return json.Unmarshal([]byte(val), dst)
	}
	return fmt.Errorf("error loading config from env: '%s' not found", key)
}
