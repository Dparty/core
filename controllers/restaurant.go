package controllers

import (
	"net/http"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	"github.com/Dparty/core/services"
	core "github.com/Dparty/model/core"
	model "github.com/Dparty/model/restaurant"
	f "github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

type RestaurantApi struct{}

// FinishBill implements coreapi.RestaurantApiInterface.
func (RestaurantApi) FinishBill(gin_context *gin.Context, id string) {
	panic("unimplemented")
}

// GetBill implements coreapi.RestaurantApiInterface.
func (RestaurantApi) GetBill(gin_context *gin.Context, id string) {
	panic("unimplemented")
}

// ListBills implements coreapi.RestaurantApiInterface.
func (RestaurantApi) ListBills(gin_context *gin.Context, restaurantId string, status string, startAt int64, endAt int64, tableId string) {
	panic("unimplemented")
}

// ListTablesBills implements coreapi.RestaurantApiInterface.
func (RestaurantApi) ListTablesBills(gin_context *gin.Context, id string) {
	panic("unimplemented")
}

// GetTable implements coreapi.RestaurantApiInterface.
func (RestaurantApi) GetTable(gin_context *gin.Context, id string) {
	panic("unimplemented")
}

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context, restaurantId string, request api.PutRestaurantRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(ctx *gin.Context, account core.Account, restaurant model.Restaurant) {
			newRestaurant, err := services.UpdateRestaurant(utils.StringToUint(restaurantId), request.Name, request.Description, request.Tags)
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			ctx.JSON(http.StatusOK, RestaurantBackward(newRestaurant))
		})
}

func (RestaurantApi) GetRestaurant(ctx *gin.Context, id string) {
	restaurant := services.GetRestaurant(utils.StringToUint(id))
	if restaurant == nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}
	ctx.JSON(http.StatusOK, RestaurantBackward(*restaurant))
}

func (RestaurantApi) CreateTable(ctx *gin.Context, restaurantId string, request api.PutTableRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
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
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account core.Account) {
			var table model.Table
			services.DB.Find(&table, utils.StringToUint(id))
			middleware.RestaurantOwner(ctx, utils.UintToString(table.RestaurantId),
				func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
					services.DB.Delete(&table)
					c.JSON(http.StatusNoContent, "")
				})
		})
}

func (RestaurantApi) ListRestaurantTable(ctx *gin.Context, restaurantId string) {
	ctx.JSON(http.StatusOK, api.TableList{
		Data: f.Reference(f.Map(services.GetRestaurant(utils.StringToUint(restaurantId)).Tables,
			func(_ int, table model.Table) api.Table {
				return api.Table{
					Id:    utils.UintToString(table.ID),
					Label: table.Label,
				}
			})),
	})
}

func (RestaurantApi) UpdateItem(ctx *gin.Context, id string, request api.PutItemRequest) {
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account core.Account) {
			if request.Pricing < 0 {
				ctx.String(http.StatusBadRequest, "")
			}
			itemId := utils.StringToUint(id)
			item, err := model.FindItem(itemId)
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			restaurant := services.GetRestaurant(item.RestaurantId)
			if account.ID != restaurant.AccountId {
				ctx.JSON(http.StatusNotAcceptable, "")
				return
			}
			item, err = services.UpdateItem(itemId, ItemForward(request))
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			ctx.JSON(http.StatusOK, ItemBackward(item))
		})
}

func (RestaurantApi) DeleteItem(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx,
		func(c *gin.Context, account core.Account) {
			itemId := utils.StringToUint(id)
			item, err := model.FindItem(itemId)
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			restaurant := services.GetRestaurant(item.RestaurantId)
			if account.ID != restaurant.AccountId {
				c.JSON(http.StatusNotAcceptable, gin.H{})
				return
			}
			services.DeleteItem(utils.StringToUint(id))
			c.String(http.StatusNoContent, "")
		})
}

func (RestaurantApi) GetItem(ctx *gin.Context, id string) {
	item, err := model.FindItem(utils.StringToUint(id))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ItemBackward(item))
}

func (RestaurantApi) CreateRestaurant(ctx *gin.Context, request api.PutRestaurantRequest) {
	middleware.IsRoot(ctx,
		func(c *gin.Context, account core.Account) {
			restaurant, _ := services.CreateRestaurant(account.ID, request.Name, request.Description, request.Tags)
			c.JSON(http.StatusCreated, RestaurantBackward(restaurant))
		})
}

func (RestaurantApi) CreateItem(ctx *gin.Context, restaurantId string, request api.PutItemRequest) {
	middleware.RestaurantOwner(ctx, restaurantId,
		func(ctx *gin.Context, account core.Account, restaurant model.Restaurant) {
			item, err := services.CreateItem(restaurant.ID, ItemForward(request))
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			if request.Pricing < 0 {
				ctx.String(http.StatusBadRequest, "")
				return
			}
			ctx.JSON(http.StatusCreated, ItemBackward(item))
		})
}

