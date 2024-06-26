package api

import (
	"database/sql"
	"fmt"
	db "github/bekeeeee/simplebank/db/sqlc"
	"github/bekeeeee/simplebank/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountId int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=1"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`

}
func (server *Server) createTransfer(ctx *gin.Context){
	 var req transferRequest

	 if err := ctx.ShouldBindJSON(&req); err!=nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	 }
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fromAccount,valid := server.validAccount(ctx,req.FromAccountId,req.Currency)

	 if !valid{
		return
	 }
	 if authPayload.Username != fromAccount.Owner {
		err:= fmt.Errorf("account [%d] does not belong to user [%s]",req.FromAccountId,authPayload.Username)
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return;
	 }
	_,valid = server.validAccount(ctx,req.ToAccountId,req.Currency)

	 if valid {
		return
	 }
	 arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID: req.ToAccountId,
		Amount: req.Amount,
		}
	result,err := server.store.TransferTx(ctx,arg) 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return;
	}
	ctx.JSON(http.StatusOK,result)

}


func (server *Server) validAccount(ctx *gin.Context,accountId int64,currency string )(db.Account,bool){
	account, err:= server.store.GetAccount(ctx,accountId)
	if err != nil {
		if err == sql.ErrNoRows{
		ctx.JSON(http.StatusNotFound,errorResponse(err))
			return account,false
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return account,false
	}
	if account.Currency != currency {
		err:= fmt.Errorf("account [%d] currency mismatch: %s vs %s",account.ID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return	 account,false
	}
	return account,true
}
