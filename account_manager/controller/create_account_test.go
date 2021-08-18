package controller

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	ctx := context.Background()
	db := createPersistence(ctx, t)

	router, err := cryptography.NewRouter(&cryptography.Configuration{
		DefaultAlgorithm: cryptography.RSA_OAEP,
		Configs: map[cryptography.AlgorithmName]cryptography.AlgorithmConfiguration{
			cryptography.RSA_OAEP: cryptography.AlgorithmConfiguration{
				Cypher: []byte("123"),
				Hash:   "sha-256",
			},
		},
	})
	require.NoError(t, err)

	svc, err := NewService(ctx, &Dependencies{
		Persistence:         db,
		CryptographicRouter: router,
	})
	require.NoError(t, err)
	col := svc.Persistence.GetCollection("account")

	defer t.Run("cleanup", func(t *testing.T) {
		err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})

	pk, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "private key should be created")

	account := models.Account{
		Addresses: []models.Address{
			models.Address("alysson@letter.me"),
			models.Address("alysson_2@letter.me"),
		},
		OwnershipKey: models.OwnershipKey("123"),
		TTL:          time.Microsecond,
		PublicKey:    cryptography.PublicKey(pk.PublicKey),
		DeviceCount:  8,
		ID:           "123",
	}

	var encryptedToken *cryptography.EncryptedBuffer
	t.Run("should create account", func(t *testing.T) {
		encryptedToken, err = svc.CreateAccount(ctx, account)
		require.NoError(t, err, "should create account")
		require.NotEmpty(t, encryptedToken, "ownershipToken should not be empty")
	})

	decryptedOwnershipKey := new(models.OwnershipKey)
	err = svc.decrypt(pk, encryptedToken, decryptedOwnershipKey)
	require.NoError(t, err, "ownership_key should be decrypted")

	t.Run("dbAccount verification", func(t *testing.T) {
		var dbAccount models.Account
		err := col.First(ctx, models.Account{
			OwnershipKey: *decryptedOwnershipKey,
		}, &dbAccount)
		require.NoError(t, err, "should create account")
		require.NotNil(t, dbAccount, "account should return from db")

		require.Equal(t, account.Addresses[:1], dbAccount.Addresses, "dbAccount should have only 1 email")

		require.Equal(t, *decryptedOwnershipKey, dbAccount.OwnershipKey, "dbAccount should have same ownership as returned")
		require.NotEqual(t, account.OwnershipKey, dbAccount.OwnershipKey, "dbAccount should not have user defined ownership")

		require.NotEqual(t, account.TTL, dbAccount.TTL, "dbAccount should not have user defined ttl")
		require.NotEqual(t, account.DeviceCount, dbAccount.DeviceCount, "dbAccount should not have user defined deviceCount")

		require.NotEqual(t, account.ID, dbAccount.ID, "dbAccount should not have user defined id")
	})

	t.Run("should not create duplicate account", func(t *testing.T) {
		token, err := svc.CreateAccount(ctx, account)
		require.Error(t, err, "should not create account")
		require.Empty(t, token, "ownershipToken should be empty")
	})
}
