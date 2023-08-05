package services

import (
	"fmt"
	"math/rand"
	"time"

	"gitea.svc.boardware.com/bwc/common/constants"
	"gitea.svc.boardware.com/bwc/common/errors"
	"gitea.svc.boardware.com/bwc/common/utils"
	model "gitea.svc.boardware.com/bwc/model/core"
)

func GetVerification(identity string, purpose constants.VerificationCodePurpose) *model.VerificationCode {
	var verificationCode model.VerificationCode
	ctx := DB.Where("identity = ? AND purpose = ?",
		identity,
		purpose,
	).Order("created_at DESC").Find(&verificationCode)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &verificationCode
}

func CreateVerificationCode(identity string, purpose constants.VerificationCodePurpose) *errors.Error {
	ctx := DB.Where("email = ?", identity).Find(&model.Account{})
	switch {
	case purpose == constants.CREATE_ACCOUNT && ctx.RowsAffected != 0:
		return errors.EmailExists()
	}
	var verificationCode model.VerificationCode
	ctx = DB.Where("identity = ? AND purpose = ?",
		identity,
		purpose,
	).Order("created_at DESC").Find(&verificationCode)
	if ctx.RowsAffected == 0 || time.Now().Unix()-verificationCode.CreatedAt.Unix() >= 60 {
		newCode := &model.VerificationCode{
			Identity: identity,
			Purpose:  purpose,
			Code:     RandomNumberString(6),
		}
		newCode.ID = utils.GenerteId()
		DB.Save(&newCode)
		emailSender.SendHtml("", "Boardware Cloud verification code",
			fmt.Sprintf(`
		<html>
		<body>
			%s
		</body>
		</html>
		`,
				newCode.Code), []string{identity}, []string{}, []string{})
		return nil
	}
	return errors.VerificationCodeFrequent()
}

const charset = "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomNumberString(length int) string {
	return StringWithCharset(length, charset)
}
