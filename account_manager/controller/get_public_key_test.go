package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_GetPublicKey(t *testing.T) {
	ctx := context.Background()
	db := createPersistence(ctx, t)
	svc, err := NewService(ctx, &Dependencies{
		Persistence: db,
	})
	require.NoError(t, err)
	col := svc.Persistence.GetCollection("account")
	defer t.Run("cleanup", func(t *testing.T) {
		_, err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})

	email := models.Address("alysson@letter.me")
	privateKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err)

	publicKey := privateKey.GetPublicKey()
	account := models.Account{
		OwnershipKey: "123",
		PublicKey:    *publicKey,
		Addresses: []models.Address{
			email,
			models.Address("alysson_2@letter.me"),
		},
	}

	t.Run("setup", func(t *testing.T) {
		_, err := col.Create(ctx, account)
		require.NoError(t, err, "should create account")

		account.Addresses = []models.Address{"priscila@gmail.com"}
		_, err = col.Create(ctx, account)
		require.NoError(t, err, "should create dummy account")
	})

	t.Run("should get account", func(t *testing.T) {
		dbPublicKey, err := svc.GetPublicKey(ctx, email)
		require.NoError(t, err, "should find account")
		require.NotNil(t, dbPublicKey, "should return publicKey")
		require.Equal(t, publicKey, dbPublicKey, "should return account")
	})
}
