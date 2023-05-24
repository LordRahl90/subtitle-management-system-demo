package sts

import (
	"gorm.io/gorm"
)

// Subtitle setup subtitle files
type Subtitle struct {
	ID             string `json:"id" gorm:"size:100;primaryKey"`
	Name           string `json:"name" gorm:"size:100;uniqueIndex"`
	Filename       string `json:"file_name"`
	SourceLanguage string `json:"source_language"`
	gorm.Model
}

// Content handles the content within each file
type Content struct {
	ID         string `json:"id" gorm:"size:100;primaryKey"`
	SubtitleID string `json:"subtitle_id"`
	TimeRange  string `json:"time_range" gorm:"index:content_time_range"`
	TimeStart  string `json:"start"`
	TimeStop   string `json:"stop"`
	Content    string `json:"content" gorm:"index:content_time_range"`
	gorm.Model
}
