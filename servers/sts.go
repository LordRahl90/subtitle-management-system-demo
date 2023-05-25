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

	for _, file := range files {
		println(file.Filename)
	}

	println("Length: ", len(files))
	println(name, sourceLang, targetLang)
}
