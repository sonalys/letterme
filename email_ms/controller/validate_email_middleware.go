package controller

import (
	"github.com/sonalys/letterme/email_ms/smtp"
)

// ValidateEmailMiddleware is used to filter only existant recipients from an envelope.
func (s *Service) ValidateEmailMiddleware(next smtp.EnvelopeHandler) smtp.EnvelopeHandler {
	return func(pipeline *smtp.EmailPipeline) error {
		for _, address := range pipeline.Envelope.ToList {
			info, err := s.getAccountInfo(s.context, address)
			if err != nil {
				return err
			}
			pipeline.ProcessingEmailList = append(pipeline.ProcessingEmailList, smtp.ProcessingEmail{
				To:          address,
				AccountInfo: info.AccountAddressInfo,
			})
		}
		return next(pipeline)
	}
}
