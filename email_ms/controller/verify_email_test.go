package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_VerifyEmail(t *testing.T) {
	ctx := context.Background()

	testCases := []testCase{
		{
			name: "fail to publish",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.messaging.On("Publish", messaging.QAccountMS, mock.Anything).
					Return(errors.New("foo/bar"))
				got, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.False(t, got)
			},
		},
		{
			name: "response error",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.messaging.On("Publish", messaging.QAccountMS, mock.Anything).
					Return(nil)

				respChan := make(chan models.Response, 1)

				respChan <- models.Response{
					Error: errors.New("foo/bar"),
				}

				th.router.On("WaitResponse", mock.Anything).
					Return(transformChannel(respChan))
				got, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.False(t, got)
			},
		},
		{
			name: "response decode error",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.messaging.On("Publish", messaging.QAccountMS, mock.Anything).
					Return(nil)

				respChan := make(chan models.Response, 1)
				delivery := models.Delivery{}
				delivery.SetBody(1)
				respChan <- models.Response{
					Message: delivery,
				}

				th.router.On("WaitResponse", mock.Anything).
					Return(transformChannel(respChan))
				got, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.False(t, got)
			},
		},
		{
			name: "all ok",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.messaging.On("Publish", messaging.QAccountMS, mock.Anything).
					Return(nil)

				respChan := make(chan models.Response, 1)
				delivery := models.Delivery{}
				contract := contracts.CheckEmailResponse{
					Exists: true,
				}
				delivery.SetBody(contract)

				respChan <- models.Response{
					Message: delivery,
				}

				th.router.On("WaitResponse", mock.Anything).
					Return(transformChannel(respChan))
				exists, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.NoError(t, err)
				assert.Equal(t, contract.Exists, exists)
			},
		},
	}
	runTest(t, testCases)
}
