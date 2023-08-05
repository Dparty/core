package services

import (
	"time"

	"gitea.svc.boardware.com/bwc/common/constants"
	"gitea.svc.boardware.com/bwc/common/errors"
	"gitea.svc.boardware.com/bwc/common/utils"
	"gitea.svc.boardware.com/bwc/model/core"
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

func (a Account) Forward() core.Account {
	return core.Account{
		ID:    a.ID,
		Email: a.Email,
		Role:  a.Role,
	}
}

func (a *Account) Backward(account core.Account) *Account {
	a.Email = account.Email
	a.ID = account.ID
	a.Role = account.Role
	return a
}

func CreateSession(email, password string) (*Session, *errors.Error) {
	var account *core.Account
	DB.First(&account, "email = ?", email)
	if account == nil {
		return nil, errors.AuthenticationError()
	}
	if !utils.PasswordsMatch(account.Password, password, account.Salt) {
		return nil, errors.AuthenticationError()
	}
	expiredAt := time.Now().AddDate(0, 0, 7).Unix()
	token, err := utils.SignJwt(
		utils.UintToString(account.ID),
		account.Email,
		string(account.Role),
		expiredAt,
	)
	if err != nil {
		return nil, errors.UndefineError()
	}
	var a Account
	a.Backward(*account)
	return &Session{
		Account:     *a.Backward(*account),
		Token:       token,
		TokenFormat: constants.JWT,
		TokeType:    constants.BEARER,
		ExpiredAt:   expiredAt,
	}, nil
}

func CreateAccount(email, password string, role constants.Role) (*Account, *errors.Error) {
	var accounts []core.Account
	DB.Find(&accounts, "email = ?", email)
	if len(accounts) > 0 {
		return nil, errors.EmailExists()
	}
	hashed, salt := utils.HashWithSalt(password)
	account := Account{
		Email: email,
		Role:  role,
	}.Forward()
	account.Password = hashed
	account.Salt = salt
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
	var coreAccount core.Account
	DB.Find(&coreAccount, id)
	if coreAccount.ID == 0 {
		return nil
	}
	var account Account
	return account.Backward(coreAccount)
}

func GetAccountByEmail(email string) {
	var account core.Account
	DB.First(&account, "email = ?", email)
}

func CreateAccountWithVerificationCode(email, code, password string) (*Account, *errors.Error) {
	verificationCode := GetVerification(email, constants.CREATE_ACCOUNT)
	if verificationCode != nil {
		verificationCode.Tries++
		DB.Save(&verificationCode)
	}
	if verificationCode == nil || verificationCode.Code != code || time.Now().Unix()-verificationCode.CreatedAt.Unix() > EXPIRED_TIME || verificationCode.Tries > 10 {
		return nil, errors.VerificationCodeError()
	}
	return CreateAccount(email, password, constants.USER)
}
