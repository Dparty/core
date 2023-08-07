package controllers

import (
	api "github.com/Dparty/core-api"
	core "github.com/Dparty/core/services"
	restaurant "github.com/Dparty/model/restaurant"

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

func RestaurantBackward(restaurant restaurant.Restaurant) api.Restaurant {
	return api.Restaurant{
		Id:          utils.UintToString(restaurant.ID),
		Name:        restaurant.Name,
		Description: restaurant.Description,
	}
}

func ItemBackward(item restaurant.Item) api.Item {
	var attributes []api.Attribute = make([]api.Attribute, 0)
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeBackward(a))
	}
	return api.Item{
		Id:         utils.UintToString(item.ID),
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: attributes,
	}
}

func ItemForward(item api.PutItemRequest) restaurant.Item {
	var attributes []restaurant.Attribute
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeForward(a))
	}
	return restaurant.Item{
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: attributes,
	}
}

func AttributeForward(attribute api.Attribute) restaurant.Attribute {
	var options []restaurant.Option
	for _, o := range attribute.Options {
		options = append(options, restaurant.Option{
			Label: o.Label,
			Extra: o.Extra,
		})
	}
	return restaurant.Attribute{
		Label:   attribute.Label,
		Options: options,
	}
}

func AttributeBackward(attribute restaurant.Attribute) api.Attribute {
	var options []api.Option
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
