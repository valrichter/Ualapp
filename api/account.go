package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/token"
	"github.com/valrichter/Ualapp/util"
)

// Auth struct to handle authentication
type Account struct {
	server *Server
}

type AccountResponse struct {
	ID            int32     `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Balance       int64     `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	AccountNumber string    `json:"account_number"`
}

func (u AccountResponse) ToAccountResponse(account *db.Account) *AccountResponse {
	return &AccountResponse{
		ID:            account.ID,
		UserID:        account.UserID,
		Balance:       account.Balance,
		Currency:      account.Currency,
		CreatedAt:     account.CreatedAt,
		AccountNumber: account.AccountNumber.String,
	}
}

func (u AccountResponse) ToAccountResponses(accounts []db.Account) []AccountResponse {
	accountResponses := make([]AccountResponse, len(accounts))

	for i := range accounts {
		accountResponses[i] = *u.ToAccountResponse(&accounts[i])
	}

	return accountResponses
}

type AccountByNumResponse struct {
	*AccountResponse
	Email string `json:"email"`
}

func (u AccountByNumResponse) ToAccountByNumResponse(account *db.GetAccountByAccountNumberRow) *AccountByNumResponse {
	return &AccountByNumResponse{
		AccountResponse: &AccountResponse{
			ID:            account.ID,
			UserID:        account.UserID,
			Balance:       account.Balance,
			Currency:      account.Currency,
			CreatedAt:     account.CreatedAt,
			AccountNumber: account.AccountNumber.String,
		},
		Email: account.Email,
	}
}

// Routing for authentication
func (account Account) router(server *Server) {
	account.server = server

	serverGroup := server.router.Group("/account").Use(AuthMiddleware(server.tokenMaker))
	serverGroup.GET("/", account.getUserAccounts)
	serverGroup.POST("/create", account.createAccount)
	serverGroup.POST("/transfer", account.createTransfer)
	serverGroup.POST("/add-money", account.addMoney)
	serverGroup.POST("/get-account-by-number", account.getAccountByAccountNumber)
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

	newAccountNumber, err := util.GenerateAccountNumber(newAccount.ID, newAccount.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		err := account.server.store.DeleteAccount(context.Background(), newAccount.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.
			JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newAccount, err = account.server.store.UpdateAccountNumber(context.Background(),
		db.UpdateAccountNumberParams{
			AccountNumber: pgtype.Text{
				String: newAccountNumber,
				Valid:  true,
			},
			ID: newAccount.ID,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, AccountResponse{}.ToAccountResponse(&newAccount))
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

type TransferRequest struct {
	FromAccountID int32 `json:"from_account_id" binding:"required"`
	ToAccountID   int32 `json:"to_account_id" binding:"required"`
	Amount        int64 `json:"amount" binding:"required"`
}

func (account *Account) createTransfer(ctx *gin.Context) {
	var req TransferRequest
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

	argMoney := db.CreateEntryParams{
		AccountID: obj.ToAccountID,
		Amount:    obj.Amount,
	}

	_, err = account.server.store.CreateEntry(context.Background(), argMoney)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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

type GetAccountByAccountNumberRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
}

func (a *Account) getAccountByAccountNumber(ctx *gin.Context) {
	var info GetAccountByAccountNumberRequest

	if err := ctx.ShouldBindJSON(&info); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	acc, err := a.server.store.GetAccountByAccountNumber(ctx, pgtype.Text{
		String: info.AccountNumber,
		Valid:  true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, AccountByNumResponse{}.ToAccountByNumResponse(&acc))
}
