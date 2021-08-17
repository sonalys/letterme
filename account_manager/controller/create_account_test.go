package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	ctx := context.Background()
	db := createPersistence(ctx, t)
	svc, err := NewService(ctx, &Dependencies{
		Persistence: db,
	})
	require.NoError(t, err)
	col := svc.Persistence.GetCollection("account")

	account := models.Account{
		Addresses: []models.Address{
			models.Address("alysson@letter.me"),
			models.Address("alysson_2@letter.me"),
		},
	}

	t.Run("should create account", func(t *testing.T) {
		token, err := svc.CreateAccount(ctx, account)
		account.OwnershipKey = token
		require.NoError(t, err, "should create account")
		require.NotEmpty(t, token, "ownershipToken should not be empty")
	})

	t.Run("dbAccount verification", func(t *testing.T) {
		var dbAccount models.Account
		err := col.First(ctx, models.Account{
			OwnershipKey: account.OwnershipKey,
		}, &dbAccount)
		require.NoError(t, err, "should create account")
		require.NotNil(t, dbAccount, "account should return from db")
		require.Equal(t, account.Addresses[:1], dbAccount.Addresses, "dbAccount should have only 1 email")
		require.Equal(t, account.OwnershipKey, dbAccount.OwnershipKey, "dbAccount should have same ownership as returned")
	})

	t.Run("should not create duplicate account", func(t *testing.T) {
		token, err := svc.CreateAccount(ctx, account)
		require.Error(t, err, "should not create account")
		require.Empty(t, token, "ownershipToken should be empty")
	})

	defer t.Run("cleanup", func(t *testing.T) {
		err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})
}
