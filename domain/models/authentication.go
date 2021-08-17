package models

import "time"

// Authentication is used to store data inside JWT
type Authentication struct {
	AccountID  DatabaseID `json:"account_id"`
	Expiration time.Time  `json:"expiration"`
}
