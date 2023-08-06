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

func CreateItem(accountId, restaurantId uint, item restaurant.Item) restaurant.Item {
	item.RestaurantId = restaurantId
	DB.Save(&item)
	return item
}

func ListRestaurants(accountId uint) []restaurant.Restaurant {
	fmt.Println(accountId)
	var restaurants []restaurant.Restaurant
	DB.Where("account_id = ?", accountId).Find(&restaurants)
	return restaurants
}
