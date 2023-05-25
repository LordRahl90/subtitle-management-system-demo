package servers

import (
	"translations/requests"

	"github.com/gin-gonic/gin"
)

func (s *Server) stsRoutes() {
	stsRoute := s.Router.Group("sts")
	{
		stsRoute.POST("", s.createSubtitle)
		stsRoute.POST("upload", s.uploadSubtitles)
	}
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

// takes in an array of subtitle files and generates subtitle files in the target language
// security:
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

	for _, file := range files {
		res, err := s.subtitleService.Upload(ctx.Request.Context(), res.ID, sourceLang, targetLang, file)
		if err != nil {
			internalError(ctx, err)
		}
		println("output file: ", res)
	}

	println("Length: ", len(files))
	println(name, sourceLang, targetLang)

	ctx.Header("Content-Disposition", "attachment;filename=hello.txt")
	ctx.Header("Content-Type", "application/text/plain")
	if _, err := ctx.Writer.Write([]byte("hello hello world")); err != nil {
		internalError(ctx, err)
		return
	}
}
