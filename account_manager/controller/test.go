package controller

import (
	"context"
	"testing"

	"github.com/sonalys/letterme/account_manager/interfaces"
	"github.com/sonalys/letterme/account_manager/persistence"
	"github.com/sonalys/letterme/account_manager/utils"
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
