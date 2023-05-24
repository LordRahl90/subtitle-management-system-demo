package servers

import (
	"translations/requests"

	"github.com/gin-gonic/gin"
)

func (s *Server) stsRoutes() {
	stsRoute := s.Router.Group("sts")
	{
		stsRoute.POST("", s.createSubtitle)
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
