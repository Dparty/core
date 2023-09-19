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
	router.POST("/feieyun", func(ctx *gin.Context) {
		// orderId := ctx
	})
	router.GET("/feieyun/feieyun_verify_3E6TRJ5g81bCsdZI.txt", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/plain")
		ctx.String(http.StatusOK, "3E6TRJ5g81bCsdZI")
	})
	api.AccountApiInterfaceMounter(router, accountApi)
	api.RestaurantApiInterfaceMounter(router, restaurantApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
