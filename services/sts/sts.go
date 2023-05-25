package sts

import (
	"context"
	"mime/multipart"

	"translations/domains/sts"
	"translations/requests"
	"translations/responses"
	"translations/services/tms"

	"gorm.io/gorm"
)

// SubtitleService service for managing the connection between the domain and presentation
type SubtitleService struct {
	translationService tms.Service
	subtitleRepository sts.Manager
}

// New returns a new service with the provided repo
func New(repo sts.Manager, tm tms.Service) Service {
	return &SubtitleService{
		subtitleRepository: repo,
		translationService: tm,
	}
}

// Create creates a new subtitle record
func (ss *SubtitleService) Create(ctx context.Context, e *requests.Subtitle) (responses.Subtitle, error) {
	var result responses.Subtitle
	sub := &sts.Subtitle{
		Name:           e.Name,
		Filename:       e.Filename,
		SourceLanguage: e.SourceLanguage,
	}

	if err := ss.subtitleRepository.Create(ctx, sub); err != nil {
		return result, err
	}
	result.ID = sub.ID
	result.Name = sub.Name
	result.SourceLanguage = sub.SourceLanguage
	result.TargetLanguage = e.TargetLanguage
	result.Content = make([]responses.Content, 0, len(e.Content))
	result.Filename = e.Filename

	for _, v := range e.Content {
		c := &sts.Content{
			SubtitleID: sub.ID,
			TimeRange:  v.TimeRange,
			Content:    v.Content,
		}
		if err := ss.subtitleRepository.CreateContent(ctx, c); err != nil {
			return result, err
		}

		// look for the content
		res, err := ss.translationService.Translate(ctx, v.Content, e.SourceLanguage, e.TargetLanguage)
		if err != nil {
			return result, err
		}
		result.Content = append(result.Content, responses.Content{
			ID:         c.ID,
			SubtitleID: sub.ID,
			TimeRange:  v.TimeRange,
			Content:    res,
		})
	}
	return result, nil
}

// Upload uploads a subtitle record
// reads the lines
// parses the lines and extract contents
// generates a translation and writes it to a file
func (ss *SubtitleService) Upload(ctx context.Context, sourceLanguage string, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", nil
	}
	defer f.Close()

	return "", nil
}

// NewWithDefault return a subtitle service with default connection
func NewWithDefault(db *gorm.DB) (Service, error) {
	repo, err := sts.New(db)
	if err != nil {
		return nil, err
	}
	tmsRepo, err := tms.NewWithDefault(db)
	if err != nil {
		return nil, err
	}
	return New(repo, tmsRepo), nil
}
