package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/token"
)

type transferRequest struct {
	FromAccountId int64  `json:"fromAccountID" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountID" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.TransferTxParms{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	fromAccount, valid := server.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.UserName {
		err = errors.New("from account does not belong to the authentcated user")
		ctx.JSON(http.StatusUnauthorized,errResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("invalid currency:%v for account:%v that has currency: %v", currency, accountID, account.Currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account, true
}
