package test

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
	"time"

	"github.com/quanxiang-cloud/organizations/internal/logic/org/account"
	"github.com/stretchr/testify/assert"
)

func (suite *AllSuite) CheckPassword() {
	rq := account.LoginAccountRequest{
		UserName: suite.userEmail,
		Password: "654321a..",
		Types:    "pwd",
	}

	res, err := suite.account.CheckPassword(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
	rq1 := account.LoginAccountRequest{
		UserName: suite.userEmail,
		Password: "12345",
		Types:    "pwd",
	}
	res1, err1 := suite.account.CheckPassword(suite.ctx, &rq1)
	assert.NotNil(suite.T(), err1)
	assert.Nil(suite.T(), res1)
	format := time.Now().Format("20060102")
	rq2 := account.LoginAccountRequest{
		UserName: suite.userEmail,
		Password: format,
		Types:    "code",
	}
	res2, err2 := suite.account.CheckPassword(suite.ctx, &rq2)
	assert.Nil(suite.T(), err2)
	assert.NotNil(suite.T(), res2)
	rq3 := account.LoginAccountRequest{
		UserName: suite.userEmail,
		Password: "123",
		Types:    "code",
	}
	res3, err3 := suite.account.CheckPassword(suite.ctx, &rq3)
	assert.NotNil(suite.T(), err3)
	assert.Nil(suite.T(), res3)
}

func (suite *AllSuite) UpdatePassword() {
	rq := account.UpdatePasswordRequest{
		UserID:      suite.userID,
		OldPassword: "654321a..",
		NewPassword: "654321a..",
	}
	res, err := suite.account.UpdatePassword(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
	rq1 := account.UpdatePasswordRequest{
		UserID:      "d5c67872-bcc2-428a-91b3-10f976caa0bc",
		OldPassword: "12345",
		NewPassword: "123456",
	}
	res1, err1 := suite.account.UpdatePassword(suite.ctx, &rq1)
	assert.NotNil(suite.T(), err1)
	assert.Nil(suite.T(), res1)

}

//TestForgetUpdatePassword 此项测试必须与message服务联测
func (suite *AllSuite) ForgetUpdatePassword() {
	rq := account.CodeRequest{
		UserName: suite.userEmail,
		Model:    suite.conf.VerificationCode.ForgetCode,
	}
	res, err := suite.account.GetCode(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq1 := account.ForgetResetRequest{
		UserName:    suite.userEmail,
		Code:        res.Code,
		NewPassword: "654321a..",
	}
	res1, err1 := suite.account.ForgetUpdatePassword(suite.ctx, &rq1)
	assert.Nil(suite.T(), err1)
	assert.NotNil(suite.T(), res1)
}

func (suite *AllSuite) AdminUpdatePassword() {

	rq := account.AdminUpdatePasswordRequest{
		UserIDs:   []string{suite.userID},
		CreatedBy: "adminUserID",
	}
	res, err := suite.account.AdminUpdatePassword(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq1 := account.AdminUpdatePasswordRequest{
		UserIDs:   []string{suite.userID},
		CreatedBy: "adminUserID",
	}
	res1, err1 := suite.account.AdminUpdatePassword(suite.ctx, &rq1)
	assert.Nil(suite.T(), err1)
	assert.NotNil(suite.T(), res1)
}

func (suite *AllSuite) GetCode() {
	rq := account.CodeRequest{
		UserName: suite.userEmail,
		Model:    suite.conf.VerificationCode.ResetCode,
	}
	res1, err1 := suite.account.GetCode(suite.ctx, &rq)
	assert.Nil(suite.T(), err1)
	assert.NotNil(suite.T(), res1)

	suite.redisClient.Del(suite.ctx, suite.conf.VerificationCode.ResetCode+":"+suite.userEmail)

}
