package sts

import (
	"context"
	"mime/multipart"
	"translations/requests"
	"translations/responses"
)

// Service interface for managing
type Service interface {
	Upload(ctx context.Context, sourceLanguage string, file *multipart.FileHeader) (string, error)
	Create(ctx context.Context, e *requests.Subtitle) (responses.Subtitle, error)
}
