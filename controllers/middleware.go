package controllers

import (
	"net/http"

	"github.com/Dparty/common/constants"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	"github.com/Dparty/model/core"
	"github.com/Dparty/model/restaurant"
	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func (Middleware) GetAccount(c *gin.Context, next func(c *gin.Context, account core.Account)) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		fault.GinHandler(c, fault.ErrUnauthorized)
		return
	}
	account, err := core.FindAccount(auth.AccountId)
	if err != nil {
		fault.GinHandler(c, err)
		return
	}
	next(c, account)
}

func (m Middleware) IsRoot(c *gin.Context, next func(c *gin.Context, account core.Account)) {
	m.GetAccount(c, func(c *gin.Context, account core.Account) {
		if account.Role != constants.ROOT {
			fault.GinHandler(c, fault.ErrUnauthorized)
			return
		}
		next(c, account)
	})
}

func (m Middleware) RestaurantOwner(c *gin.Context, restaurantId string,
	next func(c *gin.Context, account core.Account, restaurant restaurant.Restaurant)) {
	m.GetAccount(c, func(c *gin.Context, account core.Account) {
		r := restaurant.FindRestaurant(utils.StringToUint(restaurantId))
		if r == nil {
			return
		}
		if account.ID != r.Owner().ID {
			c.JSON(http.StatusNotAcceptable, gin.H{})
			return
		}
		next(c, account, *r)
	})
}

var middleware Middleware
