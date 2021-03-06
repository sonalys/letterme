package utils

import (
	"encoding/json"
	"fmt"
	"os"

	dModels "github.com/sonalys/letterme/domain/models"
)

// LoadFromEnv loads config from env into dst.
// dst must be a pointer.
func LoadFromEnv(key string, dst interface{}) error {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fmt.Errorf("error loading config from env: '%s' not found", key)
	}

	if err := json.Unmarshal([]byte(val), dst); err != nil {
		return err
	}

	if validatable, ok := dst.(dModels.Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}
