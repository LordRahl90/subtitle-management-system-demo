package servers

import (
	"errors"
	"fmt"
	"os"

	"translations/requests"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) stsRoutes() {
	stsRoute := s.Router.Group("sts")
	{
		stsRoute.POST("", s.createSubtitle)
		stsRoute.POST("upload", s.uploadSubtitles)
		stsRoute.GET("download/:filename", s.downloadFile)
		stsRoute.POST("search/words", s.searchByWords)
		stsRoute.POST("search/time-range", s.searchByTimeRange)
	}
}

func (s *Server) searchByWords(ctx *gin.Context) {
	var req requests.Search
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	res, err := s.subtitleService.FindContentBySentences(ctx, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx, "records not found")
			return
		}
		internalError(ctx, err)
		return
	}

	success(ctx, res)
}

func (s *Server) searchByTimeRange(ctx *gin.Context) {
	var req requests.Search
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	res, err := s.subtitleService.FindContentByTimeRange(ctx, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx, "records not found")
			return
		}
		internalError(ctx, err)
		return
	}

	success(ctx, res)
}

func (s *Server) createSubtitle(ctx *gin.Context) {
	var req *requests.Subtitle
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}
	res, err := s.subtitleService.Create(ctx, req)
	if err != nil {
		internalError(ctx, err)
		return
	}

	created(ctx, res)
}

// Takes in an array of subtitle files and generates subtitle files in the target language
// Security considerations:
// it's an authenticated endpoint
// it's subject to cors preview
// there is upload limit
// there should also be a check on the type of file that was uploaded
func (s *Server) uploadSubtitles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}
	name := ctx.PostForm("name")
	sourceLang := ctx.PostForm("source_language")
	targetLang := ctx.PostForm("target_language")
	files := form.File["files"]

	subtitle := requests.Subtitle{
		Name:           name,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}
	res, err := s.subtitleService.Create(ctx, &subtitle)
	if err != nil {
		internalError(ctx, err)
		return
	}

	outputFile := []string{}

	for _, file := range files {
		res, err := s.subtitleService.Upload(ctx.Request.Context(), s.outputDirectory, res.ID, sourceLang, targetLang, file)
		if err != nil {
			internalError(ctx, err)
			return
		}
		outputFile = append(outputFile, res)
	}

	// the file path are returned to be downloaded individually.
	success(ctx, outputFile)

}

func (s *Server) downloadFile(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	file := fmt.Sprintf("%s/%s", s.outputDirectory, fileName)

	content, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			notFound(ctx, "file not found")
			return
		}
		internalError(ctx, err)
		return
	}

	ctx.Header("Content-Disposition", "attachment;filename="+fileName)
	ctx.Header("Content-Type", "application/text/plain")
	if _, err := ctx.Writer.Write(content); err != nil {
		internalError(ctx, err)
		return
	}
}
