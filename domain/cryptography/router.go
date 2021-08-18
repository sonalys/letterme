package cryptography

import (
	"crypto"
	"fmt"
	"hash"
)

const CRYPTO_CYPHER_ENV = "LM_CRYPTO_CYPHER"

type AlgorithmConfiguration struct {
	Cypher []byte   `json:"cypher"`
	Hash   HashFunc `json:"hash"`
}

// Configuration is used to fetch configurations related to cryptography.
type Configuration struct {
	Configs          map[AlgorithmName]AlgorithmConfiguration `json:"configs"`
	DefaultAlgorithm AlgorithmName                            `json:"default_algorithm"`
}

// Router is responsible for routing custom deserialization for cryptographic algorithms,
// and encryption for interfaces.
type Router struct {
	defaultAlgorithm AlgorithmName
	Algorithms       map[AlgorithmName]EncryptionAlgorithm
}

func stringToHash(s HashFunc) (hash.Hash, error) {
	switch s {
	case sha256:
		return crypto.SHA256.New(), nil
	default:
		return nil, fmt.Errorf("hash not found: '%s'", s)
	}
}

func NewRouter(c *Configuration) (*Router, error) {
	router := &Router{
		defaultAlgorithm: c.DefaultAlgorithm,
		Algorithms:       make(map[AlgorithmName]EncryptionAlgorithm),
	}

	for algorithm, config := range c.Configs {
		switch algorithm {
		case RSA_OAEP:
			hashAlg, err := stringToHash(config.Hash)
			if err != nil {
				return nil, err
			}
			router.addRSA_OAEP(config.Cypher, hashAlg)
		}
	}
	return router, nil
}

func (r *Router) addRSA_OAEP(cypher []byte, hash hash.Hash) {
	r.Algorithms[RSA_OAEP] = rsa_oaep{
		cypher: cypher,
		hash:   hash,
	}
}

func (r *Router) Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error {
	if algorithm, ok := r.Algorithms[b.Algorithm]; ok {
		return algorithm.Decrypt(k, b, dst)
	} else {
		return fmt.Errorf("handler for '%s' not found", b.Algorithm)
	}
}

func (r *Router) Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error) {
	if algorithm, ok := r.Algorithms[r.defaultAlgorithm]; ok {
		return algorithm.Encrypt(k, src)
	} else {
		return nil, fmt.Errorf("handler for '%s' not found", algorithm)
	}
}

func (r *Router) EncryptAlgorithm(k *PublicKey, src interface{}, algorithm AlgorithmName) (*EncryptedBuffer, error) {
	if algorithm, ok := r.Algorithms[algorithm]; ok {
		return algorithm.Encrypt(k, src)
	} else {
		return nil, fmt.Errorf("handler for '%s' not found", algorithm)
	}
}
