package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simple-bank/db2/sqlc"
)

// when new account is created, balance is always 0
// regarding binding refer to https://gin-gonic.com/en/docs/examples/binding-and-validation/#_top
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}


// Note: (server *Server) is like self in Python, or this in Java
// so you can access it by server.createAccount(ctx *gin.Context)
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}