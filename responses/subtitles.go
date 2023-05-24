package responses

// Subtitle response DTO for subtitles
type Subtitle struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Filename       string    `json:"file_name"`
	SourceLanguage string    `json:"source_language"`
	TargetLanguage string    `json:"target_language"`
	Content        []Content `json:"content"`
}

// Content respoj se DTO for content
type Content struct {
	ID         string `json:"id"`
	SubtitleID string `json:"subtitle_id"`
	TimeRange  string `json:"time_range"`
	Content    string `json:"content"`
}
