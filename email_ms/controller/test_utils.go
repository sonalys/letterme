package controller

import (
	"reflect"
	"testing"

	"github.com/sonalys/letterme/domain/messaging"
	mMocks "github.com/sonalys/letterme/domain/messaging/mocks"
	"github.com/sonalys/letterme/domain/persistence/mocks"
	"github.com/stretchr/testify/mock"
)

func transformChannel(i chan messaging.Response) <-chan messaging.Response {
	return i
}

type testHandler struct {
	messenger   *mMocks.Messenger
	router      *mMocks.EventRouter
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
				messenger:   &mMocks.Messenger{},
				router:      &mMocks.EventRouter{},
				persistence: &mocks.Persistence{},
			}
			svc := &Service{
				Dependencies: &Dependencies{
					Messenger:   th.messenger,
					EventRouter: th.router,
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
