package servers

import (
	"net/http"

	"translations/domains/users"
	"translations/services/sts"
	"translations/services/tms"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server contains the resources meant to server the endpoints
type Server struct {
	db               *gorm.DB
	Router           *gin.Engine
	userService      users.IUserService
	translateService tms.Service
	subtitleService  sts.Service
	signingSecret    string
	outputDirectory  string
}

// New creates a new instance of server
func New(db *gorm.DB, signingSecret, outputDirectory string) (*Server, error) {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	userService, err := users.New(db)
	if err != nil {
		return nil, err
	}
	translateService, err := tms.NewWithDefault(db)
	if err != nil {
		return nil, err
	}

	subtitleService, err := sts.NewWithDefault(db, outputDirectory)
	if err != nil {
		return nil, err
	}

	s := &Server{
		db:               db,
		Router:           router,
		userService:      userService,
		translateService: translateService,
		subtitleService:  subtitleService,
		signingSecret:    signingSecret,
		outputDirectory:  outputDirectory,
	}

	s.Router.Use(s.authenticated())
	s.Router.POST("/login", s.authenticate)
	s.Router.POST("/users/create", s.createUser)
	s.tmsRoutes()
	s.stsRoutes()

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

func notFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   err,
	})
}

func internalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
