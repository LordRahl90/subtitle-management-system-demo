package tms

import (
	"context"
	"translations/domains/tms"
)

// Service is the service definition for TMS
type Service interface {
	Create(ctx context.Context, e *tms.Translation) error
	Upload(ctx context.Context) error
	Translate(ctx context.Context, e *tms.Translation) (string, error)
}
