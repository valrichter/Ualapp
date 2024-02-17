package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/valrichter/Ualapp/token"
	"github.com/valrichter/Ualapp/util"
)

//TODO: add tests for auth

// Auth struct to handle authentication
type Auth struct {
	server     *Server
	tokenMaker token.Maker
}

// Routing for authentication
func (auth Auth) router(server *Server) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Default().Fatal("cannot load config:", err)
		return
	}

	auth.server = server

	key := config.TokenSimmetricKey
	tokenMaker, err := token.NewPasetoMaker(key)
	if err != nil {
		log.Default().Fatal("cannot create token maker", err)
		return
	}

	auth.tokenMaker = tokenMaker

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

	accessToken, accessPayload, err := auth.tokenMaker.CreateToken(
		dbUser.Email,
		time.Minute*15,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user_email":   accessPayload.Username,
		"expires_at":   accessPayload.ExpiredAt,
	})

}
