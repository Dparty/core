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

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context, restaurantId string, request api.PutRestaurantRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
			newRestaurant, err := services.UpdateRestaurant(utils.StringToUint(restaurantId), request.Name, *request.Description)
			if err != nil {
				err.GinHandler(c)
				return
			}
			c.JSON(http.StatusOK, RestaurantBackward(newRestaurant))
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
			if table := services.CreateTable(restaurant.ID, request.Label); table != nil {
				c.JSON(http.StatusCreated, api.Table{
					Id:    utils.UintToString(table.ID),
					Label: table.Label,
				})
			} else {
				c.String(http.StatusConflict, "")
			}
		})
}

func (RestaurantApi) UpdateTable(ctx *gin.Context, id string, gin_body api.PutTableRequest) {

}

func (RestaurantApi) DeleteTable(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx, func(ctx *gin.Context, account api.Account) {
		var table restaurant.Table
		services.DB.Find(&table, utils.StringToUint(id))
		middleware.RestaurantOwner(ctx, utils.UintToString(table.RestaurantId),
			func(c *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
				services.DB.Delete(&table)
			})
	})
}

func (RestaurantApi) ListRestaurantTable(ctx *gin.Context, restaurantId string) {
	tables, err := services.ListRestaurantTable(utils.StringToUint(restaurantId))
	if err != nil {
		err.GinHandler(ctx)
		return
	}
	ctx.JSON(http.StatusOK, api.TableList{
		Data: f.Reference(f.Map(tables, func(_ int, table restaurant.Table) api.Table {
			return api.Table{
				Id:    utils.UintToString(table.ID),
				Label: table.Label,
			}
		})),
	})
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

func (RestaurantApi) GetItem(ctx *gin.Context, id string) {
	item, err := services.GetItem(utils.StringToUint(id))
	if err != nil {
		err.GinHandler(ctx)
		return
	}
	ctx.JSON(http.StatusOK, ItemBackward(item))
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
	res, err := services.GetRestaurant(utils.StringToUint(id))
	if err != nil {
		err.GinHandler(ctx)
	}
	ctx.JSON(http.StatusOK, api.ItemList{
		Data: f.Map(services.ListRestaurantItems(res.ID),
			func(_ int, item restaurant.Item) api.Item {
				return ItemBackward(item)
			}),
	})
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
	middleware.RestaurantOwner(ctx, restaurantId, func(ctx *gin.Context, account api.Account, restaurant restaurant.Restaurant) {
		if err := services.DeleteRestaurant(utils.StringToUint(restaurantId)); err != nil {
			err.GinHandler(ctx)
			return
		}
		ctx.JSON(http.StatusNoContent, "")
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
