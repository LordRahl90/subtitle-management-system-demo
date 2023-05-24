package tms

import "context"

// Manager defines the expected functions of the translation service
type Manager interface {
	Create(ctx context.Context, e *Translation) error
	Find(ctx context.Context, sourceLang, targetLang, sentence string) (*Translation, error)
	FindByID(ctx context.Context, id string) (*Translation, error)
	Update(ctx context.Context, e *Translation) error
}
