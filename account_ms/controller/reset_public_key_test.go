package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/account_ms/models"

	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_ResetPublicKey(t *testing.T) {
	ctx := context.Background()
	svc, err := InitializeFromEnv(ctx)
	require.NoError(t, err)

	col := svc.Persistence.GetCollection(accountCollection)
	defer t.Run("cleanup", func(t *testing.T) {
		_, err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})

	oldKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "should create old private key")

	account := dModels.Account{
		OwnershipKey: dModels.OwnershipKey("123"),
		Addresses: []dModels.Address{
			dModels.Address("alysson@letter.me"),
			dModels.Address("alysson_2@letter.me"),
		},
		PublicKey: *oldKey.GetPublicKey(),
	}

	_, err = col.Create(ctx, account)
	require.NoError(t, err, "should create account")

	newKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "should create new private key")

	newPublicKey := newKey.GetPublicKey()

	err = svc.ResetPublicKey(ctx, models.ResetPublicKeyRequest{
		OwnershipKey: account.OwnershipKey,
		PublicKey:    *newPublicKey,
	})
	require.NoError(t, err, "should change public key")

	dbAccount := new(dModels.Account)
	err = col.First(ctx, filter{}, &dbAccount)
	require.NoError(t, err, "should find account")

	require.Equal(t, *newPublicKey, dbAccount.PublicKey, "public key should be updated")
}
