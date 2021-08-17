package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load config from env into dst.
// dst must be a pointer.
func LoadFromEnv(key string, dst interface{}) error {
	if val, ok := os.LookupEnv(key); ok {
		if err := json.Unmarshal([]byte(val), dst); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error loading config from env: '%s' not found", key)
}
