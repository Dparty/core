package controllers

import (
	"net/http"

	"github.com/Dparty/common/constants"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	core "github.com/Dparty/core/services"
	model "github.com/Dparty/model/core"
	"github.com/gin-gonic/gin"
)

type AccountApi struct{}

func (AccountApi) CreateSession(c *gin.Context, createSessionRequest api.CreateSessionRequest) {
	session, sessionError := core.CreateSession(
		createSessionRequest.Email,
		createSessionRequest.Password,
	)
	if sessionError != nil {
		fault.GinHandler(c, sessionError)
		return
	}
	c.JSON(http.StatusCreated, SessionBackward(*session))
}

func (AccountApi) CreateAccount(c *gin.Context, createAccountRequest api.CreateAccountRequest) {
	middleware.IsRoot(c, func(_ *gin.Context, _ model.Account) {
		var createAccountRequest api.CreateAccountRequest
		err := c.ShouldBindJSON(&createAccountRequest)
		if err != nil {
			// TODO: Error message
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		role := constants.USER
		if createAccountRequest.Role != nil {
			role = constants.Role(*createAccountRequest.Role)
		}
		a, createError := core.CreateAccount(
			createAccountRequest.Email,
			createAccountRequest.Password,
			role,
		)
		if createError != nil {
			fault.GinHandler(c, createError)
			return
		}
		c.JSON(http.StatusCreated, AccountBackward(*a))
	})
}

func (AccountApi) ListAccount(gin_context *gin.Context, ordering api.Ordering, index int64, limit int64) {
	// TODO: List account api
}

func (AccountApi) GetAccount(c *gin.Context) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		fault.GinHandler(c, fault.ErrUnauthorized)
		return
	}
	account := core.GetAccountById(auth.AccountId)
	if account == nil {
		fault.GinHandler(c, fault.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, AccountBackward(*account))
}

func (AccountApi) GetAccountById(id string) *api.Account {
	account := core.GetAccountById(utils.StringToUint(id))
	if account == nil {
		return nil
	}
	a := AccountBackward(*account)
	return &a
}

func (a AccountApi) VerifySession(c *gin.Context, sessionVerificationRequest api.SessionVerificationRequest) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		fault.GinHandler(c, fault.ErrUnauthorized)
		return
	}
	account := a.GetAccountById(utils.UintToString(auth.AccountId))
	if account == nil {
		c.JSON(401, gin.H{})
		return
	}
	c.JSON(http.StatusOK, api.Session{
		Account: *account,
	})
}

func (a AccountApi) UpdatePassword(c *gin.Context, request api.UpdatePasswordRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			if err := core.UpdatePassword(account.ID, request.Password, request.NewPassword); err != nil {
				fault.GinHandler(c, err)
				return
			}
			c.String(http.StatusNoContent, "")
		})
}
