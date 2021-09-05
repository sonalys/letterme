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
				th.router.On("Communicate", messaging.QAccountMS, mock.Anything, mock.Anything).
					Return(errors.New("foo/bar"))
				got, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.Error(t, err)
				assert.False(t, got)
			},
		},
		{
			name: "all ok",
			run: func(t *testing.T, s *Service, th *testHandler) {
				th.router.On("Communicate", messaging.QAccountMS, mock.Anything, mock.Anything).
					Run(mockSetDST(2, contracts.CheckEmailResponse{Exists: true})).
					Return(nil)

				exists, err := s.VerifyEmail(ctx, models.Address("foo/bar"))
				assert.NoError(t, err)
				assert.Equal(t, true, exists)
			},
		},
	}
	runTest(t, testCases)
}
