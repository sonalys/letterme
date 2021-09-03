package smtp

import (
	"errors"

	"github.com/sonalys/letterme/domain/models"
)

// parseEmailAddress parses <email@example.com> to models.Address.
//
// TODO: maybe should be in domain?
func parseEmailAddress(buf []byte) (*models.Address, error) {
	size := len(buf)
	if size < 4 || size > 253 {
		return nil, errors.New("address size must be between 4 and 253")
	}

	strip := buf[1 : len(buf)-1]
	addr := models.Address(strip)
	return &addr, addr.Validate()
}
