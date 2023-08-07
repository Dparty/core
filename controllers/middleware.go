package controllers

import (
	"net/http"

	api "github.com/Dparty/core-api"
	"github.com/Dparty/core/services"
	"github.com/Dparty/model/restaurant"

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

func (m Middleware) IsRoot(c *gin.Context, next func(c *gin.Context, account api.Account)) {
	m.GetAccount(c, func(c *gin.Context, account api.Account) {
		if account.Role != api.ROOT {
			errors.PermissionError().GinHandler(c)
			return
		}
		next(c, account)
	})
}

func (m Middleware) RestaurantOwner(c *gin.Context, restaurantId string, next func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant)) {
	m.GetAccount(c, func(c *gin.Context, account api.Account) {
		id := utils.StringToUint(restaurantId)
		restaurant, err := services.GetRestaurant(id)
		if err != nil {
			err.GinHandler(c)
			return
		}
		if utils.StringToUint(account.Id) != restaurant.AccountId {
			c.JSON(http.StatusNotAcceptable, gin.H{})
			return
		}
		next(c, account, restaurant)
	})
}

var middleware Middleware
