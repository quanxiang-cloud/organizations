package account

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

//import (
//	"context"
//	"github.com/go-redis/redis/v8"
//	"github.com/quanxiang-cloud/cabin/logger"
//	"github.com/quanxiang-cloud/organizations/pkg/configs"
//	"github.com/stretchr/testify/suite"
//	"gorm.io/gorm"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//)
//
//type AccountSuite struct {
//	suite.Suite
//	conf        configs.Config
//	redisClient redis.UniversalClient
//	log         logger.AdaptedLogger
//	account     Account
//	Name        string
//	Ctx         context.Context
//	UserID      string
//}
//
//func NewAccountSuite(ctx context.Context, conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger, name, userID string) *AccountSuite {
//	return &AccountSuite{
//		Ctx:         ctx,
//		conf:        conf,
//		redisClient: redisClient,
//		log:         log,
//		account:     NewAccount(conf, db, redisClient),
//		Name:        name,
//		UserID:      userID,
//	}
//}
//
//func (suite *AccountSuite) CheckPassword() {
//	rq := LoginAccountRequest{
//		UserName: suite.Name,
//		Password: "654321a..",
//		Types:    "pwd",
//	}
//
//	res, err := suite.account.CheckPassword(suite.Ctx, &rq)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), res)
//	rq1 := LoginAccountRequest{
//		UserName: suite.Name,
//		Password: "12345",
//		Types:    "pwd",
//	}
//	res1, err1 := suite.account.CheckPassword(suite.Ctx, &rq1)
//	assert.NotNil(suite.T(), err1)
//	assert.Nil(suite.T(), res1)
//	format := time.Now().Format("20060102")
//	rq2 := LoginAccountRequest{
//		UserName: suite.Name,
//		Password: format,
//		Types:    "code",
//	}
//	res2, err2 := suite.account.CheckPassword(suite.Ctx, &rq2)
//	assert.Nil(suite.T(), err2)
//	assert.NotNil(suite.T(), res2)
//	rq3 := LoginAccountRequest{
//		UserName: suite.Name,
//		Password: "123",
//		Types:    "code",
//	}
//	res3, err3 := suite.account.CheckPassword(suite.Ctx, &rq3)
//	assert.NotNil(suite.T(), err3)
//	assert.Nil(suite.T(), res3)
//}
//
//func (suite *AccountSuite) UpdatePassword() {
//	rq := UpdatePasswordRequest{
//		UserID:      suite.UserID,
//		OldPassword: "654321Aa..",
//		NewPassword: "654321Aa..",
//	}
//	res, err := suite.account.UpdatePassword(suite.Ctx, &rq)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), res)
//	rq1 := UpdatePasswordRequest{
//		UserID:      suite.UserID,
//		OldPassword: "12345",
//		NewPassword: "123456",
//	}
//	res1, err1 := suite.account.UpdatePassword(suite.Ctx, &rq1)
//	assert.NotNil(suite.T(), err1)
//	assert.Nil(suite.T(), res1)
//
//}
//
////TestForgetUpdatePassword 此项测试必须与message服务联测
//func (suite *AccountSuite) ForgetUpdatePassword() {
//	rq := CodeRequest{
//		UserName: suite.Name,
//		Model:    suite.conf.VerificationCode.ForgetCode,
//	}
//	res, err := suite.account.GetCode(suite.Ctx, &rq)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), res)
//
//	rq1 := ForgetResetRequest{
//		UserName:    suite.Name,
//		Code:        res.Code,
//		NewPassword: "654321a..",
//	}
//	res1, err1 := suite.account.ForgetUpdatePassword(suite.Ctx, &rq1)
//	assert.Nil(suite.T(), err1)
//	assert.NotNil(suite.T(), res1)
//}
//
//func (suite *AccountSuite) AdminUpdatePassword() {
//
//	rq := AdminUpdatePasswordRequest{
//		UserIDs:   []string{suite.UserID},
//		CreatedBy: "adminUserID",
//	}
//	res, err := suite.account.AdminUpdatePassword(suite.Ctx, &rq)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), res)
//
//	rq1 := AdminUpdatePasswordRequest{
//		UserIDs:   []string{suite.UserID},
//		CreatedBy: "adminUserID",
//	}
//	res1, err1 := suite.account.AdminUpdatePassword(suite.Ctx, &rq1)
//	assert.Nil(suite.T(), err1)
//	assert.NotNil(suite.T(), res1)
//}
//
//func (suite *AccountSuite) GetCode() {
//	rq := CodeRequest{
//		UserName: suite.Name,
//		Model:    suite.conf.VerificationCode.ResetCode,
//	}
//	res1, err1 := suite.account.GetCode(suite.Ctx, &rq)
//	assert.Nil(suite.T(), err1)
//	assert.NotNil(suite.T(), res1)
//
//	suite.redisClient.Del(suite.Ctx, suite.conf.VerificationCode.ResetCode+":"+suite.Name)
//
//}
