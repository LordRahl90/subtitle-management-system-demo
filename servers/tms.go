package servers

import (
	"translations/domains/tms"
	"translations/requests"

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

	res, err := s.translateService.Translate(ctx, e)
	if err != nil {
		internalError(ctx, err)
		return
	}
	success(ctx, res)
}
