package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/valrichter/Ualapp/util"
)

// Auth struct to handle authentication
type Auth struct {
	server *Server
}

// Routing for authentication
func (auth Auth) router(server *Server) {
	auth.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/login", auth.login)
}

// Login function for authentication
func (auth Auth) login(ctx *gin.Context) {
	user := new(userRequest)
	if err := ctx.ShouldBindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dbUser, err := auth.server.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Incorrect mail",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(user.Password, dbUser.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect password",
		})
		return
	}

}
