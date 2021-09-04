package smtp

import (
	"errors"
	"testing"

	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/assert"
)

func Test_Pipeline(t *testing.T) {
	type testCase struct {
		name           string
		env            *models.UnencryptedEmail
		middlewares    func(tC testCase) []EnvelopeMiddleware
		expectedOutput *models.UnencryptedEmail
		expectedError  error
	}
	testCases := []testCase{
		{
			name: "single middleware, no error",
			env: &models.UnencryptedEmail{
				Title: []byte("bananas"),
			},
			expectedOutput: &models.UnencryptedEmail{
				Title: []byte("bananas"),
				From:  models.Address("bananas@letter.me"),
			},
			expectedError: nil,
			middlewares: func(tC testCase) []EnvelopeMiddleware {
				return []EnvelopeMiddleware{
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.From = tC.expectedOutput.From
							return next(env)
						}
					},
				}
			},
		},
		{
			name:           "single middleware, error",
			env:            &models.UnencryptedEmail{},
			expectedOutput: &models.UnencryptedEmail{},
			expectedError:  errors.New("foo/bar"),
			middlewares: func(tC testCase) []EnvelopeMiddleware {
				return []EnvelopeMiddleware{
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							return tC.expectedError
						}
					},
				}
			},
		},
		{
			name: "two middleware, no error",
			env:  &models.UnencryptedEmail{},
			expectedOutput: &models.UnencryptedEmail{
				Title:  []byte("alysson"),
				ToList: []models.Address{"bananas"},
			},
			expectedError: nil,
			middlewares: func(tC testCase) []EnvelopeMiddleware {
				return []EnvelopeMiddleware{
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.Title = tC.expectedOutput.Title
							return next(env)
						}
					},
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.ToList = tC.expectedOutput.ToList
							return next(env)
						}
					},
				}
			},
		},
		{
			name: "two middleware, error on first",
			env:  &models.UnencryptedEmail{},
			expectedOutput: &models.UnencryptedEmail{
				Title: []byte("alysson"),
			},
			expectedError: errors.New("bananas"),
			middlewares: func(tC testCase) []EnvelopeMiddleware {
				return []EnvelopeMiddleware{
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.Title = tC.expectedOutput.Title
							return tC.expectedError
						}
					},
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.ToList = []models.Address{"bananas"}
							return next(env)
						}
					},
				}
			},
		},
		{
			name: "two middleware, error on last",
			env:  &models.UnencryptedEmail{},
			expectedOutput: &models.UnencryptedEmail{
				Title:  []byte("alysson"),
				ToList: []models.Address{"bananas"},
			},
			expectedError: errors.New("bananas"),
			middlewares: func(tC testCase) []EnvelopeMiddleware {
				return []EnvelopeMiddleware{
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.ToList = tC.expectedOutput.ToList
							return next(env)
						}
					},
					func(next EnvelopeHandler) EnvelopeHandler {
						return func(env *models.UnencryptedEmail) error {
							env.Title = tC.expectedOutput.Title
							return tC.expectedError
						}
					},
				}
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			pipeline := Pipeline{}
			middlewares := tC.middlewares(tC)
			pipeline.AddMiddlewares(middlewares...)

			err := pipeline.Start(tC.env)
			assert.True(t, errors.Is(err, tC.expectedError))
			assert.Equal(t, tC.expectedOutput, tC.env)
		})
	}
}
