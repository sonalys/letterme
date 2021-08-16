package models

// EncryptionAlgorithm represents all compatible encryption algorithms.
type EncryptionAlgorithm string

const (
	RSA_OAEP EncryptionAlgorithm = "rsa_oaep"
)

// EncryptedBuffer is used to keep compatibility in case of encryption algorithm changes to the system.
type EncryptedBuffer struct {
	Buffer    []byte              `json:"buffer"`
	Algorithm EncryptionAlgorithm `json:"algorithm"`
}
