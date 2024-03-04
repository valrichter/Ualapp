package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/token"
)

//TODO: add tests for users

type User struct {
	server *Server
}

func (u User) router(server *Server) {
	u.server = server

	serverGroup := server.router.Group("/users", AuthMiddleware(server.token))
	serverGroup.GET("", u.listUsers)
	serverGroup.GET("me", u.getLoggedInUser)
	serverGroup.PATCH("username", u.updateUsername)
}

// listUsers lists all users of database
func (u *User) listUsers(ctx *gin.Context) {
	arg := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := u.server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return list of user
	allUsers := []UserResponse{}
	for _, user := range users {
		allUsers = append(allUsers, newUserResponse(user))
	}

	ctx.JSON(http.StatusOK, allUsers)
}

// getLoggedInUser gets the logged user
func (u *User) getLoggedInUser(ctx *gin.Context) {
	userId, err := u.server.GetActiveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := u.server.store.GetUserById(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))

}

type UpdateUsernameType struct {
	Username string `json:"username" binding:"required"`
}

func (u *User) updateUsername(ctx *gin.Context) {
	userId, err := u.server.GetActiveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var userInfo UpdateUsernameType
	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUsernameParams{
		ID: userId,
		Username: pgtype.Text{
			String: userInfo.Username,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	}

	user, err := u.server.store.UpdateUsername(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

func (s *Server) GetActiveUserID(ctx *gin.Context) (int32, error) {
	// TODO: Refactor middleware authorization
	// authorizationPayload = user_id
	payload := ctx.MustGet("user_id")
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return 0, fmt.Errorf("unauthorized to access resource")
	}

	user, err := s.store.GetUserByEmail(ctx, payload.(*token.Payload).Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized to access resource",
			})
			return 0, fmt.Errorf("unauthorized to access resource")
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 0, err
	}

	return user.ID, nil
}

// userResponse struct to create a response for a new user
type UserResponse struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// TODO: fix newUserResponse and toUserResponse
// newUserResponse creates a new userResponse
func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}