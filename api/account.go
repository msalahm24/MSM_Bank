package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountReq struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *server) listAccounts(ctx *gin.Context) {
	var req listAccountReq
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx,arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,accounts)
}

type deleteAccountReq struct{
	ID int64 `form:"id" binding:"required"`
}
func(server *server)deleteAccount(ctx *gin.Context){
	var req deleteAccountReq
	err:=ctx.ShouldBindQuery(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	err = server.store.DeleteAccount(ctx,req.ID)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"status":fmt.Sprintf("Account with id %v deleted",req.ID)})
}
type putAccountRequest struct {
	Balance    int64 `json:"balance" binding:"required"`
	ID int64 `json:"ID" binding:"required"`
}
func(server *server)putAccount(ctx *gin.Context){
	var req putAccountRequest

	err :=ctx.ShouldBind(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	arg :=db.UpdateAccountParams{
		Balance: req.Balance,
		ID: req.ID,
	}
	account,err := server.store.UpdateAccount(ctx,arg)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return
		}
		ctx.JSON(http.StatusNotFound,errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)
}
