package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/valrichter/Ualapp/token"
)

// Auth struct to handle authentication
type Account struct {
	server *Server
}

// Routing for authentication
func (account Account) router(server *Server) {

	account.server = server

	serverGroup := server.router.Group("/account", AuthMiddleware(server.token))
	serverGroup.POST("/create-account", account.createAccount)
}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required"`
}

func (account Account) createAccount(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}

	user, err := account.server.store.GetUserByEmail(ctx, payload.(*token.Payload).Username)
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

	acc := new(AccountRequest)
	if err := ctx.ShouldBindJSON(&acc); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}
