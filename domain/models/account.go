package models

import (
	"crypto/rsa"
	"time"
)

// Account represents an user account,
// Can have many addresses and devices, it won't have any password, just the public key.
//
// Must be linked to a unique code, which can be used to delete the account or re-upload a new private certificate later.
type Account struct {
	ID DatabaseID `json:"id"`
	// Addresses are the many email addresses posessed by this account.
	Addresses []Address `json:"addresses"`
	// PublicKey is used to encrypt all data sent to this user.
	PublicKey PublicKey `json:"public_key"`
	// Ownershipkey is used to re-upload a new private key to recover the used addresses, all the previous data is lost however.
	// It must be used only to this mean, for authentication, use JWT.
	OwnershipKey string `json:"ownership_key"`
	// DeviceCount is used to keep emails and attachments into backend for multiple read confirmations before deleting it.
	DeviceCount uint8 `json:"device_count"`
	// TTL informs how many time messages sent to this user will persist,
	// this information will be fetched and inserted into email.valid_until
	TTL time.Duration `json:"ttl"`
}

// AccountAddressInfo is used to fetch information about a given address.
type AccountAddressInfo struct {
	Address       Address `json:"address"`
	rsa.PublicKey `json:"public_key"`
}
