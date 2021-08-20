package domain

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

// CryptoConfig is used to fetch configurations related to
type CryptoConfig struct {
	Configs          map[AlgorithmName]AlgorithmConfiguration `json:"configs"`
	DefaultAlgorithm AlgorithmName                            `json:"default_algorithm"`
}

func (c CryptoConfig) Validate() error {
	var errList []error
	if len(c.Configs) == 0 {
		errList = append(errList, newEmptyFieldError("configs"))
	}

	if c.DefaultAlgorithm == "" {
		errList = append(errList, newEmptyFieldError("session_name"))
	}
	if len(errList) > 0 {
		return newInvalidConfigError(c, errList)
	}
	return nil
}

// CryptoRouter is responsible for routing custom deserialization for cryptographic algorithms,
// and encryption for interfaces.
type CryptoRouter struct {
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

func NewCryptoRouter(c *CryptoConfig) (*CryptoRouter, error) {
	router := &CryptoRouter{
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

func (r *CryptoRouter) addRSA_OAEP(cypher []byte, hash hash.Hash) {
	r.Algorithms[RSA_OAEP] = rsa_oaep{
		cypher: cypher,
		hash:   hash,
	}
}

func (r *CryptoRouter) Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error {
	if algorithm, ok := r.Algorithms[b.Algorithm]; ok {
		return algorithm.Decrypt(k, b, dst)
	} else {
		return fmt.Errorf("handler for '%s' not found", b.Algorithm)
	}
}

func (r *CryptoRouter) Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error) {
	if algorithm, ok := r.Algorithms[r.defaultAlgorithm]; ok {
		return algorithm.Encrypt(k, src)
	} else {
		return nil, fmt.Errorf("handler for '%s' not found", algorithm)
	}
}

func (r *CryptoRouter) EncryptAlgorithm(k *PublicKey, src interface{}, algorithm AlgorithmName) (*EncryptedBuffer, error) {
	if algorithm, ok := r.Algorithms[algorithm]; ok {
		return algorithm.Encrypt(k, src)
	} else {
		return nil, fmt.Errorf("handler for '%s' not found", algorithm)
	}
}