package controllers

import (
	"fmt"
	"net/http"

	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	"github.com/Dparty/core/services"
	"github.com/Dparty/model/restaurant"
	"github.com/gin-gonic/gin"
)

type RestaurantApi struct{}

var restaurantApi RestaurantApi

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context, id string, gin_body api.PutRestaurantRequest) {

}

func (RestaurantApi) GetRestaurant(ctx *gin.Context) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		fmt.Println(account)
	})
}

func (RestaurantApi) CreateTable(ctx *gin.Context, id string, gin_body api.PutTableRequest) {

}

func (RestaurantApi) UpdateTable(ctx *gin.Context, id string, gin_body api.PutTableRequest) {

}

func (RestaurantApi) UpdateItem(ctx *gin.Context, id string, gin_body api.PutItemRequest) {

}

func (RestaurantApi) CreateRestaurant(ctx *gin.Context, request api.PutRestaurantRequest) {
	middleware.IsRoot(ctx, func(c *gin.Context, account api.Account) {
		restaurant, _ := services.CreateRestaurant(utils.StringToUint(account.Id), request.Name, request.Description)
		c.JSON(http.StatusCreated, RestaurantBackward(restaurant))
	})
}

func (RestaurantApi) CreateItem(ctx *gin.Context, restaurantId string, request api.PutItemRequest) {
	middleware.IsRoot(ctx, func(c *gin.Context, account api.Account) {
		services.CreateItem(utils.StringToUint(account.Id), utils.StringToUint(restaurantId),
			restaurant.Item{
				Name:    request.Name,
				Pricing: request.Pricing,
			})
	})
}

func (RestaurantApi) ListRestaurantItems(ctx *gin.Context) {

}
func (RestaurantApi) ListRestaurants(ctx *gin.Context) {
	middleware.GetAccount(ctx, func(c *gin.Context, account api.Account) {
		services.ListRestaurants(utils.StringToUint(account.Id))
	})
}
