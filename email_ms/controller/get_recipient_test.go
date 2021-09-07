package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_VerifyEmailExistence(t *testing.T) {
	ctx := context.Background()

	testCases := []testCase{
		{
			name: "fail to publish",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.router.On("Communicate", messaging.AccountMS, mock.Anything, mock.Anything).
					Return(errors.New("foo/bar"))
				exists, pk, err := s.getRecipient(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.Nil(t, pk)
				assert.False(t, exists)
			},
		},
		{
			name: "all ok",
			run: func(t *testing.T, s *Service, th *testHandler) {
				privKey, err := cryptography.NewPrivateKey(2048)
				require.NoError(t, err)
				pubKey := privKey.GetPublicKey()
				th.router.On("Communicate", messaging.AccountMS, mock.Anything, mock.Anything).
					Run(mockSetDST(2, contracts.CheckEmailResponse{Exists: true, PublicKey: pubKey})).
					Return(nil)

				exists, pk, err := s.getRecipient(ctx, models.Address("foo/bar"))
				assert.NoError(t, err)
				assert.Equal(t, true, exists)
				assert.Equal(t, pubKey, pk)
			},
		},
	}
	runTest(t, testCases)
}
