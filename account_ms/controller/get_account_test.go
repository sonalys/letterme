package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain"
	"github.com/stretchr/testify/require"
)

func Test_GetAccount(t *testing.T) {
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

	account := domain.Account{
		OwnershipKey: "123",
		Addresses: []domain.Address{
			domain.Address("alysson@letter.me"),
			domain.Address("alysson_2@letter.me"),
		},
	}

	t.Run("setup", func(t *testing.T) {
		_, err := col.Create(ctx, account)
		require.NoError(t, err, "should create account")

		account.OwnershipKey = "456"
		_, err = col.Create(ctx, account)
		require.NoError(t, err, "should create dummy account")
	})

	t.Run("should get account", func(t *testing.T) {
		dbAccount, err := svc.GetAccount(ctx, account.OwnershipKey)
		require.NoError(t, err, "should delete account")
		require.NotNil(t, dbAccount, "should return account")
		require.Equal(t, account.OwnershipKey, dbAccount.OwnershipKey, "should return account")
	})
}
