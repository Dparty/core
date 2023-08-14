package controllers

import (
	api "github.com/Dparty/core-api"
	core "github.com/Dparty/core/services"
	restaurantapi "github.com/Dparty/model/restaurant"

	"github.com/Dparty/common/utils"
)

func AccountBackward(account core.Account) api.Account {
	return api.Account{
		Id:    utils.UintToString(account.ID),
		Email: account.Email,
		Role:  api.Role(account.Role),
	}
}

func SessionBackward(session core.Session) api.Session {
	return api.Session{
		Account:     AccountBackward(session.Account),
		Token:       session.Token,
		TokenType:   string(session.TokeType),
		TokenFormat: string(session.TokenFormat),
		ExpiredAt:   session.ExpiredAt,
	}
}

func PaginationBackward(pagination core.Pagination) api.Pagination {
	return api.Pagination{
		Index: pagination.Index,
		Limit: pagination.Limit,
		Total: pagination.Total,
	}
}

func RestaurantBackward(restaurant restaurantapi.Restaurant) api.Restaurant {
	var tags []string = make([]string, 0)
	// if len(restaurant.Tags) != 0 {
	// 	tags = restaurant.Tags
	// }
	return api.Restaurant{
		Id:          utils.UintToString(restaurant.ID),
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Tags:        tags,
	}
}

func ItemBackward(item restaurantapi.Item) api.Item {
	var attributes []api.Attribute = make([]api.Attribute, 0)
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeBackward(a))
	}
	var images []string = make([]string, 0)
	for _, image := range item.Images {
		images = append(images, image)
	}
	var tags []string = make([]string, 0)
	if len(item.Tags) != 0 {
		tags = item.Tags
	}
	return api.Item{
		Id:         utils.UintToString(item.ID),
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: attributes,
		Images:     images,
		Tags:       tags,
	}
}

func ItemForward(item api.PutItemRequest) restaurantapi.Item {
	var attributes []restaurantapi.Attribute = make([]restaurantapi.Attribute, 0)
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeForward(a))
	}
	return restaurantapi.Item{
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: attributes,
		Tags:       item.Tags,
	}
}

func AttributeForward(attribute api.Attribute) restaurantapi.Attribute {
	var options []restaurantapi.Option = make([]restaurantapi.Option, 0)
	for _, o := range attribute.Options {
		options = append(options, restaurantapi.Option{
			Label: o.Label,
			Extra: o.Extra,
		})
	}
	return restaurantapi.Attribute{
		Label:   attribute.Label,
		Options: options,
	}
}

func AttributeBackward(attribute restaurantapi.Attribute) api.Attribute {
	var options []api.Option = make([]api.Option, 0)
	for _, o := range attribute.Options {
		options = append(options, api.Option{
			Label: o.Label,
			Extra: o.Extra,
		})
	}
	return api.Attribute{
		Label:   attribute.Label,
		Options: options,
	}
}
