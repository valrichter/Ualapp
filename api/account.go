package api

import (
	"context"
	"errors"
	"fmt"
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

	serverGroup := server.router.Group("/account", AuthMiddleware(server.tokenMaker))
	serverGroup.POST("/create", account.createAccount)
	serverGroup.GET("/", account.getUserAccounts)
	serverGroup.POST("/transfer", account.createTransfer)
	serverGroup.POST("/add-balance", account.addMoney)
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
		Balance:  0,
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
	userId, err := account.server.GetActiveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}

	accounts, err := account.server.store.GetAccountByUserId(context.Background(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type transferRequest struct {
	FromAccountID int32 `json:"from_account_id" binding:"required"`
	ToAccountID   int32 `json:"to_account_id" binding:"required"`
	Amount        int64 `json:"amount" binding:"required"`
}

func (account *Account) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := account.validAccount(ctx, req.FromAccountID)
	if !valid {
		return
	}

	toAccount, valid := account.validAccount(ctx, req.ToAccountID)
	if !valid {
		return
	}

	if fromAccount.Balance < req.Amount {
		err := fmt.Errorf("balance is not enough, %d < %d", fromAccount.Balance, req.Amount)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if fromAccount.Currency != toAccount.Currency {
		err := fmt.Errorf("from account currency mismatch, %s != %s", fromAccount.Currency, toAccount.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxRequest{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := account.server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (account *Account) validAccount(ctx *gin.Context, accountID int32) (db.Account, bool) {
	acc, err := account.server.store.GetAccountById(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// account doesn't exist
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return acc, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return acc, false
	}

	return acc, true
}

type AddMoneyRequest struct {
	ToAccountID int32  `json:"to_account_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	Reference   string `json:"reference" binding:"required"`
}

func (account *Account) addMoney(ctx *gin.Context) {
	userId, err := account.server.GetActiveUserID(ctx)
	if err != nil {
		return
	}

	obj := AddMoneyRequest{}
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	acc, err := account.server.store.GetAccountById(ctx, obj.ToAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// account doesn't exist
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if acc.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access resource",
		})
		return
	}

	argMoney := db.CreateMoneyRecordParams{
		UserID:    userId,
		Status:    "pending",
		Amount:    obj.Amount,
		Reference: obj.Reference,
	}

	_, err = account.server.store.CreateMoneyRecord(context.Background(), argMoney)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			// 23505 is for unique_violation
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Record with that reference already exists",
				})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// check money rerod to confirm trasaction status

	argBalance := db.UpdateAccountBalanceParams{
		ID:     obj.ToAccountID,
		Amount: obj.Amount,
	}

	_, err = account.server.store.UpdateAccountBalance(context.Background(), argBalance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Money added successfully"})
}
