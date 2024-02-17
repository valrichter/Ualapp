package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valrichter/Ualapp/db/sqlc"
)

// Server serves HTTP requests for our banking service
// Contains a store to access the database
// Contains a gin engine to serve HTTP requests

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
// NewServer creates a new HTTP server and setup routing
func NewHTTPServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	// create routes
	server.setupRouter()

	return server, nil
}

// setupRouter sets up the routing for the HTTP server
func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to Ualapp!"})
	})

	// One way to handle routes
	// * Users
	// router.POST("/create_user", server.createUser)
	router.GET("/list_users", server.listUsers)
	router.GET("/users/me", server.getLoggedInUser)
	server.router = router

	// Another way to handle routes
	Auth{}.router(server)

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
