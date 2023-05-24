package tms

import (
	"context"

	"github.com/google/uuid"
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

// Find finds a translation
func (tr *TranslationRepository) Find(ctx context.Context, sourceLang string, targetLang string, sentence string) (*Translation, error) {
	var result Translation
	err := tr.db.Where("source_language = ? AND target_language = ? AND source = ?",
		sourceLang, targetLang, sentence).First(&result).Error

	return &result, err
}

// FindByID finds the translation by ID
func (tr *TranslationRepository) FindByID(ctx context.Context, id string) (*Translation, error) {
	var result Translation
	err := tr.db.Where("id = ? ", id).First(&result).Error

	return &result, err
}

// Update updates the target sentence
func (tr *TranslationRepository) Update(ctx context.Context, e *Translation) error {
	return tr.db.
		Model(&Translation{}).
		Where("id = ?", e.ID).
		Update("target", e.Target).Error
}

// Create saves a new translation content to the database
func (tr *TranslationRepository) Create(ctx context.Context, e *Translation) error {
	e.ID = uuid.NewString()
	return tr.db.WithContext(ctx).Save(&e).Error
}
