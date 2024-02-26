package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/token"
)

//TODO: add tests for users

// userResponse struct to create a response for a new user
type userResponse struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// TODO: Refactor getLoggedInUser (middleware auth)
// getLoggedInUser gets the logged user
func (server *Server) getLoggedInUser(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.(*token.Payload).Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized to access resource",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// TODO: FIX THIS, DO NOT RETURN PASSWORD
	ctx.JSON(http.StatusOK, user)

}
