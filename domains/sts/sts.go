package sts

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subtitle contains the subtitle elements
type Repository struct {
	db *gorm.DB
}

// New returns a new instance of
func New(db *gorm.DB) (Manager, error) {
	if err := db.AutoMigrate(&Subtitle{}, &Content{}); err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

// FindContentByTimestamp finds a subtitle content by the name and timerange
func (rp *Repository) FindContentByTimeRange(ctx context.Context, subtitleID string, timerange ...string) ([]Content, error) {
	var result []Content
	err := rp.db.WithContext(ctx).
		Where("subtitle_id = ? AND time_range IN ?",
			subtitleID, timerange).
		Find(&result).Error
	return result, err
}

// FindSubtitle finds a subtitle with the subtitle's name
func (rp *Repository) FindSubtitle(ctx context.Context, name string) (*Subtitle, error) {
	var result *Subtitle
	err := rp.db.WithContext(ctx).Where("name = ?", name).First(&result).Error
	return result, err
}

// Create creates a new subtitle
func (rp *Repository) Create(ctx context.Context, s *Subtitle) error {
	s.ID = uuid.NewString()
	return rp.db.WithContext(ctx).Save(&s).Error
}

// CreateContent creates a new content
func (rp *Repository) CreateContent(ctx context.Context, c *Content) error {
	c.ID = uuid.NewString()
	return rp.db.WithContext(ctx).Save(&c).Error
}
