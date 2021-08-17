package controller

import "github.com/sonalys/letterme/domain/models"

// CreateAccount receives a new account model, it should be valid, and it's address should not exist already.
func (s *Service) CreateAccount(account models.Account) (ownershipToken string, err error) {
	return "", nil
}
