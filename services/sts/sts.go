package sts

import (
	"bufio"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"translations/domains/sts"
	"translations/requests"
	"translations/responses"
	"translations/services/tms"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// SubtitleService service for managing the connection between the domain and presentation
type SubtitleService struct {
	translationService tms.Service
	subtitleRepository sts.Manager
	outputDirectory    string
}

// New returns a new service with the provided repo
func New(repo sts.Manager, tm tms.Service, outputDirectory string) Service {
	return &SubtitleService{
		subtitleRepository: repo,
		translationService: tm,
		outputDirectory:    outputDirectory,
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
// returns the file name to the calling method
func (ss *SubtitleService) Upload(ctx context.Context, outputDirectory, subtitleID, sourceLanguage, targetLanguage string, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", nil
	}
	defer f.Close()
	fileName := file.Filename + "-" + primitive.NewObjectID().Hex() + ".txt"
	outputPath := outputDirectory + fileName
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
			return "", err
		}
	}

	// to optimize and also to build a cloud native solution,
	// this file should be uploaded to a cloud storage and the path
	// sent to the client.
	// This way, the pods can die gracefully without having to lose any(much) file(s) in transit.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		ct := parseLine(ctx, line, subtitleID)
		if ct == nil {
			continue
		}
		if err := ss.subtitleRepository.CreateContent(ctx, ct); err != nil {
			return "", err
		}

		// lookup the translation as well
		res, err := ss.translationService.Translate(ctx, ct.Content, sourceLanguage, targetLanguage)
		if err != nil {
			return "", err
		}

		if _, err := outputFile.Write(
			[]byte(fmt.Sprintf("%s [%s] %s\n", ct.ContentSeq, ct.TimeRange, res))); err != nil {
			return "", nil
		}
	}

	return fileName, nil
}

func parseLine(ctx context.Context, line, subtitleID string) *sts.Content {
	if line == "" {
		return nil
	}
	bs := strings.Index(line, "[")
	be := strings.Index(line, "]")
	seq := strings.Trim(line[:bs], " ")
	tr := line[bs+1 : be]
	content := strings.Trim(line[be+1:], " ")
	return &sts.Content{
		SubtitleID: subtitleID,
		ContentSeq: seq,
		TimeRange:  tr,
		Content:    content,
	}
}

// NewWithDefault return a subtitle service with default connection
// outputs is the default
func NewWithDefault(db *gorm.DB, outputDirectory string) (Service, error) {
	repo, err := sts.New(db)
	if err != nil {
		return nil, err
	}
	tmsRepo, err := tms.NewWithDefault(db)
	if err != nil {
		return nil, err
	}
	return New(repo, tmsRepo, outputDirectory), nil
}
