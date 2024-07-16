package api

import (
	"bank_system/token"
	"bank_system/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config     util.Config
	database   *gorm.DB
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(db *gorm.DB, config util.Config) (*Server, error) {
	tokenmaker, err := token.NewJWTMaker(config.Token.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	server := &Server{
		config:     config,
		database:   db,
		tokenMaker: tokenmaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/users/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.POST("/users/register", server.createUser)
	router.GET("/users/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/users/login", server.loginUser)
	router.POST("/users/renewAccessToken", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/account", server.accountPage)
	server.router = router
}

func (server *Server) Start() error {
	if server.config.Server.Port[0] != ':' {
		server.config.Server.Port = ":" + server.config.Server.Port
	}
	return server.router.Run(server.config.Server.Port)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
