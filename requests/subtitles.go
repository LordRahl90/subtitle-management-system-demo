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
	ContentSeq string `json:"content_seq"`
	SubtitleID string `json:"subtitle_id"`
	TimeRange  string `json:"time_range"`
	Content    string `json:"content"`
}

// Search request DTO for search
type Search struct {
	Name           string   `json:"name"`
	TimeRange      []string `json:"time_range"`
	Source         []string `json:"source,omitempty"`
	SourceLanguage string   `json:"source_language,omitempty"`
	TargetLanguage string   `json:"target_language,omitempty"`
}
