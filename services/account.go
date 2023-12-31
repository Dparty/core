package services

import (
	"time"

	"github.com/Dparty/common/constants"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	"github.com/Dparty/model/core"
)

const EXPIRED_TIME = 60 * 5

type Account struct {
	ID    uint           `json:"id"`
	Email string         `json:"email"`
	Role  constants.Role `json:"role"`
}

type Session struct {
	Account     Account               `json:"account"`
	Token       string                `json:"token"`
	TokeType    constants.TokenType   `json:"tokenType"`
	TokenFormat constants.TokenFormat `json:"tokenFormat"`
	ExpiredAt   int64                 `json:"expiredAt"`
	CreatedAt   int64                 `json:"createdAt"`
}

func (a *Account) Backward(account core.Account) *Account {
	a.Email = account.Email
	a.ID = account.ID
	a.Role = account.Role
	return a
}

func CreateSession(email, password string) (*Session, error) {
	account, err := core.FindAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return nil, fault.ErrUnauthorized
	}
	expiredAt := time.Now().AddDate(0, 0, 7).Unix()
	token, err := utils.SignJwt(
		utils.UintToString(account.ID),
		account.Email,
		string(account.Role),
		expiredAt,
	)
	if err != nil {
		return nil, fault.ErrUndefined
	}
	return &Session{
		Token:       token,
		TokenFormat: constants.JWT,
		TokeType:    constants.BEARER,
		ExpiredAt:   expiredAt,
	}, nil
}

func UpdatePassword(accountId uint, oldPassword, newPassword string) error {
	account, err := core.FindAccount(accountId)
	if err != nil {
		return err
	}
	if !utils.PasswordsMatch(account.Password, oldPassword, account.Salt) {
		return fault.ErrUnauthorized
	}
	hashed, salt := utils.HashWithSalt(newPassword)
	account.Password = hashed
	account.Salt = salt
	DB.Save(&account)
	return nil
}

func UpdatePasswordForce(accountId uint, newPassword string) error {
	account, err := core.FindAccount(accountId)
	if err != nil {
		return err
	}
	hashed, salt := utils.HashWithSalt(newPassword)
	account.Password = hashed
	account.Salt = salt
	DB.Save(&account)
	return nil
}

func CreateAccount(email, password string, role constants.Role) (*Account, error) {
	_, err := core.FindAccountByEmail(email)
	if err == nil {
		return nil, fault.ErrEmailExists
	}
	hashed, salt := utils.HashWithSalt(password)
	account := core.Account{
		Email:    email,
		Role:     role,
		Password: hashed,
		Salt:     salt,
	}
	account.ID = utils.GenerteId()
	if role != "" {
		account.Role = role
	} else {
		account.Role = constants.USER
	}
	DB.Create(&account)
	var back Account
	return back.Backward(account), nil
}

func GetAccountById(id uint) *Account {
	coreAccount, err := core.FindAccount(id)
	if err != nil {
		return nil
	}
	var account Account
	return account.Backward(coreAccount)
}
