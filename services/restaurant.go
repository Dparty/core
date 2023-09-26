package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/Dparty/common/errors"
	"github.com/Dparty/common/utils"
	model "github.com/Dparty/model/restaurant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func CreateRestaurant(accountId uint, name, description string, tags []string) (model.Restaurant, *errors.Error) {
	restaurant := model.Restaurant{
		Name:        name,
		Description: description,
	}
	restaurant.AccountId = accountId
	DB.Save(&restaurant)
	return restaurant, nil
}

func UpdateRestaurant(restaurantId uint, name, description string, tags []string) (model.Restaurant, *errors.Error) {
	var restaurant model.Restaurant
	DB.First(&restaurant, restaurantId)
	restaurant.Name = name
	restaurant.Description = description
	ctx := DB.Save(&restaurant)
	if ctx.RowsAffected == 0 {
		return restaurant, errors.NotFoundError()
	}
	return restaurant, nil
}

func DeleteRestaurant(id uint) *errors.Error {
	restaurant := GetRestaurant(id)
	DB.Delete(&restaurant)
	return nil
}

func validItem(item model.Item) bool {
	var attributes map[string]bool = make(map[string]bool)
	for _, attribute := range item.Attributes {
		_, ok := attributes[attribute.Label]
		if ok {
			return false
		}
		var options map[string]bool = make(map[string]bool)
		for _, option := range attribute.Options {
			_, ok := options[option.Label]
			if ok {
				return false
			}
			options[option.Label] = true
		}
		attributes[attribute.Label] = true
	}
	return true
}

func CreateItem(restaurantId uint, item model.Item) (model.Item, *errors.Error) {
	item.RestaurantId = restaurantId
	if !validItem(item) {
		return model.Item{}, &errors.Error{
			StatusCode: 400,
			Code:       4000,
			Message:    "重複屬性",
		}
	}
	DB.Save(&item)
	return item, nil
}

func GetItem(id uint) (model.Item, *errors.Error) {
	var item model.Item
	ctx := DB.Find(&item, id)
	if ctx.RowsAffected == 0 {
		return item, errors.NotFoundError()
	}
	return item, nil
}

// func GetOrderItem(id uint, options model.Options) (model.Order, *errors.Error) {
// 	item, err := GetItem(id)
// 	return model.Order{Name: item.Name, Pricing: item.Pricing}, err
// }

func UpdateItem(id uint, item model.Item) (model.Item, *errors.Error) {
	item.ID = id
	var old model.Item
	DB.Find(&old, id)
	old.ID = id
	old.Name = item.Name
	old.Pricing = item.Pricing
	old.Attributes = item.Attributes
	old.Images = item.Images
	old.Tags = item.Tags
	old.Printers = item.Printers
	if !validItem(item) {
		return old, &errors.Error{
			StatusCode: 400,
			Code:       4000,
			Message:    "重複屬性",
		}
	}
	DB.Save(&old)
	return old, nil
}

func DeleteItem(id uint) {
	DB.Delete(&model.Item{}, id)
}

func GetRestaurant(id uint) *model.Restaurant {
	return model.FindRestaurant(id)
}

func ListRestaurants(accountId uint) []model.Restaurant {
	var restaurants []model.Restaurant = make([]model.Restaurant, 0)
	DB.Where("account_id = ?", accountId).Find(&restaurants)
	return restaurants
}

func ListRestaurantItems(id uint) []model.Item {
	var items []model.Item = make([]model.Item, 0)
	DB.Where("restaurant_id = ?", id).Find(&items)
	return items
}

func GetFileContentType(ouput *os.File) (string, error) {
	buf := make([]byte, 512)
	_, err := ouput.Read(buf)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buf)
	return contentType, nil
}

func UploadItemImage(id uint, file *multipart.FileHeader) string {
	var item model.Item
	DB.Find(&item, id)
	imageId := utils.GenerteId()
	path := "items/" + utils.UintToString(imageId)
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", Bucket, CosClient.Region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  CosClient.SecretID,
			SecretKey: CosClient.SecretKey,
		},
	})
	f, _ := file.Open()
	client.Object.Put(context.Background(), path, f,
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: file.Header.Get("content-type"),
			},
		})
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", Bucket, CosClient.Region, path)
	item.Images = []string{url}
	DB.Save(&item)
	return url
}

func CreateTable(restaurantId uint, label string) *model.Table {
	table := &model.Table{
		RestaurantId: restaurantId,
		Label:        label,
	}
	ctx := DB.Where("restaurant_id = ? AND label = ?", restaurantId, label).Find(&table)
	if ctx.RowsAffected != 0 {
		return nil
	}
	DB.Save(&table)
	return table
}

