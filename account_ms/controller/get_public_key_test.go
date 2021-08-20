package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain"
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

	email := domain.Address("alysson@letter.me")
	privateKey, err := domain.NewPrivateKey(2048)
	require.NoError(t, err)

	publicKey := privateKey.GetPublicKey()
	account := domain.Account{
		OwnershipKey: "123",
		PublicKey:    *publicKey,
		Addresses: []domain.Address{
			email,
			domain.Address("alysson_2@letter.me"),
		},
	}

	t.Run("setup", func(t *testing.T) {
		_, err := col.Create(ctx, account)
		require.NoError(t, err, "should create account")

		account.Addresses = []domain.Address{"priscila@gmail.com"}
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
