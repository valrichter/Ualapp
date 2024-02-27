package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/valrichter/Ualapp/db/sqlc"
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
	serverGroup.POST("/create", account.createAccount)
	serverGroup.GET("/", account.getUserAccounts)
}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (account *Account) createAccount(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}

	// TODO: Refactor createAccount
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

	arg := db.CreateAccountParams{
		UserID:   user.ID,
		Currency: acc.Currency,
	}

	newAccount, err := account.server.store.CreateAccount(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			// check if user with that email already exists (23505 is for unique_violation)
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "you already have an account with that currency",
				})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusCreated, newAccount)
}

func (account *Account) getUserAccounts(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey)
	if payload == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}
}
