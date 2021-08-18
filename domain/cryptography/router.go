package cryptography

import (
	"fmt"
	"hash"
)

type Router struct {
	Algorithms map[AlgorithmName]EncryptionAlgorithm
}

func NewRouter() *Router {
	return &Router{
		Algorithms: make(map[AlgorithmName]EncryptionAlgorithm),
	}
}

// AddRSA_OAEP configures a new rsa-oaep handler for this router.
func (r *Router) AddRSA_OAEP(cypher []byte, hash hash.Hash) {
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

func (r *Router) Encrypt(k *PublicKey, algorithm AlgorithmName, src interface{}) (*EncryptedBuffer, error) {
	if algorithm, ok := r.Algorithms[algorithm]; ok {
		return algorithm.Encrypt(k, src)
	} else {
		return nil, fmt.Errorf("handler for '%s' not found", algorithm)
	}
}
