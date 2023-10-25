package controllers

import (
	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	core "github.com/Dparty/core/services"
	model "github.com/Dparty/model/restaurant"
	f "github.com/chenyunda218/golambda"
)

func BillBackward(bill model.Bill) api.Bill {
	return api.Bill{}
}

func OrderBackward(order model.Order) api.Order {
	return api.Order{
		Item: ItemBackward(order.Item),
		Options: f.Map(order.Specification, func(_ int, option model.Pair) api.Pair {
			return api.Pair{
				Left:  option.Left,
				Right: option.Right,
			}
		}),
	}
}

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

func RestaurantBackward(restaurant model.Restaurant) api.Restaurant {
	return api.Restaurant{
		Id:          utils.UintToString(restaurant.ID),
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Tags:        restaurant.Tags,
		Tables: f.Map(restaurant.Tables, func(_ int, table model.Table) api.Table {
			return api.Table{
				Label: table.Label,
				Id:    utils.UintToString(table.ID),
			}
		}),
		Items: f.Map(restaurant.GetItems(), func(_ int, item model.Item) api.Item {
			return ItemBackward(item)
		}),
	}
}

func ItemBackward(item model.Item) api.Item {
	var attributes []api.Attribute = make([]api.Attribute, 0)
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeBackward(a))
	}
	var images = make([]string, 0)
	for _, image := range item.Images {
		images = append(images, image)
	}
	var tags = make([]string, 0)
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
		Printers: f.Map(item.Printers, func(_ int, s uint) string {
			return utils.UintToString(s)
		}),
	}
}

func ItemForward(item api.PutItemRequest) model.Item {
	var attributes []model.Attribute = make([]model.Attribute, 0)
	for _, a := range item.Attributes {
		attributes = append(attributes, AttributeForward(a))
	}
	return model.Item{
		Images:     item.Images,
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: attributes,
		Tags:       item.Tags,
		Printers: f.Map(item.Printers, func(_ int, printer string) uint {
			return utils.StringToUint(printer)
		}),
	}
}

func AttributeForward(attribute api.Attribute) model.Attribute {
	var options []model.Option = make([]model.Option, 0)
	for _, o := range attribute.Options {
		options = append(options, model.Option{
			Label: o.Label,
			Extra: o.Extra,
		})
	}
	return model.Attribute{
		Label:   attribute.Label,
		Options: options,
	}
}

func AttributeBackward(attribute model.Attribute) api.Attribute {
	var options = make([]api.Option, 0)
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
func PrinterBackward(printer model.Printer) api.Printer {
	return api.Printer{
		Id:          utils.UintToString(printer.ID),
		Name:        printer.Name,
		Description: printer.Description,
		Type:        api.PrinterType(printer.Type),
		Sn:          printer.Sn,
	}
}
