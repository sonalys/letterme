package domain

import (
	"crypto"
	"fmt"
	"hash"
)

const CRYPTO_CYPHER_ENV = "LM_CRYPTO_CONFIG"

type AlgorithmConfiguration struct {
	Cypher []byte   `json:"cypher"`
	Hash   HashFunc `json:"hash"`
}

func (c AlgorithmConfiguration) Validate() error {
	var errList []error
	if len(c.Cypher) == 0 {
		errList = append(errList, newEmptyFieldError("cypher"))
	}

	if c.Hash == "" {
		errList = append(errList, newEmptyFieldError("hash"))
	}

	if len(errList) > 0 {
		return newInvalidConfigError(c, errList)
	}
	return nil
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

	for key := range c.Configs {
		if key != RSA_OAEP {
			errList = append(errList, newInvalidFieldError(
				fmt.Sprintf("configs[%s]", key),
				fmt.Errorf("'%s' is not a valid hash name", key),
			))
		}
		if err := c.Configs[key].Validate(); err != nil {
			errList = append(
				errList,
				newInvalidFieldError(
					fmt.Sprintf("configs[%s]", key),
					err,
				),
			)
		}
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
	case SHA256:
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
		// nolint // will have more algorithms in the future, we will keep switch-case
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

func (r *CryptoRouter) addRSA_OAEP(cypher []byte, hashAlgorithm hash.Hash) {
	r.Algorithms[RSA_OAEP] = rsa_oaep{
		cypher: cypher,
		hash:   hashAlgorithm,
	}
}

func (r *CryptoRouter) Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error {
	algorithm, ok := r.Algorithms[b.Algorithm]
	if !ok {
		return fmt.Errorf("handler for '%s' not found", b.Algorithm)
	}
	return algorithm.Decrypt(k, b, dst)
}

func (r *CryptoRouter) Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error) {
	algorithm, ok := r.Algorithms[r.defaultAlgorithm]
	if !ok {
		return nil, fmt.Errorf("handler for '%s' not found", r.defaultAlgorithm)
	}
	return algorithm.Encrypt(k, src)
}

func (r *CryptoRouter) EncryptAlgorithm(k *PublicKey, src interface{}, name AlgorithmName) (*EncryptedBuffer, error) {
	algorithm, ok := r.Algorithms[name]
	if !ok {
		return nil, fmt.Errorf("handler for '%s' not found", name)
	}
	return algorithm.Encrypt(k, src)
}
