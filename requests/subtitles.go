package requests

// Subtitle request DTO for subtitles
type Subtitle struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name"`
	Filename       string    `json:"file_name"`
	SourceLanguage string    `json:"source_language"`
	TargetLanguage string    `json:"target_language"`
	Content        []Content `json:"content"`
}

// Content request DTO for content
type Content struct {
	SubtitleID string `json:"subtitle_id"`
	TimeRange  string `json:"time_range"`
	Content    string `json:"content"`
}
