package controllers

import (
	"net/http"

	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	"github.com/Dparty/core/services"
	restaurant "github.com/Dparty/model/restaurant"
	f "github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

type RestaurantApi struct{}

var restaurantApi RestaurantApi

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context, restaurantId string, gin_body api.PutRestaurantRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
			// services.CreateTa
		})
}

func (RestaurantApi) GetRestaurant(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		restaurant, err := services.GetRestaurant(utils.StringToUint(id))
		if err != nil {
			err.GinHandler(ctx)
			return
		}
		ctx.JSON(http.StatusOK, RestaurantBackward(restaurant))
	})
}

func (RestaurantApi) CreateTable(ctx *gin.Context, restaurantId string, request api.PutTableRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
			ok := services.CreateTable(restaurant.ID, request.Label)
			if ok {
				c.JSON(http.StatusCreated, api.Table{
					Label: request.Label,
				})
			} else {
				c.String(http.StatusConflict, "")
			}
		})
}

func (RestaurantApi) UpdateTable(ctx *gin.Context, id string, gin_body api.PutTableRequest) {

}

func (RestaurantApi) UpdateItem(ctx *gin.Context, id string, request api.PutItemRequest) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		itemId := utils.StringToUint(id)
		item, err := services.GetItem(itemId)
		if err != nil {
			err.GinHandler(c)
			return
		}
		restaurant, _ := services.GetRestaurant(item.RestaurantId)
		if utils.StringToUint(account.Id) != restaurant.AccountId {
			c.JSON(http.StatusNotAcceptable, gin.H{})
			return
		}
		item = services.UpdateItem(itemId, ItemForward(request))
		c.JSON(http.StatusOK, ItemBackward(item))
	})
}

func (RestaurantApi) DeleteItem(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		itemId := utils.StringToUint(id)
		item, err := services.GetItem(itemId)
		if err != nil {
			err.GinHandler(c)
			return
		}
		restaurant, _ := services.GetRestaurant(item.RestaurantId)
		if utils.StringToUint(account.Id) != restaurant.AccountId {
			c.JSON(http.StatusNotAcceptable, gin.H{})
			return
		}
		services.DeleteItem(utils.StringToUint(id))
		c.String(http.StatusNoContent, "")
	})
}

func (RestaurantApi) CreateRestaurant(ctx *gin.Context, request api.PutRestaurantRequest) {
	middleware.IsRoot(ctx, func(c *gin.Context, account api.Account) {
		description := ""
		if request.Description != nil {
			description = *request.Description
		}
		restaurant, _ := services.CreateRestaurant(utils.StringToUint(account.Id), request.Name, description)
		c.JSON(http.StatusCreated, RestaurantBackward(restaurant))
	})
}

func (RestaurantApi) CreateItem(ctx *gin.Context, restaurantId string, request api.PutItemRequest) {
	middleware.RestaurantOwner(ctx, restaurantId, func(ctx *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
		ctx.JSON(http.StatusCreated, ItemBackward(services.CreateItem(restaurant.ID, ItemForward(request))))
	})
}

func (RestaurantApi) ListRestaurantItems(ctx *gin.Context, id string) {
	ctx.JSON(http.StatusOK,
		f.Map(services.ListRestaurantItems(
			utils.StringToUint(id)),
			func(_ int, item restaurant.Item) api.Item {
				return ItemBackward(item)
			}))
}

func (RestaurantApi) ListRestaurants(ctx *gin.Context) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		restaurants := services.ListRestaurants(utils.StringToUint(account.Id))
		var restauratnList api.RestaurantList = api.RestaurantList{
			Data: make([]api.Restaurant, 0),
		}
		for _, r := range restaurants {
			restauratnList.Data = append(restauratnList.Data, RestaurantBackward(r))
		}
		c.JSON(http.StatusOK, restauratnList)
	})
}
func (RestaurantApi) DeleteRestaurant(ctx *gin.Context, restaurantId string) {
	middleware.RestaurantOwner(ctx, restaurantId, func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {

	})
}

func (RestaurantApi) UploadItemImage(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx, func(ctx *gin.Context, account api.Account) {
		itemId := utils.StringToUint(id)
		item, err := services.GetItem(itemId)
		if err != nil {
			err.GinHandler(ctx)
			return
		}
		middleware.RestaurantOwner(ctx, utils.UintToString(item.RestaurantId), func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
			ctx.JSON(http.StatusCreated, api.Uploading{Url: services.UploadItemImage(item.ID)})
		})
	})
}
