package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	server *Server
}

func (a Auth) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/login", a.login)
}

func (a Auth) login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Route"})
}
