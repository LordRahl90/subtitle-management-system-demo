// Package tms
// This is a service package that is the intermediary between the presenter and the domain
package tms

import (
	"context"
	"errors"
	"translations/domains/tms"

	"gorm.io/gorm"
)

var _ Service = (*TranslationService)(nil)

// TranslationService service that manages the interface between upload and preentation
type TranslationService struct {
	translationRepo tms.Manager
}

// New returns a new instance of translation service
func New(repo tms.Manager) Service {
	return &TranslationService{
		translationRepo: repo,
	}
}

// NewWithDefault returns a new translation service with default setup successfully.
func NewWithDefault(db *gorm.DB) (Service, error) {
	tmsRepo, err := tms.New(db)
	if err != nil {
		return nil, err
	}

	return New(tmsRepo), nil
}

// Create implements Service
func (ts *TranslationService) Create(ctx context.Context, e *tms.Translation) error {
	return ts.translationRepo.Create(ctx, e)
}

// Translate takes in the required data and translates it to the target language.
// It returns the source string if the target isn't found
func (ts *TranslationService) Translate(ctx context.Context, e *tms.Translation) (string, error) {
	res, err := ts.translationRepo.Find(ctx, e.SourceLanguage, e.TargetLanguage, e.Source)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.Source, nil
		}
		return "", err
	}
	return res.Target, nil
}

// Upload handles the uploading process
func (ts *TranslationService) Upload(ctx context.Context) error {
	panic("unimplemented")
}
