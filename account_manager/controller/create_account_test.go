package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/account_manager/models"

	domain "github.com/sonalys/letterme/domain"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	ctx := context.Background()
	db := createPersistence(ctx, t)

	router, err := domain.NewCryptoRouter(&domain.CryptoConfig{
		DefaultAlgorithm: domain.RSA_OAEP,
		Configs: map[domain.AlgorithmName]domain.AlgorithmConfiguration{
			domain.RSA_OAEP: {
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
		_, err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})

	pk, err := domain.NewPrivateKey(2048)
	require.NoError(t, err, "private key should be created")

	account := models.CreateAccountRequest{
		Address:   "alysson@letter.me",
		PublicKey: *pk.GetPublicKey(),
	}

	var encryptedToken *domain.EncryptedBuffer
	t.Run("should create account", func(t *testing.T) {
		encryptedToken, err = svc.CreateAccount(ctx, account)
		require.NoError(t, err, "should create account")
		require.NotEmpty(t, encryptedToken, "ownershipToken should not be empty")
	})

	decryptedOwnershipKey := new(domain.OwnershipKey)
	err = svc.decrypt(pk, encryptedToken, decryptedOwnershipKey)
	require.NoError(t, err, "ownership_key should be decrypted")

	t.Run("dbAccount verification", func(t *testing.T) {
		var dbAccount domain.Account
		err := col.First(ctx, domain.Account{
			OwnershipKey: *decryptedOwnershipKey,
		}, &dbAccount)
		require.NoError(t, err, "should create account")
		require.NotNil(t, dbAccount, "account should return from db")

		require.Equal(t, []domain.Address{account.Address}, dbAccount.Addresses, "dbAccount should have only 1 email")

		require.Equal(t, *decryptedOwnershipKey, dbAccount.OwnershipKey, "dbAccount should have same ownership as returned")

		require.Equal(t, uint8(1), dbAccount.DeviceCount, "dbAccount should not have user defined deviceCount")
	})

	t.Run("should not create duplicate account", func(t *testing.T) {
		token, err := svc.CreateAccount(ctx, account)
		require.Error(t, err, "should not create account")
		require.Empty(t, token, "ownershipToken should be empty")
	})

	t.Run("should not create invalid email", func(t *testing.T) {
		account.Address = domain.Address("alysson¨@letter.me")
		token, err := svc.CreateAccount(ctx, account)
		require.Error(t, err, "should not create account")
		require.Empty(t, token, "ownershipToken should be empty")
	})

	t.Run("should not create email from other domains", func(t *testing.T) {
		account.Address = domain.Address("alysson@gmail.com")
		token, err := svc.CreateAccount(ctx, account)
		require.Error(t, err, "should not create account")
		require.Empty(t, token, "ownershipToken should be empty")
	})
}
