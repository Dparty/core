package controllers

import (
	"fmt"
	"net/http"

	"github.com/Dparty/common/server"
	api "github.com/Dparty/core-api"
	"github.com/Dparty/core/services"
	"github.com/Dparty/feieyun"
	"github.com/Dparty/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

var db *gorm.DB

func Init(inject *gorm.DB) {
	db = inject
	model.Init(db)
	services.Init(db)
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})
	router.POST("/feieyun/callback", func(ctx *gin.Context) {
		var feieyunCallback feieyun.FeieyunCallback
		if err := ctx.ShouldBind(&feieyunCallback); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(feieyunCallback.Status)
		ctx.JSON(http.StatusOK, gin.H{"data": "info.A "})
	})
	router.GET("/feieyun/feieyun_verify_3E6TRJ5g81bCsdZI.txt", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/plain")
		ctx.String(http.StatusOK, "3E6TRJ5g81bCsdZI")
	})
	router.GET("/cooperate/:id", cooperate)
	var accountApi AccountApi
	api.AccountApiInterfaceMounter(router, accountApi)
	var restaurantApi RestaurantApi
	api.RestaurantApiInterfaceMounter(router, restaurantApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
