package api

import (
	"bank_system/token"
	"bank_system/util"
	"fmt"

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
	fmt.Println(server.config.Token.SecretKey)
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
