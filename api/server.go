package api

import (
	"bank_system/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config   util.Config
	database *gorm.DB
	router   *gin.Engine
}

func NewServer(db *gorm.DB, config util.Config) (*Server, error) {
	server := &Server{
		config:   config,
		database: db,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

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
