package requests

// Translation Request DTO for translation
type Translation struct {
	Source         string `json:"source" binding:"required"`
	Target         string `json:"target"`
	SourceLanguage string `json:"sourceLanguage" binding:"required"`
	TargetLanguage string `json:"targetLanguage"`
}
