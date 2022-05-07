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

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/mock"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"github.com/quanxiang-cloud/organizations/pkg/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AccountSuite struct {
	suite.Suite
	conf        configs.Config
	redisClient redis.UniversalClient
	log         logger.AdaptedLogger
	account     Account
	Ctx         context.Context
	db          *gorm.DB
	t           gomock.TestReporter
}

func TestAccount(t *testing.T) {
	d := new(AccountSuite)
	d.t = t
	suite.Run(t, d)
}

func (suite *AccountSuite) SetupTest() {
	conf, err := configs.NewConfig("../../../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), conf)
	ctx := context.Background()
	ctx = header2.SetContext(ctx, user.TenantID, "")
	conn, _, err := sqlmock.New()

	db, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      conn,
	}), &gorm.Config{})
	suite.db = db

	suite.Ctx = ctx
	suite.conf = *conf
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	mock := redismock.NewNiceMock(client)
	suite.redisClient = mock

}

func (suite *AccountSuite) TestCheckPassword() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &LoginAccountRequest{
		UserName: "test1@test.com",
		Password: "654321a..",
		Types:    "pwd",
	}
	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	res, err := suite.account.CheckPassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq.Password = "123456"
	rq.Types = "code"

	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.redisClient.SetEX(suite.Ctx, suite.conf.VerificationCode.LoginCode+":"+rq.UserName, "123456", suite.conf.VerificationCode.ExpireTime*time.Second)
	res, err = suite.account.CheckPassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *AccountSuite) TestUpdatePassword() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().SelectByUserID(gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().UpdatePasswordByUserID(gomock.Any(), gomock.Any()),
	)

	rq := &UpdatePasswordRequest{
		UserID:      "1",
		OldPassword: "654321a..",
		NewPassword: "654321Aa..",
	}
	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	res, err := suite.account.UpdatePassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *AccountSuite) TestFirstUpdatePassword() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().SelectByUserID(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().UpdatePasswordByUserID(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &FirstSetPasswordRequest{
		UserID:      "0",
		OldPassword: "654321a..",
		NewPassword: "654321Aa..",
	}
	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	res, err := suite.account.FirstUpdatePassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *AccountSuite) TestForgetUpdatePassword() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().UpdatePasswordByUserID(gomock.Any(), gomock.Any()),
	)

	rq := &ForgetResetRequest{
		UserName:    "test1@test.com",
		Code:        "123456",
		NewPassword: "654321Aa..",
	}
	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	suite.redisClient.SetEX(suite.Ctx, suite.conf.VerificationCode.ForgetCode+":"+rq.UserName, "123456", suite.conf.VerificationCode.ExpireTime*time.Second)
	res, err := suite.account.ForgetUpdatePassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *AccountSuite) TestAdminUpdatePassword() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().UpdatePasswordByUserID(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)

	rq := &AdminUpdatePasswordRequest{
		UserIDs: []string{"1"},
		SendMessage: []user.SendMessage{
			user.SendMessage{
				UserID:      "1",
				SendChannel: 0,
			},
		},
	}
	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
	}
	res, err := suite.account.AdminUpdatePassword(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *AccountSuite) TestGetCode() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)

	suite.account = &account{
		DB:          suite.db,
		conf:        suite.conf,
		accountRepo: accountRepo,
		user:        userRepo,
		redisClient: suite.redisClient,
		message:     message.NewMessage(suite.conf.InternalNet),
	}

	rq := &CodeRequest{
		UserName: "test1@test.com",
		Model:    "code:login",
	}

	res, err := suite.account.GetCode(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq.Model = "code:reset"
	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)
	res, err = suite.account.GetCode(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq.Model = "code:forget"
	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)
	res, err = suite.account.GetCode(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

	rq.UserName = "testnull@test.com"
	rq.Model = "code:register"
	gomock.InOrder(

		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)
	res, err = suite.account.GetCode(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}
