package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/token"
	"github.com/valrichter/Ualapp/util"
)

// TODO: add tests for server

// Server serves HTTP requests for our banking service
// Contains a store to access the database
// Contains a gin engine to serve HTTP requests

type Server struct {
	store  db.Store
	router *gin.Engine
	config util.Config
	token  token.Maker
}

// NewServer creates a new HTTP server and setup routing
// NewServer creates a new HTTP server and setup routing
func NewHTTPServer(store db.Store) (*Server, error) {

	config, err := util.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	token, err := token.NewPasetoMaker(config.TokenSimmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:  store,
		token:  token,
		config: config,
	}

	// create routes
	server.setupRouter()

	return server, nil
}

// setupRouter sets up the routing for the HTTP server
func (server *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.Use(cors.Default())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to Ualapp!"})
	})

	// One way to handle routes
	// * Users
	server.router = router
	serverGroup := server.router.Group("/users", AuthMiddleware(server.token))
	serverGroup.GET("", server.listUsers)
	serverGroup.GET("me", server.getLoggedInUser)
	serverGroup.PATCH("username")

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