func (RestaurantApi) ListRestaurantItems(ctx *gin.Context, id string) {
	res := services.GetRestaurant(utils.StringToUint(id))
	ctx.JSON(http.StatusOK, api.ItemList{
		Data: f.Map(services.ListRestaurantItems(res.ID),
			func(_ int, item model.Item) api.Item {
				return ItemBackward(item)
			}),
	})
}

func (RestaurantApi) ListRestaurants(ctx *gin.Context) {
	middleware.GetAccount(ctx,
		func(c *gin.Context, account core.Account) {
			restaurants := services.ListRestaurants(account.ID)
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
	middleware.RestaurantOwner(ctx, restaurantId,
		func(ctx *gin.Context, account core.Account, restaurant model.Restaurant) {
			if err := services.DeleteRestaurant(utils.StringToUint(restaurantId)); err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			ctx.JSON(http.StatusNoContent, "")
		})
}

func (RestaurantApi) UploadItemImage(ctx *gin.Context, id string) {
	middleware.GetAccount(ctx,
		func(ctx *gin.Context, account core.Account) {
			itemId := utils.StringToUint(id)
			item, err := model.FindItem(itemId)
			if err != nil {
				fault.GinHandler(ctx, err)
				return
			}
			middleware.RestaurantOwner(ctx, utils.UintToString(item.RestaurantId),
				func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
					file, _ := c.FormFile("file")
					ctx.JSON(http.StatusCreated, api.Uploading{Url: services.UploadItemImage(item.ID, file)})
				})
		})
}

func (RestaurantApi) CreateBill(ctx *gin.Context, tableId string, request api.CreateBillRequest) {
	table := services.GetTable(utils.StringToUint(tableId))
	if len(request.Orders) == 0 {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	restaurant := services.GetRestaurant(table.RestaurantId)
	var orders model.Orders = make(model.Orders, 0)
	for _, order := range request.Orders {
		pairs := PairsToMap(order.Options)
		order, err := services.CreateOrder(restaurant.ID, utils.StringToUint(order.ItemId), pairs)
		if err != nil {
			fault.GinHandler(ctx, err)
			return
		}
		orders = append(orders, order)
	}
	bill, _ := services.CreateBill(restaurant.Name, table, orders)
	ctx.JSON(http.StatusCreated, api.Bill{
		Id:     utils.UintToString(bill.ID),
		Orders: make([]api.Order, 0),
	})
}

func (RestaurantApi) CreatePrinter(c *gin.Context, id string, request api.PutPrinterRequest) {
	middleware.RestaurantOwner(c, id,
		func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
			printer, err := services.CreatePrinter(model.Printer{
				RestaurantId: restaurant.ID,
				Sn:           request.Sn,
				Name:         request.Name,
				Description:  request.Description,
				Type:         model.PrinterType(request.Type),
			})
			if err != nil {
				fault.GinHandler(c, err)
				return
			}
			c.JSON(http.StatusCreated, api.Printer{
				Id:          utils.UintToString(printer.ID),
				Sn:          printer.Sn,
				Name:        printer.Name,
				Description: printer.Description,
			})
		})
}

func (RestaurantApi) ListPrinters(c *gin.Context, id string) {
	middleware.RestaurantOwner(c, id,
		func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
			c.JSON(http.StatusOK, api.PrinterList{
				Data: f.Map(services.ListPrinters(restaurant.ID),
					func(_ int, printer model.Printer) api.Printer {
						return PrinterBackward(printer)
					}),
			})
		})
}

func (RestaurantApi) UpdatePrinter(c *gin.Context, id string, request api.PutPrinterRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account core.Account) {
			printer := services.GetPrinter(utils.StringToUint(id))
			middleware.RestaurantOwner(c, utils.UintToString(printer.RestaurantId),
				func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
					printer := model.Printer{
						Name:        request.Name,
						Description: request.Description,
						Sn:          request.Sn,
						Type:        model.PrinterType(request.Type),
					}
					p := services.UpdatePrinter(utils.StringToUint(id), printer)
					c.JSON(http.StatusOK, PrinterBackward(p))
				})
		})
}

func (RestaurantApi) DeletePrinter(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account core.Account) {
			printer := services.GetPrinter(utils.StringToUint(id))
			middleware.RestaurantOwner(c, utils.UintToString(printer.RestaurantId),
				func(c *gin.Context, account core.Account, restaurant model.Restaurant) {
					err := services.DeletePrinter(printer.ID)
					if err != nil {
						fault.GinHandler(c, err)
						return
					}
					c.String(http.StatusNoContent, "")
				})
		})
}
