package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/sonalys/letterme/account_ms/mocks"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
	"github.com/sonalys/letterme/domain/persistence/mongo"
	"github.com/stretchr/testify/require"
)

func Test_verifyEmailExistence(t *testing.T) {
	ctx := context.Background()
	testCases := []testCase{
		{
			name:     "invalid data",
			data:     "123",
			expError: errDecode,
		},
		{
			name: "not found",
			data: contracts.CheckEmailRequest{
				Address: models.Address("bananas@letter.me"),
			},
			preRun: func(t *testing.T, mock *mocks.Service) {
				mock.On("GetPublicKey", ctx, models.Address("bananas@letter.me")).
					Return(nil, mongo.ErrNotFound)
			},
			expResp: contracts.CheckEmailResponse{Exists: false},
		},
		{
			name: "custom error",
			data: contracts.CheckEmailRequest{
				Address: models.Address("bananas@letter.me"),
			},
			preRun: func(t *testing.T, mock *mocks.Service) {
				mock.On("GetPublicKey", ctx, models.Address("bananas@letter.me")).
					Return(nil, errors.New("foo/bar"))
			},
			expError: errInternal,
		},
		{
			name: "found",
			data: contracts.CheckEmailRequest{
				Address: models.Address("bananas@letter.me"),
			},
			preRun: func(t *testing.T, mock *mocks.Service) {
				mock.On("GetPublicKey", ctx, models.Address("bananas@letter.me")).
					Return(nil, nil)
			},
			expResp: contracts.CheckEmailResponse{Exists: true},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			msg := messaging.Delivery{}
			err := msg.SetBody(tC.data)
			require.NoError(t, err)

			mock := &mocks.Service{}
			handler := &Handler{mock}

			if tC.preRun != nil {
				tC.preRun(t, mock)
			}
			got, err := handler.verifyEmailExistence(ctx, msg)
			require.Equal(t, tC.expResp, got)
			if tC.expError != nil {
				require.True(t, errors.Is(err, tC.expError))
			}

			mock.AssertExpectations(t)
		})
	}
}
