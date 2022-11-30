package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
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
	if !server.validAccount(ctx, req.FromAccountId, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}
	if account.Currency != currency {
		err = fmt.Errorf("invalid currency:%v for account:%v that has currency: %v", currency, accountID, account.Currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}
	return true
}
