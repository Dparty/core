package controllers

import (
	api "github.com/Dparty/core-api"
	"github.com/gin-gonic/gin"
)

type RestaurantApi struct{}

var restaurantApi RestaurantApi

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context, gin_body api.PutRestaurantRequest) {

}
func (RestaurantApi) GetRestaurant(ctx *gin.Context) {

}
func (RestaurantApi) CreateTable(ctx *gin.Context, gin_body api.PutTableRequest) {

}
func (RestaurantApi) UpdateTable(ctx *gin.Context, gin_body api.PutTableRequest) {

}
func (RestaurantApi) UpdateItem(ctx *gin.Context, gin_body api.PutItemRequest) {

}
func (RestaurantApi) CreateRestaurant(ctx *gin.Context, gin_body api.PutRestaurantRequest) {

}
func (RestaurantApi) CreateItem(ctx *gin.Context, gin_body api.PutItemRequest) {

}
func (RestaurantApi) ListRestaurantItems(ctx *gin.Context) {

}
