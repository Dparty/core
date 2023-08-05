package controllers

import (
	api "github.com/Dparty/core-api"

	"github.com/Dparty/common/errors"

	"github.com/Dparty/common/utils"

	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func (Middleware) GetAccount(c *gin.Context, next func(c *gin.Context, account api.Account)) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		errors.UnauthorizedError().GinHandler(c)
		return
	}
	account := accountApi.GetAccountById(utils.UintToString(auth.AccountId))
	next(c, *account)
}

func (Middleware) IsRoot(c *gin.Context, next func(c *gin.Context, account api.Account)) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		errors.UnauthorizedError().GinHandler(c)
		return
	}
	account := accountApi.GetAccountById(utils.UintToString(auth.AccountId))
	if account.Role != api.ROOT {
		errors.PermissionError().GinHandler(c)
		return
	}
	next(c, *account)
}

var middleware Middleware
