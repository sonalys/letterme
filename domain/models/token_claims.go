package models

// TokenClaims are the letter.me customized jwt token claims.
// Need to have address and expiry date.
type TokenClaims struct {
	ExpiresAt int64
	Address   Address
}

func (c TokenClaims) Valid() error {
	return nil
}
