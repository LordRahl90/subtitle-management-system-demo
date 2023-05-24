package servers

import (
	"net/http"

	"translations/domains/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server contains the resources meant to server the endpoints
type Server struct {
	db          *gorm.DB
	Router      *gin.Engine
	userService users.IUserService
}

// New creates a new instance of server
func New(db *gorm.DB) (*Server, error) {
	router := gin.Default()
	userService, err := users.New(db)
	if err != nil {
		return nil, err
	}

	s := &Server{
		db:          db,
		Router:      router,
		userService: userService,
	}

	s.Router.POST("/login", s.authenticate)
	s.Router.POST("/users/create", s.createUser)
	s.Router.Use(s.authenticated())

	return s, nil
}

func success(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}

func created(c *gin.Context, response interface{}) {
	c.JSON(http.StatusCreated, response)
}

func unAuthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

func badRequestFromError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "not found",
	})
}

func internalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
