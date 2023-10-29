package services

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	model "github.com/Dparty/model/restaurant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func CreateRestaurant(accountId uint, name, description string, tags []string) (model.Restaurant, error) {
	restaurant := model.Restaurant{
		Name:        name,
		Description: description,
		Tags:        tags,
	}
	restaurant.AccountId = accountId
	DB.Save(&restaurant)
	return restaurant, nil
}

func UpdateRestaurant(restaurantId uint, name, description string, tags []string) (model.Restaurant, error) {
	var restaurant model.Restaurant
	DB.First(&restaurant, restaurantId)
	restaurant.Name = name
	restaurant.Description = description
	restaurant.Tags = tags
	ctx := DB.Save(&restaurant)
	if ctx.RowsAffected == 0 {
		return restaurant, fault.ErrNotFound
	}
	return restaurant, nil
}

func DeleteRestaurant(id uint) error {
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

func CreateItem(restaurantId uint, item model.Item) (model.Item, error) {
	item.RestaurantId = restaurantId
	if !validItem(item) {
		return model.Item{}, fault.ErrUndefined
	}
	DB.Save(&item)
	return item, nil
}

func UpdateItem(id uint, item model.Item) (model.Item, error) {
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
		return old, fault.ErrUndefined
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

func CreateOrder(restaurantId, itemId uint, optionsMap map[string]string) (model.Order, error) {
	item, err := model.FindItem(itemId)
	var options []model.Pair = make([]model.Pair, 0)
	if item.RestaurantId != restaurantId {
		return model.Order{}, fault.ErrNotFound
	}
	for k, v := range optionsMap {
		option, err := item.Attributes.GetOption(k, v)
		if err != nil {
			return model.Order{}, fault.ErrNotFound
		}
		options = append(options, option)
	}
	return model.Order{
		Item:          item,
		Specification: options,
	}, err
}

func CreateBill(restaurantName string, table model.Table, orders model.Orders) (model.Bill, error) {
	var lastBill model.Bill
	ctx := DB.Where("restaurant_id = ?", table.RestaurantId).Order("created_at DESC").Find(&lastBill)
	var pickUpCode int64 = 0
	if ctx.RowsAffected != 0 {
		pickUpCode = (lastBill.PickUpCode + 1) % 1000
	}
	bill := model.Bill{RestaurantId: table.RestaurantId,
		Orders: orders, TableLabel: table.Label,
		TableId:    table.ID,
		PickUpCode: pickUpCode}
	DB.Save(&bill)
	PrintBill(restaurantName, bill, table, false)
	return bill, nil
}

type OrderNumber struct {
	Number int
	Order  model.Order
}

func Equal(o model.Order, order model.Order) bool {
	if o.Item.ID != order.Item.ID {
		return false
	}
	om := o.SpecificationToMap()
	tm := order.SpecificationToMap()
	if len(om) != len(tm) {
		return false
	}
	for k, v := range om {
		if tm[k] != v {
			return false
		}
	}
	return true
}

func PrintHelper(order model.Order, orders []OrderNumber) []OrderNumber {
	for i, o := range orders {
		if Equal(o.Order, order) {
			orders[i].Number++
			return orders
		}
	}
	orders = append(orders, OrderNumber{
		Number: 1,
		Order:  order,
	})
	return orders
}

func PrintBill(restaurantName string, bill model.Bill, table model.Table, reprint bool) {
	var printers []model.Printer
	timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	orderNumbers := make([]OrderNumber, 0)
	for _, order := range bill.Orders {
		orderNumbers = PrintHelper(order, orderNumbers)
	}
	DB.Where("restaurant_id = ?", table.RestaurantId).Find(&printers)
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurantName)
	content += fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
	content += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
	content += "--------------------------------<BR>"
	var printersString map[uint]string = make(map[uint]string)
	for _, order := range orderNumbers {
		content += fmt.Sprintf("<B>%s %.2fX%d</B><BR>", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)
		attributes := ""
		attributesWithoutMonth := ""
		for _, option := range order.Order.Specification {
			attributes += fmt.Sprintf("<B>|-- %s +%.2f</B><BR>", option.Right, float64(order.Order.Extra(option))/100)
			attributesWithoutMonth += fmt.Sprintf("<B>|--  %s</B><BR>", option.Right)
		}
		content += attributes
		for _, printer := range order.Order.Item.Printers {
			_, ok := printersString[printer]
			if !ok {
				printersString[printer] = fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
				printersString[printer] += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
			}
			printersString[printer] += fmt.Sprintf("<B>%sX%d</B><BR>", order.Order.Item.Name, order.Number)
			printersString[printer] += attributesWithoutMonth
		}
	}
	for k, v := range printersString {
		foodPrinter := GetPrinter(k)
		p, _ := printerFactory.Connect(foodPrinter.Sn)
		p.Print(v+"<BR>"+timestring, "")
	}
	content += "--------------------------------<BR>"
	content += fmt.Sprintf("<B>合計: %.2f元 (-5%%)</B><BR>", math.Floor(float64(bill.Total())/100*0.95))
	content += timestring
	for _, printer := range printers {
		if printer.Type == "BILL" {
			p, _ := printerFactory.Connect(printer.Sn)
			p.Print(content, "")
		}
	}
}

func CreatePrinter(printer model.Printer) (model.Printer, error) {
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

func DeletePrinter(id uint) error {
	printer := GetPrinter(id)
	if printer.InUsed() {
		return fault.ErrBadRequest
	}
	ctx := DB.Delete(&printer)
	if ctx.RowsAffected == 0 {
		return fault.ErrNotFound
	}
	return nil
}
