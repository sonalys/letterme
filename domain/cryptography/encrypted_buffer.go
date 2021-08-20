package cryptography

// AlgorithmName represents all compatible encryption algorithms.
type AlgorithmName string

const (
	RSA_OAEP AlgorithmName = "rsa_oaep"
)

// HashFunc represents used hash function.
type HashFunc string

const (
	SHA256 HashFunc = "sha-256"
)

// EncryptedBuffer is used to keep compatibility in case of encryption algorithm changes to the system.
type EncryptedBuffer struct {
	Buffer    []byte        `json:"buffer"`
	Algorithm AlgorithmName `json:"algorithm"`
	Hash      HashFunc      `json:"hash"`
}
