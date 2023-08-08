package services

import (
	"fmt"

	"github.com/Dparty/common/errors"
	"github.com/Dparty/common/utils"
	model "github.com/Dparty/model/restaurant"
)

func CreateRestaurant(accountId uint, name, description string) (model.Restaurant, *errors.Error) {
	restaurant := model.Restaurant{
		Name:        name,
		Description: description,
	}
	restaurant.AccountId = accountId
	DB.Save(&restaurant)
	return restaurant, nil
}

func CreateItem(restaurantId uint, item model.Item) model.Item {
	item.RestaurantId = restaurantId
	DB.Save(&item)
	return item
}

func GetItem(id uint) (model.Item, *errors.Error) {
	var item model.Item
	ctx := DB.Find(&item, id)
	if ctx.RowsAffected == 0 {
		return item, errors.NotFoundError()
	}
	return item, nil
}

func UpdateItem(id uint, item model.Item) model.Item {
	item.ID = id
	var old model.Item
	DB.Find(&old, id)
	old.ID = id
	old.Name = item.Name
	old.Pricing = item.Pricing
	old.Attributes = item.Attributes
	DB.Save(&old)
	return old
}

func DeleteItem(id uint) {
	DB.Delete(&model.Item{}, id)
}

func GetRestaurant(id uint) (model.Restaurant, *errors.Error) {
	var restaurant model.Restaurant
	ctx := DB.Find(&restaurant, id)
	if ctx.RowsAffected == 0 {
		return restaurant, errors.NotFoundError()
	}
	return restaurant, nil
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

func UploadItemImage(id uint) string {
	var item model.Item
	DB.Find(&item, id)
	imageId := utils.GenerteId()
	path := "items/" + utils.UintToString(imageId)
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", Bucket, CosClient.Region, path)
	item.Images = append(item.Images, url)
	DB.Save(&item)
	return CosClient.CreatePresignedURL(Bucket, path)
}

func CreateTable(restaurantId uint, label string) bool {
	table := &model.Table{
		RestaurantId: restaurantId,
		Label:        label,
	}
	ctx := DB.Find(&table)
	if ctx.RowsAffected != 0 {
		return false
	}
	DB.Save(&table)
	return true
}
