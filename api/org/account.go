package org

/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/account"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// Account api
type Account struct {
	account     account.Account
	conf        configs.Config
	log         logger.AdaptedLogger
	redisClient redis.UniversalClient
}

// NewAccountAPI new
func NewAccountAPI(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger) Account {
	return Account{
		account:     account.NewAccount(conf, db, redisClient),
		conf:        conf,
		log:         log,
		redisClient: redisClient,
	}
}

// CheckPWD check login pwd
func (a *Account) CheckPWD(c *gin.Context) {
	r := new(account.LoginAccountRequest)
	err := c.ShouldBindJSON(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Header = c.Request.Header.Clone()

	userAccount, err := a.account.CheckPassword(ginheader.MutateContext(c), r)
	if err != nil {
		//todo 记录日志和操作日志
		resp.Format(nil, err).Context(c)
		return
	}
	//todo 记录日志和操作日志
	resp.Format(userAccount, nil).Context(c)
	return
}

// AdminResetPassword admin reset password
func (a *Account) AdminResetPassword(c *gin.Context) {
	r := new(account.AdminUpdatePasswordRequest)
	err := c.ShouldBindJSON(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	r.Header = c.Request.Header.Clone()
	userAccount, err := a.account.AdminUpdatePassword(ginheader.MutateContext(c), r)
	resp.Format(userAccount, err).Context(c)
	return
}

// UserResetPassword user reset password
func (a *Account) UserResetPassword(c *gin.Context) {
	r := new(account.UpdatePasswordRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.UserID = profile.UserID
	userAccount, err := a.account.UpdatePassword(ginheader.MutateContext(c), r)
	resp.Format(userAccount, err).Context(c)
	return
}

// UserFirstResetPassword user first login to reset pasword
func (a *Account) UserFirstResetPassword(c *gin.Context) {
	r := new(account.FirstSetPasswordRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.UserID = profile.UserID
	userAccount, err := a.account.FirstUpdatePassword(ginheader.MutateContext(c), r)
	resp.Format(userAccount, err).Context(c)
	return
}

// LoginGetCode get login code
func (a *Account) LoginGetCode(c *gin.Context) {
	r := new(account.CodeRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Model = a.conf.VerificationCode.LoginCode
	_, err = a.account.GetCode(ginheader.MutateContext(c), r)
	resp.Format(nil, err).Context(c)
	return
}

// ResetPasswordGetCode get reset password code
func (a *Account) ResetPasswordGetCode(c *gin.Context) {
	r := new(account.CodeRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Model = a.conf.VerificationCode.ResetCode
	_, err = a.account.GetCode(ginheader.MutateContext(c), r)
	resp.Format(nil, err).Context(c)
	return
}

// ForgetCode get forget code
func (a *Account) ForgetCode(c *gin.Context) {
	r := new(account.CodeRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Model = a.conf.VerificationCode.ForgetCode
	_, err = a.account.GetCode(ginheader.MutateContext(c), r)
	resp.Format(nil, err).Context(c)
	return
}

// RegisterCode  get register code
func (a *Account) RegisterCode(c *gin.Context) {
	r := new(account.CodeRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Model = a.conf.VerificationCode.RegisterCode
	_, err = a.account.GetCode(ginheader.MutateContext(c), r)
	resp.Format(nil, err).Context(c)
	return
}

// UserForgetResetPassword user reset password when forgot password
func (a *Account) UserForgetResetPassword(c *gin.Context) {
	r := new(account.ForgetResetRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Header = c.Request.Header.Clone()
	userAccount, err := a.account.ForgetUpdatePassword(ginheader.MutateContext(c), r)
	resp.Format(userAccount, err).Context(c)
	return
}
