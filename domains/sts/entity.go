package sts

import (
	"gorm.io/gorm"
)

// Subtitle setup subtitle files
type Subtitle struct {
	ID             string `json:"id" gorm:"size:100;primaryKey"`
	Name           string `json:"name" gorm:"size:100;uniqueIndex"`
	Filename       string `json:"file_name,omitempty"`
	SourceLanguage string `json:"source_language"`
	gorm.Model
}

// Content handles the content within each file
type Content struct {
	ID         string `json:"id" gorm:"size:100;primaryKey"`
	ContentSeq string `json:"content_seq" gorm:"index:content_seq"`
	SubtitleID string `json:"subtitle_id"`
	TimeRange  string `json:"time_range" gorm:"index:content_time_range"`
	Content    string `json:"content" gorm:"index:content_time_range"`
	gorm.Model
}
