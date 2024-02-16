package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// userRequeststruct to create a new user
type userRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// userResponse struct to create a response for a new user
type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// newUserResponse creates a new userResponse
func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// createUser creates a new user on database
func (server *Server) createUser(ctx *gin.Context) {
	var req userRequest

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
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pgx.PgError); ok {

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

// listUsers lists all users of database
func (server *Server) listUsers(ctx *gin.Context) {
	arg := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return list of user
	allUsers := []userResponse{}
	for _, user := range users {
		allUsers = append(allUsers, newUserResponse(user))
	}

	ctx.JSON(http.StatusOK, allUsers)
}
