package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valrichter/Ualapp/db/sqlc"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	ctx.JSON(http.StatusOK, users)
}
