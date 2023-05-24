package sts

import (
	"translations/domains/sts"

	"gorm.io/gorm"
)

// SubtitleService service for managing the connection between the domain and presentation
type SubtitleService struct {
	subtitleRepository sts.Manager
}

// New returns a new service with the provided repo
func New(repo sts.Manager) Service {
	return &SubtitleService{
		subtitleRepository: repo,
	}
}

// NewWithDefault return a subtitle service with default connection
func NewWithDefault(db *gorm.DB) (Service, error) {
	repo, err := sts.New(db)
	if err != nil {
		return nil, err
	}
	return New(repo), nil
}
