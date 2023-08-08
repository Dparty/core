package controllers

import (
	"net/http"

	"github.com/Dparty/common/server"
	api "github.com/Dparty/core-api"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	api.AccountApiInterfaceMounter(router, accountApi)
	api.RestaurantApiInterfaceMounter(router, restaurantApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
