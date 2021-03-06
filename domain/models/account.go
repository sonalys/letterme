package models

import (
	"crypto/rsa"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
)

// Account represents an user account,
// Can have many addresses and devices, it won't have any password, just the public key.
//
// Must be linked to a unique code, which can be used to delete the account or re-upload a new private certificate later.
type Account struct {
	ID DatabaseID `json:"id" bson:"_id,omitempty"`
	// Addresses are the many email addresses possessed by this account.
	Addresses []Address `json:"addresses" bson:"addresses,omitempty"`
	// PublicKey is used to encrypt all data sent to this user.
	PublicKey cryptography.PublicKey `json:"public_key" bson:"publicKey,omitempty"`
	// Ownershipkey is used to re-upload a new private key to recover the used addresses, all the previous data is lost however.
	// It must be used only to this mean, for authentication, use JWT.
	OwnershipKey OwnershipKey `json:"ownership_key" bson:"ownershipKey,omitempty"`
	// DeviceCount is used to keep emails and attachments into backend for multiple read confirmations before deleting it.
	DeviceCount uint8 `json:"device_count" bson:"deviceCount,omitempty"`
	// TTL informs how many time messages sent to this user will persist,
	// this information will be fetched and inserted into email.valid_until
	TTL time.Duration `json:"ttl" bson:"ttl,omitempty"`
}

// AccountAddressInfo is used to fetch information about a given address.
type AccountAddressInfo struct {
	Address       Address `json:"address"`
	rsa.PublicKey `json:"public_key"`
}
