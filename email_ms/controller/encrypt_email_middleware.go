package controller

import (
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
	"github.com/sonalys/letterme/email_ms/smtp"
)

// CheckDestinataryMiddleware encrypts the email body to the specified account public key.
func (s *Service) EncryptEmailMiddleware(next smtp.EnvelopeHandler) smtp.EnvelopeHandler {
	return func(pipeline *smtp.EmailPipeline) error {
		for i := range pipeline.ProcessingEmailList {
			pendingEmail := &pipeline.ProcessingEmailList[i]
			encryptedEmail, err := s.encryptEnvelope(pipeline.Envelope, pendingEmail.PublicKey)
			if err != nil {
				return err
			}
			encryptedEmail.To = pendingEmail.To
			pendingEmail.Email = encryptedEmail
		}
		return next(pipeline)
	}
}

// TODO: move this to domain models.UnencryptedEmail.
func (s *Service) encryptEnvelope(env *models.UnencryptedEmail, pk *cryptography.PublicKey) (*models.Email, error) {
	var err error
	encryptedFrom, err := s.encrypt(pk, env.From)
	if err != nil {
		return nil, err
	}

	encryptedToList, err := s.encrypt(pk, env.ToList)
	if err != nil {
		return nil, err
	}

	encryptedTitle, err := s.encrypt(pk, env.Title)
	if err != nil {
		return nil, err
	}

	var encryptedBody *cryptography.EncryptedBuffer
	if len(env.Text) > 0 {
		encryptedBody, err = s.encrypt(pk, env.Text)
		if err != nil {
			return nil, err
		}
	} else {
		encryptedBody, err = s.encrypt(pk, env.HTML)
		if err != nil {
			return nil, err
		}
	}

	return &models.Email{
		From:   encryptedFrom,
		ToList: encryptedToList,
		Title:  encryptedTitle,
		Body:   encryptedBody,
	}, nil
}
