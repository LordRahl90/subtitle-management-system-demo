package tms

import "gorm.io/gorm"

// Translation keeps track of the translations within the system
type Translation struct {
	ID             string `json:"id" gorm:"size:100;primaryKey"`
	Source         string `json:"source"`
	Target         string `json:"target"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	gorm.Model
}
