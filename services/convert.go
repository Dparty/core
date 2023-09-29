package services

import (
	"github.com/Dparty/model/core"
	"gorm.io/gorm"
)

func AccountForward(a Account) core.Account {
	return core.Account{
		Model: gorm.Model{
			ID: a.ID,
		},
		Email: a.Email,
		Role:  a.Role,
	}
}
