package tms

import (
	"context"
	"mime/multipart"

	"translations/domains/tms"
)

// Service is the service definition for TMS
type Service interface {
	Create(ctx context.Context, e *tms.Translation) error
	Upload(ctx context.Context, file *multipart.FileHeader) error
	Translate(ctx context.Context, source, sourceLang, targetLang string) (string, error)
}
