package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_DeleteAccount(t *testing.T) {
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

	account := models.Account{
		OwnershipKey: "123",
		Addresses: []models.Address{
			models.Address("alysson@letter.me"),
			models.Address("alysson_2@letter.me"),
		},
	}

	t.Run("setup", func(t *testing.T) {
		_, err := col.Create(ctx, account)
		require.NoError(t, err, "should create account")
	})

	t.Run("should delete account", func(t *testing.T) {
		err := svc.DeleteAccount(ctx, account.OwnershipKey)
		require.NoError(t, err, "should delete account")
	})
}
