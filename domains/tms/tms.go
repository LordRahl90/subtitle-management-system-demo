package tms

import (
	"context"

	"gorm.io/gorm"
)

var _ Manager = (*TranslationRepository)(nil)

// TranslationRepository repository for managing translations in the database
type TranslationRepository struct {
	db *gorm.DB
}

// New returns a new instance of translation service
func New(db *gorm.DB) (Manager, error) {
	if err := db.AutoMigrate(&Translation{}); err != nil {
		return nil, err
	}
	return &TranslationRepository{
		db: db,
	}, nil
}

// Find implements Manager
func (tr *TranslationRepository) Find(ctx context.Context, sourceLang string, targetLang string, sentence string) (*Translation, error) {
	var result Translation
	err := tr.db.Where("source_language = ? AND target_language = ? AND source = ?",
		sourceLang, targetLang, sentence).First(&result).Error

	return &result, err
}

// FindByID implements Manager
func (tr *TranslationRepository) FindByID(ctx context.Context, id string) (*Translation, error) {
	var result Translation
	err := tr.db.Where("id = ? ", id).First(&result).Error

	return &result, err
}

// Update implements Manager
func (tr *TranslationRepository) Update(ctx context.Context, e *Translation) error {
	panic("unimplemented")
}

// Create saves a new translation content to the database
func (tr *TranslationRepository) Create(ctx context.Context, e Translation) error {
	return tr.db.WithContext(ctx).Save(&e).Error
}
