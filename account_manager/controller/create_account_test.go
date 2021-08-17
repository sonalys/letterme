package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/account_manager/interfaces"
	"github.com/sonalys/letterme/account_manager/persistence"
	"github.com/sonalys/letterme/account_manager/utils"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func createPersistence(ctx context.Context, t *testing.T) interfaces.Persistence {
	var cfg persistence.Configuration
	if err := utils.LoadFromEnv(persistence.MONGO_ENV, &cfg); err != nil {
		require.Fail(t, err.Error())
	}

	mongo, err := persistence.NewMongo(ctx, &cfg)
	require.NoError(t, err, "should create without errors")

	return mongo
}

func Test_CreateAccount(t *testing.T) {
	ctx := context.Background()
	db := createPersistence(ctx, t)

	svc, err := NewService(ctx, &Dependencies{
		Persistence: db,
	})
	require.NoError(t, err)

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
		dbAccount, err := svc.GetAccount(ctx, account.OwnershipKey)
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

	defer func() {
		err := svc.DeleteAccount(ctx, account.OwnershipKey)
		require.NoError(t, err, "should create account")
	}()
}
