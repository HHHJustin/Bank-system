package api

import (
	"bank_system/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) accountPage(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	c.HTML(http.StatusOK, "account.html", gin.H{
		"username": authPayload.Username,
	})
}
