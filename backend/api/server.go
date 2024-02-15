package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// Server serves HTTP requests for our banking service
// Contains a store to access the database
// Contains a gin engine to serve HTTP requests

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewGinServer(envPath string) *Server {
	config, err := util.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot load config: %s", err))
	}

	connPoll, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database: %s", err))
	}

	store := db.NewPostgreSQLStore(connPoll)

	g := gin.Default()

	return &Server{store: store, router: g}
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(port int) {
	server.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Ualapp"})
	})
	server.router.Run(fmt.Sprintf(":%d", port))
}
