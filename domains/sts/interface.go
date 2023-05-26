package sts

import (
	"context"
)

// Manager manages interface for Subtitle Translations
type Manager interface {
	Create(ctx context.Context, s *Subtitle) error
	CreateContent(ctx context.Context, c *Content) error
	FindSubtitle(ctx context.Context, name, sourceLanguage string) (*Subtitle, error)
	FindContentByTimeRange(ctx context.Context, subtitleID string, timerange ...string) ([]Content, error)
	FindContentBySentences(ctx context.Context, subtitleID string, words ...string) ([]Content, error)
}
