package services

import (
	"fmt"

	"github.com/Dparty/common/errors"
	"github.com/Dparty/common/utils"
	"github.com/Dparty/model/restaurant"
)

func CreateRestaurant(accountId uint, name, description string) (restaurant.Restaurant, *errors.Error) {
	restaurant := restaurant.Restaurant{
		Name:        name,
		Description: description,
	}
	restaurant.ID = utils.GenerteId()
	restaurant.AccountId = accountId
	DB.Save(&restaurant)
	return restaurant, nil
}

func CreateItem(restaurantId uint, item restaurant.Item) restaurant.Item {
	item.RestaurantId = restaurantId
	item.ID = utils.GenerteId()
	DB.Save(&item)
	return item
}

func GetItem(id uint) (restaurant.Item, *errors.Error) {
	var item restaurant.Item
	ctx := DB.Find(&item, id)
	if ctx.RowsAffected == 0 {
		return item, errors.NotFoundError()
	}
	return item, nil
}

func UpdateItem(id uint, item restaurant.Item) restaurant.Item {
	item.ID = id
	var old restaurant.Item
	DB.Find(&old, id)
	old.ID = id
	old.Name = item.Name
	old.Pricing = item.Pricing
	old.Attributes = item.Attributes
	DB.Save(&old)
	return old
}

func GetRestaurant(id uint) (restaurant.Restaurant, *errors.Error) {
	var restaurant restaurant.Restaurant
	ctx := DB.Find(&restaurant, id)
	if ctx.RowsAffected == 0 {
		return restaurant, errors.NotFoundError()
	}
	return restaurant, nil
}

func ListRestaurants(accountId uint) []restaurant.Restaurant {
	var restaurants []restaurant.Restaurant = make([]restaurant.Restaurant, 0)
	DB.Where("account_id = ?", accountId).Find(&restaurants)
	return restaurants
}

func ListRestaurantItems(id uint) []restaurant.Item {
	var items []restaurant.Item = make([]restaurant.Item, 0)
	DB.Where("restaurant_id = ?", id).Find(&items)
	return items
}

func UploadItemImage(id uint) string {
	var item restaurant.Item
	DB.Find(&item, id)
	imageId := utils.GenerteId()
	path := "items/" + utils.UintToString(imageId)
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", Bucket, CosClient.Region, path)
	item.Images = append(item.Images, url)
	DB.Save(&item)
	return CosClient.CreatePresignedURL(Bucket, path)
}
