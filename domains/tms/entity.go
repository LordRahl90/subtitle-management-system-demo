package tms

import "gorm.io/gorm"

// Translation keeps track of the translations within the system
type Translation struct {
	ID             string `json:"id" gorm:"size:100;primaryKey"`
	Source         string `json:"source" gorm:"index:source_source_lang_target_lang"`
	Target         string `json:"target"`
	SourceLanguage string `json:"source_language" gorm:"index:source_source_lang_target_lang"`
	TargetLanguage string `json:"target_language" gorm:"index:source_source_lang_target_lang"`
	gorm.Model
}
