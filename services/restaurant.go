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

func CreateRestaurant(accountId uint, name, description string) (model.Restaurant, *errors.Error) {
	restaurant := model.Restaurant{
		Name:        name,
		Description: description,
	}
	restaurant.AccountId = accountId
	DB.Save(&restaurant)
	return restaurant, nil
}

func UpdateRestaurant(restaurantId uint, name, description string) (model.Restaurant, *errors.Error) {
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
	restaurant, err := GetRestaurant(id)
	if err != nil {
		return err
	}
	DB.Delete(&restaurant)
	return nil
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
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", Bucket, CosClient.Region))
	if err != nil {
		fmt.Println(1, err)
		return ""
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  CosClient.SecretID,
			SecretKey: CosClient.SecretKey,
		},
	})
	f, err := file.Open()
	if err != nil {
		fmt.Println(2, err)
		return ""
	}
	_, err = client.Object.Put(context.Background(), path, f,
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: file.Header.Get("content-type"),
			},
		})
	if err != nil {
		fmt.Println(3, err)
		return ""
	}
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", Bucket, CosClient.Region, path)
	item.Images = append(item.Images, url)
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

func ListRestaurantTable(restaurantId uint) ([]model.Table, *errors.Error) {
	var tables []model.Table = make([]model.Table, 0)
	if _, err := GetRestaurant(restaurantId); err != nil {
		return tables, err
	}
	DB.Where("restaurant_id = ?", restaurantId).Find(&tables)
	return tables, nil
}
