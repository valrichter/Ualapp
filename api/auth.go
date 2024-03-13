package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/token"
	"github.com/valrichter/Ualapp/util"
)

// TODO: add tests for auth
// TODO: Refactor auth

// Auth struct to handle authentication
type Auth struct {
	server *Server
}

// Routing for authentication
func (auth Auth) router(server *Server) {

	auth.server = server

	key := server.config.TokenSimmetricKey
	tokenMaker, err := token.NewPasetoMaker(key)
	if err != nil {
		log.Default().Fatal("cannot create token maker", err)
		return
	}

	auth.server.tokenMaker = tokenMaker

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/register", auth.register)
	serverGroup.POST("/login", auth.login)
}

// userRequest struct to create a new user
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login function for authentication
func (auth Auth) login(ctx *gin.Context) {
	var user UserRequest

	// TODO: Fix Handle errors
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": strings.Split(err.Error(), ":"),
		})
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

	accessToken, accessPayload, err := auth.server.tokenMaker.CreateToken(dbUser.Email, time.Minute*15)

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

// createUser creates a new user on database
func (auth Auth) register(ctx *gin.Context) {
	var req UserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID:             uuid.New(),
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := auth.server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			// check if user with that email already exists (23505 is for unique_violation)
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "user with that email already exists",
				})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// return created user
	response := newUserResponse(user)
	ctx.JSON(http.StatusCreated, response)
}
