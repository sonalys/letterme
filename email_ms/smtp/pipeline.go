package smtp

import (
	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain/models"
)

// EnvelopeMiddleware is a middleware that can be chained and create a pipeline.
type EnvelopeMiddleware func(next EnvelopeHandler) EnvelopeHandler

// EnvelopeHandler is a handler that processes a pipeline.
type EnvelopeHandler func(envelope *models.UnencryptedEmail) error

type Pipeline []EnvelopeHandler

// emptyMiddleware is used as tail for the pipeline.
func emptyMiddleware(envelope *models.UnencryptedEmail) error {
	return nil
}

// AddMiddleware adds a new middleware to the envelope pipeline.
func (p *Pipeline) AddMiddlewares(middlewares ...EnvelopeMiddleware) {
	size := len(middlewares)

	if size == 0 {
		return
	}

	*p = make(Pipeline, size)

	// the last middleware points to empty.
	(*p)[size-1] = middlewares[size-1](emptyMiddleware)

	// the penultimate middleware points to the last and etc.
	for i := 1; i < size; i++ {
		previousHandler := (*p)[size-i]
		currentMiddleware := middlewares[size-1-i]
		*p = append(*p, currentMiddleware(previousHandler))
	}
}

// Start should execute the pipeline.
func (p *Pipeline) Start(envelope *models.UnencryptedEmail) error {
	lastIndex := len(*p) - 1
	if err := (*p)[lastIndex](envelope); err != nil {
		return errors.Wrap(err, "failed to process envelope pipeline")
	}
	return nil
}
