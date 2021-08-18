package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_ResetPublicKey(t *testing.T) {
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

	oldKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "should create old private key")

	account := models.Account{
		OwnershipKey: models.OwnershipKey("123"),
		Addresses: []models.Address{
			models.Address("alysson@letter.me"),
			models.Address("alysson_2@letter.me"),
		},
		PublicKey: *oldKey.GetPublicKey(),
	}

	_, err = col.Create(ctx, account)
	require.NoError(t, err, "should create account")

	newKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "should create new private key")

	newPublicKey := newKey.GetPublicKey()

	err = svc.ResetPublicKey(ctx, account.OwnershipKey, *newPublicKey)
	require.NoError(t, err, "should change public key")

	dbAccount := new(models.Account)
	err = col.First(ctx, filter{}, &dbAccount)
	require.NoError(t, err, "should find account")

	require.Equal(t, *newPublicKey, dbAccount.PublicKey, "public key should be updated")
}
