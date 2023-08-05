package controllers

import (
	"reflect"

	api "github.com/Dparty/core-api"
	core "github.com/Dparty/core/services"

	"github.com/Dparty/common/utils"
)

func Converter(f any) any {
	typeOfA := reflect.TypeOf(f)
	switch typeOfA.PkgPath() {
	case "github.com/boardware-cloud/controllers/model":
		switch typeOfA.Name() {
		}
	case "github.com/boardware-cloud/services/core":
		switch typeOfA.Name() {
		case "Account":
			return AccountBackward(f.(core.Account))
		case "Session":
			return SessionBackward(f.(core.Session))
		}
	}
	panic("No Converter")
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
