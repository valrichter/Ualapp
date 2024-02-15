package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/valrichter/Ualapp/db/sqlc"
)

// Server serves HTTP requests for our banking service
// Contains a store to access the database
// Contains a gin engine to serve HTTP requests

type Server struct {
	store db.Store
	gin   *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewGinServer(port int) {
	g := gin.Default()

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Run(fmt.Sprintf(":%d", port))
}
