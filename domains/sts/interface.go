package sts

import (
	"context"
)

type Manager interface {
	Create(ctx context.Context, s *Subtitle) error
	CreateContent(ctx context.Context, c *Content) error
	FindSubtitle(ctx context.Context, name string) (*Subtitle, error)
	FindContentByTimeRange(ctx context.Context, subtitleID string, timerange ...string) ([]Content, error)
}
