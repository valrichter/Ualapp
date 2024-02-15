package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
