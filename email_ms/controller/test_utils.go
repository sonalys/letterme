package controller

import (
	"reflect"
	"testing"

	"github.com/sonalys/letterme/domain/mocks"
	"github.com/sonalys/letterme/domain/models"
	"github.com/stretchr/testify/mock"
)

func transformChannel(i chan models.Response) <-chan models.Response {
	return i
}

type testHandler struct {
	messaging   *mocks.Messaging
	router      *mocks.Router
	persistence *mocks.Persistence
}

type testCase struct {
	name string
	run  func(t *testing.T, s *Service, th *testHandler)
}

func runTest(t *testing.T, testList []testCase) {
	for _, tC := range testList {
		t.Run(tC.name, func(t *testing.T) {
			th := testHandler{
				messaging:   &mocks.Messaging{},
				router:      &mocks.Router{},
				persistence: &mocks.Persistence{},
			}
			svc := &Service{
				Dependencies: &Dependencies{
					Messaging:   th.messaging,
					Router:      th.router,
					Persistence: th.persistence,
				},
			}
			tC.run(t, svc, &th)
		})
	}
}

func mockSetDST(index int, dstVal interface{}) func(mock.Arguments) {
	return func(args mock.Arguments) {
		argVal := args.Get(index)
		dstVal := reflect.ValueOf(dstVal)
		reflect.ValueOf(argVal).Elem().Set(dstVal)
	}
}
