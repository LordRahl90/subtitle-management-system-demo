package servers

import (
	"fmt"
	"strings"

	"translations/domains/core"

	"github.com/gin-gonic/gin"
)

// var nonAuthPath = map[string]struct{}{
// 	"/products": {},
// }

func (s *Server) authenticated() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		// if _, ok := nonAuthPath[ctx.FullPath()]; ok {
		// 	ctx.Next()
		// 	return
		// }
		headers := ctx.Request.Header
		authHeader, ok := headers["Authorization"]
		if !ok {
			unAuthorized(ctx, fmt.Errorf("authorization token not provided"))
			ctx.Abort()
			return
		}
		authData := strings.Split(authHeader[0], " ")
		if len(authData) != 2 {
			unAuthorized(ctx, fmt.Errorf("invalid token format provided"))
			ctx.Abort()
			return
		}
		authToken := authData[1]
		userInfo, err := core.Decode(authToken)
		if err != nil || userInfo == nil {
			unAuthorized(ctx, err)
			ctx.Abort()
			return
		}
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

// CORSMiddleware middleware for CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") //in production, this will be mapped to the domain
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Access-Control-Allow-Headers, Access-Control-Allow-Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
