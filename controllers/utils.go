package controllers

import (
	"reflect"
	"strings"

	"github.com/Dparty/common/utils"
	api "github.com/Dparty/core-api"
	"github.com/gin-gonic/gin"
)

type Headers struct {
	Authorization string
}

type AuthenticationStatus string

const (
	Authorized   = "Authorized"
	Unauthorized = "Unauthorized"
)

type Authentication struct {
	Status    AuthenticationStatus
	Role      api.Role
	AccountId uint
	Email     string
}

func Authorize(c *gin.Context) Authentication {
	var headers Headers
	c.ShouldBindHeader(&headers)
	authorization := headers.Authorization
	splited := strings.Split(authorization, " ")
	if authorization == "" || len(splited) != 2 {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return AuthorizeByJWT(splited[1])
}

func AuthorizeByJWT(token string) Authentication {
	claims, err := utils.VerifyJwt(token)
	if err != nil {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return Authentication{
		Status:    Authorized,
		Email:     claims["email"].(string),
		AccountId: utils.StringToUint(claims["id"].(string)),
		Role:      api.Role(claims["role"].(string)),
	}
}

func PairsToMap(s []api.Pair) map[string]string {
	output := make(map[string]string)
	for _, option := range s {
		output[option.Left] = option.Right
	}
	return output
}

func SpecificationEqual(a api.Specification, b api.Specification) bool {
	if a.ItemId != b.ItemId {
		return false
	}
	return reflect.DeepEqual(PairsToMap(a.Options), PairsToMap(b.Options))
}

func RemoveFirstOrder(orders []api.Specification, target api.Specification) []api.Specification {

	return orders
}
