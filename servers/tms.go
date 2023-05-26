package servers

import (
	"translations/domains/tms"
	"translations/requests"
	"translations/responses"

	"github.com/gin-gonic/gin"
)

func (s *Server) tmsRoutes() {
	tmsRoute := s.Router.Group("tms")
	{
		tmsRoute.POST("", s.createTranslation)
		tmsRoute.POST("/upload", s.uploadTranslation)
		tmsRoute.POST("/translate", s.translate)
	}
}

func (s *Server) createTranslation(ctx *gin.Context) {
	var req requests.Translation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}
	e := &tms.Translation{
		Source:         req.Source,
		Target:         req.Target,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
	}
	if err := s.translateService.Create(ctx.Request.Context(), e); err != nil {
		internalError(ctx, err)
		return
	}
	created(ctx, e)
}

func (s *Server) uploadTranslation(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}

	if err := s.translateService.Upload(ctx, file); err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, "translations uploaded successfully")
}

func (s *Server) translate(ctx *gin.Context) {
	var req requests.Translation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	e := &tms.Translation{
		Source:         req.Source,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
	}

	res, err := s.translateService.Translate(ctx, e.Source, e.SourceLanguage, e.TargetLanguage)
	if err != nil {
		internalError(ctx, err)
		return
	}
	success(ctx, responses.TranslateResponse{Target: res})
}
