package controller

import (
	"context"
	"testing"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/require"
)

func Test_Authenticate(t *testing.T) {
	ctx := context.Background()
	svc, err := InitializeFromEnv(ctx)
	require.NoError(t, err)

	col := svc.Persistence.GetCollection(accountCollection)

	defer t.Run("cleanup", func(t *testing.T) {
		_, err := col.Delete(ctx, filter{})
		require.NoError(t, err, "should clear collection")
	})

	clientKey, err := cryptography.NewPrivateKey(2048)
	require.NoError(t, err, "private key should be created")

	account := dModels.Account{
		Addresses: []dModels.Address{dModels.Address("alysson@letter.me")},
		PublicKey: *clientKey.GetPublicKey(),
	}

	_, err = col.Create(ctx, account)
	require.NoError(t, err, "should create account")

	t.Run("should authenticate", func(t *testing.T) {
		encryptedJWT, err := svc.Authenticate(ctx, account.Addresses[0])
		require.NoError(t, err)
		require.NotNil(t, encryptedJWT)

		var jwtToken string
		err = svc.decrypt(clientKey, encryptedJWT, &jwtToken)
		require.NoError(t, err, "should decrypt jwt token client side")

		claims := new(dModels.TokenClaims)
		err = svc.Authenticator.ReadToken(jwtToken, claims)
		require.NoError(t, err, "should read claims without errors")
		require.Equal(t, account.Addresses[0], claims.Address, "claim address should match")
		require.GreaterOrEqual(t, time.Now().Add(time.Hour).Unix(), claims.ExpiresAt, "claim expiry date should be expired after %s", time.Hour)
	})
}
