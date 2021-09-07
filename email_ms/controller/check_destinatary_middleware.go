package controller

import (
	"github.com/sonalys/letterme/email_ms/smtp"
)

// CheckDestinataryMiddleware is used to filter only existant recipients from an envelope.
func (s *Service) CheckDestinataryMiddleware(next smtp.EnvelopeHandler) smtp.EnvelopeHandler {
	return func(pipeline *smtp.EmailPipeline) error {
		for _, address := range pipeline.Envelope.ToList {
			exists, publicKey, err := s.getRecipient(s.context, address)
			if err != nil {
				return err
			}
			if exists {
				pipeline.ProcessingEmailList = append(pipeline.ProcessingEmailList, smtp.ProcessingEmail{
					To:        address,
					PublicKey: publicKey,
				})
			}
		}
		return next(pipeline)
	}
}
