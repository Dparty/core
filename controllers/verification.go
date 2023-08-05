package controllers

import (
	"net/http"

	"gitea.svc.boardware.com/bwc/common/constants"
	api "gitea.svc.boardware.com/bwc/core-api"
	"gitea.svc.boardware.com/bwc/core/services"
	"github.com/gin-gonic/gin"
)

type VerificationApi struct{}

var verificationApi VerificationApi

const CREATE_INTERVAL = 60

func (VerificationApi) CreateVerificationCode(c *gin.Context, request api.CreateVerificationCodeRequest) {
	purpose := constants.VerificationCodePurpose(request.Purpose)
	if request.Email == nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	err := services.CreateVerificationCode(*request.Email, purpose)
	if err != nil {
		err.GinHandler(c)
		return
	}
	c.JSON(http.StatusCreated, api.CreateVerificationCodeRespones{
		Email:   request.Email,
		Purpose: request.Purpose,
		Result:  api.SUCCESS_CREATED,
	})
}
