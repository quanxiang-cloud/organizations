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
//
package user

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/mock"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type UserSuite struct {
	suite.Suite
	Ctx         context.Context
	user        User
	db          *gorm.DB
	redisClient redis.UniversalClient
	log         logger.AdaptedLogger
	conf        configs.Config
	t           gomock.TestReporter
}

func TestUser(t *testing.T) {
	d := new(UserSuite)
	d.t = t
	suite.Run(t, d)
}

func (suite *UserSuite) SetupTest() {
	conf, err := configs.NewConfig("../../../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), conf)
	ctx := context.Background()
	ctx = header2.SetContext(ctx, TenantID, "")
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

// 联动部门
func (suite *UserSuite) TestAddUser() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	userTenantRepo := mock.NewMockUserTenantRelationRepo(ctl)
	gomock.InOrder(
		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().Add(gomock.Any(), gomock.Any()),
		userLeaderRepo.EXPECT().Add(gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().Insert(gomock.Any(), gomock.Any()),
		userTenantRepo.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	rq := &AddUserRequest{
		Name:      "SuiteTestName",
		Email:     "SuiteTestName@yunify.com",
		Phone:     "13666668888",
		SelfEmail: "SuiteTestName@yunify.com",
		SendMessage: SendMessage{
			SendChannel: 0,
		},
		Dep: []DepRequest{{
			DepID: "1",
			Attr:  "test",
		}},
		Leader: []LeaderRequest{{
			UserID: "1",
			Attr:   "test",
		}},
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		accountReo:     accountRepo,
		userDepRepo:    userDepRepo,
		userTenantRepo: userTenantRepo,
		userLeaderRepo: userLeaderRepo,
	}
	res, err := suite.user.Add(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 联动部门
func (suite *UserSuite) TestUpdateUser() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	accountRepo := mock.NewMockAccountRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	userTenantRepo := mock.NewMockUserTenantRelationRepo(ctl)
	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().DeleteByUserIDs(gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().Add(gomock.Any(), gomock.Any()).AnyTimes(),
		userLeaderRepo.EXPECT().DeleteByUserIDs(gomock.Any(), gomock.Any()),
		userLeaderRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()).AnyTimes(),
		userLeaderRepo.EXPECT().Add(gomock.Any(), gomock.Any()).AnyTimes(),
		userLeaderRepo.EXPECT().SelectByLeaderID(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
	)
	rq := &UpdateUserRequest{
		ID:        "2",
		Name:      "test1",
		Email:     "test11@test.com",
		Phone:     "13666668888",
		SelfEmail: "SuiteTestName@yunify.com",
		Dep: []DepRequest{{
			DepID: "1",
			Attr:  "test",
		}},
		Leader: []LeaderRequest{{
			UserID: "1",
			Attr:   "test",
		}},
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		accountReo:     accountRepo,
		userDepRepo:    userDepRepo,
		userTenantRepo: userTenantRepo,
		userLeaderRepo: userLeaderRepo,
	}
	res, err := suite.user.Update(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestUpdateAvatar() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		userRepo.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	rq := &UpdateUserAvatarRequest{
		ID:     "2",
		Avatar: "avatar",
	}
	suite.user = &user{
		DB:       suite.db,
		userRepo: userRepo,
	}
	res, err := suite.user.UpdateAvatar(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestPageList() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(
		userDepRepo.EXPECT().SelectByDEPID(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().PageList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &SearchListUserRequest{
		DepID: "1",
		Page:  1,
		Limit: 100,
	}
	suite.user = &user{
		DB:          suite.db,
		userRepo:    userRepo,
		userDepRepo: userDepRepo,
		depRepo:     depRepo,
	}
	res, err := suite.user.PageList(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestAdminSelectByID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		columnRepo.EXPECT().GetFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		depRepo.EXPECT().PageList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()),
		depRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
		userLeaderRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()).AnyTimes(),
	)

	rq := &SearchOneUserRequest{
		ID: "1",
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
	}
	res, err := suite.user.AdminSelectByID(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestUserSelectByID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		columnRepo.EXPECT().GetFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		depRepo.EXPECT().PageList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()),
		depRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
		userLeaderRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()).AnyTimes(),
	)

	rq := &ViewerSearchOneUserRequest{
		ID: "1",
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
	}
	res, err := suite.user.UserSelectByID(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestUpdateUserStatus() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		userRepo.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()),
	)

	rq := &StatusRequest{
		ID:        "1",
		UseStatus: 1,
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.UpdateUserStatus(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestUpdateUsersStatus() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
		accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes(),
	)

	rq := &ListStatusRequest{
		IDS:       []string{"1"},
		UseStatus: 1,
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.UpdateUsersStatus(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestAdminChangeUsersDEP() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		userDepRepo.EXPECT().SelectByUserIDAndDepID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
		userDepRepo.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes(),
		userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &ChangeUsersDEPRequest{
		UsersID:  []string{"1"},
		OldDepID: "1",
		NewDepID: "2",
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.AdminChangeUsersDEP(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestOthGetOneUser() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
		depRepo.EXPECT().PageList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()).AnyTimes(),
		depRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),

		userLeaderRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()).AnyTimes(),
	)

	rq := &TokenUserRequest{
		ID: "1",
	}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.OthGetOneUser(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestTemplate() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		columnRepo.EXPECT().GetXlsxField(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &GetTemplateFileRequest{}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.Template(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestIndexCount() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)

	gomock.InOrder(
		userRepo.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		depRepo.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &IndexCountRequest{}
	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
	}
	res, err := suite.user.IndexCount(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestRegister() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)
	userTenantRepo := mock.NewMockUserTenantRelationRepo(ctl)
	mockLandlord := mock.NewMockLandlord(ctl)
	gomock.InOrder(
		accountRepo.EXPECT().SelectByAccount(gomock.Any(), gomock.Any()),
		mockLandlord.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()),
		accountRepo.EXPECT().Insert(gomock.Any(), gomock.Any()),
		userTenantRepo.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &RegisterRequest{
		Name:     "test1213",
		Email:    "test1213@test.com",
		Password: "123456Aa..",
		Code:     "123456",
	}

	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
		landlord:       mockLandlord,
		userTenantRepo: userTenantRepo,
		conf:           suite.conf,
	}
	suite.redisClient.SetEX(suite.Ctx, suite.conf.VerificationCode.RegisterCode+":"+rq.Email, "123456", suite.conf.VerificationCode.ExpireTime*time.Second)
	res, err := suite.user.Register(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

func (suite *UserSuite) TestGetUsersByIDs() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	columnRepo := mock.NewMockUserTableColumnsRepo(ctl)
	depRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)
	userLeaderRepo := mock.NewMockUserLeaderRelationRepo(ctl)
	accountRepo := mock.NewMockAccountRepo(ctl)
	userTenantRepo := mock.NewMockUserTenantRelationRepo(ctl)
	mockLandlord := mock.NewMockLandlord(ctl)
	gomock.InOrder(
		userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().SelectByUserIDs(gomock.Any(), gomock.Any()),
		depRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	rq := &GetUsersByIDsRequest{
		IDs: []string{"1"},
	}

	suite.user = &user{
		DB:             suite.db,
		userRepo:       userRepo,
		userDepRepo:    userDepRepo,
		depRepo:        depRepo,
		columnRepo:     columnRepo,
		userLeaderRepo: userLeaderRepo,
		accountReo:     accountRepo,
		redisClient:    suite.redisClient,
		landlord:       mockLandlord,
		userTenantRepo: userTenantRepo,
	}
	res, err := suite.user.GetUsersByIDs(suite.Ctx, rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}
