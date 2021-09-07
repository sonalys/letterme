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

func Test_getAccountInfo(t *testing.T) {
	ctx := context.Background()

	testCases := []testCase{
		{
			name: "fail to publish",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.router.On("Communicate", messaging.AccountMS, mock.Anything, mock.Anything).
					Return(errors.New("foo/bar"))
				info, err := s.getAccountInfo(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.Nil(t, info)
			},
		},
		{
			name: "all ok",
			run: func(t *testing.T, s *Service, th *testHandler) {
				privKey, err := cryptography.NewPrivateKey(2048)
				require.NoError(t, err)
				pubKey := privKey.GetPublicKey()

				resp := contracts.GetAccountInfoResponse{AccountAddressInfo: &models.AccountAddressInfo{
					PublicKey: pubKey,
				}}
				th.router.On("Communicate", messaging.AccountMS, mock.Anything, mock.Anything).
					Run(mockSetDST(2, resp)).
					Return(nil)

				info, err := s.getAccountInfo(ctx, models.Address("foo/bar"))
				assert.NoError(t, err)
				assert.EqualValues(t, resp, *info)
			},
		},
	}
	runTest(t, testCases)
}