func GetTable(id uint) model.Table {
	var table model.Table
	DB.Find(&table, id)
	return table
}

// func ListRestaurantTable(restaurantId uint) ([]model.Table, *errors.Error) {
// 	var tables []model.Table = make([]model.Table, 0)
// 	if _, err := GetRestaurant(restaurantId); err != nil {
// 		return tables, err
// 	}
// 	DB.Where("restaurant_id = ?", restaurantId).Find(&tables)
// 	return tables, nil
// }

func CreateOrder(restaurantId, itemId uint, optionsMap map[string]string) (model.Order, *errors.Error) {
	item, err := GetItem(itemId)
	var options model.Options = make(model.Options, 0)
	if item.RestaurantId != restaurantId {
		return model.Order{}, errors.NotFoundError()
	}
	for k, v := range optionsMap {
		option, err := item.Attributes.GetOption(k, v)
		if err != nil {
			return model.Order{}, errors.NotFoundError()
		}
		options = append(options, option)
	}
	return model.Order{
		Item:    item,
		Options: options,
	}, err
}

func CreateBill(restaurantName string, table model.Table, orders model.Orders) (model.Bill, *errors.Error) {
	var lastBill model.Bill
	ctx := DB.Where("restaurant_id = ?", table.RestaurantId).Order("created_at DESC").Find(&lastBill)
	var pickUpCode int64 = 0
	if ctx.RowsAffected != 0 {
		pickUpCode = (lastBill.PickUpCode + 1) % 1000
	}
	bill := model.Bill{RestaurantId: table.RestaurantId,
		Orders: orders, TableLabel: table.Label,
		PickUpCode: pickUpCode}
	DB.Save(&bill)
	PrintBill(restaurantName, bill, table, false)
	return bill, nil
}

func PrintBill(restaurantName string, bill model.Bill, table model.Table, reprint bool) {
	var printers []model.Printer
	DB.Where("restaurant_id = ?", table.RestaurantId).Find(&printers)
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurantName)
	content += fmt.Sprintf("<C><L><B>餐號: A%d</B></L></C>", bill.PickUpCode)
	content += fmt.Sprintf("<C><L><B>桌號: %s</B></L></C>", table.Label)
	content += "--------------------------------<BR>"
	var printersString map[uint]string = make(map[uint]string)
	for _, order := range bill.Orders {
		content += fmt.Sprintf("%s %.2f<BR>", order.Item.Name, float64(order.Item.Pricing)/100)
		attributes := ""
		for _, option := range order.Options {
			attributes += fmt.Sprintf("  |--   %s +%.2f<BR>", option.Label, float64(option.Extra)/100)
		}
		content += attributes
		for _, printer := range order.Item.Printers {
			_, ok := printersString[printer]
			if !ok {
				printersString[printer] = fmt.Sprintf("<C><L><B>桌號: %s</B></L></C><BR>", table.Label)
				printersString[printer] += fmt.Sprintf("<C><L><B>A%d</B></L></C><BR>", bill.PickUpCode)
			}
			printersString[printer] += order.Item.Name + "<BR>"
			printersString[printer] += attributes
		}
	}
	for k, v := range printersString {
		foodPrinter := GetPrinter(k)
		BillPrinter.Print(foodPrinter.Sn, v, "https://ordering-api-uat.sum-foods.com/feieyun/callback")
	}
	content += "--------------------------------<BR>"
	content += fmt.Sprintf("合計:%.2f元<BR>", float64(bill.Total())/100)
	for _, printer := range printers {
		if printer.Type == "BILL" {
			BillPrinter.Print(printer.Sn, content, "https://ordering-api-uat.sum-foods.com/feieyun/callback")
		}
	}
}

func CreatePrinter(printer model.Printer) (model.Printer, *errors.Error) {
	if ctx := DB.Where("sn = ?", printer.Sn).Find(&model.Printer{}); ctx.RowsAffected != 0 {
		return model.Printer{}, errors.PrinterSnDuplecateError()
	}
	DB.Save(&printer)
	return printer, nil
}

func UpdatePrinter(id uint, printer model.Printer) model.Printer {
	printer.ID = id
	DB.Save(&printer)
	return printer
}

func ListPrinters(restaurantId uint) []model.Printer {
	var printers []model.Printer
	DB.Where("restaurant_id = ?", restaurantId).Find(&printers)
	return printers
}

func GetPrinter(id uint) model.Printer {
	var printer model.Printer
	DB.Where("id = ?", id).Find(&printer)
	return printer
}

func DeletePrinter(id uint) *errors.Error {
	printer := GetPrinter(id)
	ctx := DB.Delete(&printer)
	if ctx.RowsAffected == 0 {
		return errors.NotFoundError()
	}
	return nil
}
