package controllers

import (
	"github.com/Dparty/common/server"
	api "github.com/Dparty/core-api"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	api.AccountApiInterfaceMounter(router, accountApi)
	api.RestaurantApiInterfaceMounter(router, restaurantApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
