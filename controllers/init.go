package controllers

import (
	"fmt"
	"net/http"

	"github.com/Dparty/common/server"
	api "github.com/Dparty/core-api"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

type FeieyunCallback struct {
	OrderId string `form:"orderId" json:"orderId" binding:"required"`
	Status  int    `form:"status" json:"status" binding:"required"`
	Stime   int    `form:"stime" json:"stime" binding:"required"`
	Sign    string `form:"sign" json:"sign" binding:"required"`
}

func Init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	router.POST("/feieyun/callback", func(ctx *gin.Context) {
		var feieyunCallback FeieyunCallback
		if err := ctx.ShouldBind(&feieyunCallback); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(feieyunCallback)
		ctx.JSON(200, gin.H{"data": "info.A "})
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
