// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	account_managermodels "github.com/sonalys/letterme/account_ms/models"

	cryptography "github.com/sonalys/letterme/domain/cryptography"

	mock "github.com/stretchr/testify/mock"

	models "github.com/sonalys/letterme/domain/models"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddNewDevice provides a mock function with given fields: ctx, accountID
func (_m *Service) AddNewDevice(ctx context.Context, accountID models.DatabaseID) (*cryptography.EncryptedBuffer, error) {
	ret := _m.Called(ctx, accountID)

	var r0 *cryptography.EncryptedBuffer
	if rf, ok := ret.Get(0).(func(context.Context, models.DatabaseID) *cryptography.EncryptedBuffer); ok {
		r0 = rf(ctx, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cryptography.EncryptedBuffer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.DatabaseID) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authenticate provides a mock function with given fields: ctx, Address
func (_m *Service) Authenticate(ctx context.Context, Address models.Address) (*cryptography.EncryptedBuffer, error) {
	ret := _m.Called(ctx, Address)

	var r0 *cryptography.EncryptedBuffer
	if rf, ok := ret.Get(0).(func(context.Context, models.Address) *cryptography.EncryptedBuffer); ok {
		r0 = rf(ctx, Address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cryptography.EncryptedBuffer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Address) error); ok {
		r1 = rf(ctx, Address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAccount provides a mock function with given fields: ctx, account
func (_m *Service) CreateAccount(ctx context.Context, account account_managermodels.CreateAccountRequest) (cryptography.EncryptedBuffer, error) {
	ret := _m.Called(ctx, account)

	var r0 cryptography.EncryptedBuffer
	if rf, ok := ret.Get(0).(func(context.Context, account_managermodels.CreateAccountRequest) cryptography.EncryptedBuffer); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Get(0).(cryptography.EncryptedBuffer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, account_managermodels.CreateAccountRequest) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccount provides a mock function with given fields: ctx, ownershipToken
func (_m *Service) DeleteAccount(ctx context.Context, ownershipToken models.OwnershipKey) error {
	ret := _m.Called(ctx, ownershipToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OwnershipKey) error); ok {
		r0 = rf(ctx, ownershipToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccount provides a mock function with given fields: ctx, ownershipToken
func (_m *Service) GetAccount(ctx context.Context, ownershipToken models.OwnershipKey) (models.Account, error) {
	ret := _m.Called(ctx, ownershipToken)

	var r0 models.Account
	if rf, ok := ret.Get(0).(func(context.Context, models.OwnershipKey) models.Account); ok {
		r0 = rf(ctx, ownershipToken)
	} else {
		r0 = ret.Get(0).(models.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.OwnershipKey) error); ok {
		r1 = rf(ctx, ownershipToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicKey provides a mock function with given fields: ctx, address
func (_m *Service) GetPublicKey(ctx context.Context, address models.Address) (*cryptography.PublicKey, error) {
	ret := _m.Called(ctx, address)

	var r0 *cryptography.PublicKey
	if rf, ok := ret.Get(0).(func(context.Context, models.Address) *cryptography.PublicKey); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cryptography.PublicKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Address) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetPublicKey provides a mock function with given fields: ctx, req
func (_m *Service) ResetPublicKey(ctx context.Context, req account_managermodels.ResetPublicKeyRequest) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, account_managermodels.ResetPublicKeyRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
