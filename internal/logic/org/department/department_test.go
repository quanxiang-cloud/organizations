package department

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
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/mock"

	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"

	"gorm.io/driver/mysql"
)

type DepartmentSuite struct {
	suite.Suite
	conf        configs.Config
	redisClient redis.Cmdable
	log         logger.AdaptedLogger
	department  Department
	Ctx         context.Context
	db          *gorm.DB
	t           gomock.TestReporter
}

func NewAccountSuite(ctx context.Context, conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger, name, userID string) *DepartmentSuite {
	return &DepartmentSuite{
		Ctx:         ctx,
		conf:        conf,
		redisClient: redisClient,
		log:         log,
	}
}

func TestDepartment(t *testing.T) {
	d := new(DepartmentSuite)
	d.t = t
	suite.Run(t, d)
}

func (suite *DepartmentSuite) SetupTest() {
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
	redisMock := redismock.NewMock()

	suite.redisClient = redisMock
	suite.Ctx = ctx

}

func (suite *DepartmentSuite) TestAdminSelectDEPByCondition() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := AdminSearchListRequest{
		Page:  1,
		Limit: 10,
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().PageList(suite.Ctx, suite.db, 0, 1, 10),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.PageList(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测5
func (suite *DepartmentSuite) TestPageList() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := AdminSearchListRequest{
		Page:  1,
		Limit: 10,
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().PageList(suite.Ctx, suite.db, rq.UseStatus, 1, 10),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.PageList(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 需要联动测试
func (suite *DepartmentSuite) TestUserSelectByID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := SearchOneRequest{
		ID: "1",
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().Get(suite.Ctx, suite.db, rq.ID),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}

	res, err := suite.department.UserSelectByID(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测
func (suite *DepartmentSuite) TestAdminSelectByID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := SearchOneRequest{
		ID: "1",
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().Get(suite.Ctx, suite.db, rq.ID),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.AdminSelectByID(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测
func (suite *DepartmentSuite) TestUserSelectByPID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := SearchListByPIDRequest{
		PID:   "1",
		Page:  1,
		Limit: 10,
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().SelectByPID(suite.Ctx, suite.db, rq.PID, consts.NormalStatus, 1, 10),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.UserSelectByPID(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测5
func (suite *DepartmentSuite) TestAdminSelectByPID() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := SearchListByPIDRequest{
		PID:   "1",
		Page:  1,
		Limit: 10,
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	gomock.InOrder(
		departmentRepo.EXPECT().SelectByPID(suite.Ctx, suite.db, rq.PID, rq.UseStatus, 1, 10),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.AdminSelectByPID(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 联测1
func (suite *DepartmentSuite) TestAddDEP() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := AddRequest{
		Name:      "testadd",
		PID:       "",
		UseStatus: 1,
		CreatBy:   "testUserID",
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		departmentRepo.EXPECT().SelectSupper(suite.Ctx, suite.db),
		departmentRepo.EXPECT().SelectByPIDAndName(suite.Ctx, suite.db, rq.PID, rq.Name),
		//departmentRepo.EXPECT().Get(suite.Ctx, suite.db, rq.PID),
		departmentRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	res, err := suite.department.Add(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 联测1
func (suite *DepartmentSuite) TestUpdateDEP() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := UpdateRequest{
		ID:        "3",
		Name:      "testUpdate",
		PID:       "1",
		UseStatus: 1,
		UpdateBy:  "testUserID",
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(
		departmentRepo.EXPECT().Get(suite.Ctx, suite.db, rq.ID),
		departmentRepo.EXPECT().SelectByPID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		departmentRepo.EXPECT().SelectByPIDAndName(suite.Ctx, suite.db, rq.PID, rq.Name),
		departmentRepo.EXPECT().Get(suite.Ctx, suite.db, gomock.Any()),
		departmentRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()),
		departmentRepo.EXPECT().SelectByPID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		departmentRepo.EXPECT().SelectByPIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(2),
		userDepRepo.EXPECT().SelectByDEPID(gomock.Any(), gomock.Any()).MaxTimes(2),
		userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:          suite.db,
		depRepo:     departmentRepo,
		userDepRepo: userDepRepo,
		userRepo:    userRepo,
	}
	res, err := suite.department.Update(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 联测1
func (suite *DepartmentSuite) TestDelDEP() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()
	rq := DelOneRequest{
		ID: "4",
	}
	departmentRepo := mock.NewMockDepartmentRepo(ctl)
	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)

	gomock.InOrder(
		departmentRepo.EXPECT().Get(suite.Ctx, suite.db, rq.ID),
		departmentRepo.EXPECT().SelectByPID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().DeleteByDepIDs(gomock.Any(), gomock.Any()),
		departmentRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:          suite.db,
		depRepo:     departmentRepo,
		userDepRepo: userDepRepo,
	}
	res, err := suite.department.Delete(suite.Ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *DepartmentSuite) TestTree() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	departmentRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		departmentRepo.EXPECT().PageList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	request := TreeRequest{}
	tree, err := suite.department.Tree(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tree)

}

func (suite *DepartmentSuite) TestGetDepByIDs() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	departmentRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		departmentRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: departmentRepo,
	}
	request := GetByIDsRequest{
		IDs: []string{"1", "2"},
	}
	tree, err := suite.department.GetDepByIDs(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tree)

}

func (suite *DepartmentSuite) TestSetDEPLeader() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(
		userDepRepo.EXPECT().SelectByUserIDAndDepID(gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().Add(gomock.Any(), gomock.Any()),
		//userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:          suite.db,
		userDepRepo: userDepRepo,
		userRepo:    userRepo,
	}
	request := SetDEPLeaderRequest{
		UserID: "1",
		DepID:  "1",
		Attr:   "leader",
	}
	res, err := suite.department.SetDEPLeader(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *DepartmentSuite) TestCancelDEPLeader() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	userDepRepo := mock.NewMockUserDepartmentRelationRepo(ctl)
	userRepo := mock.NewMockUserRepo(ctl)

	gomock.InOrder(
		userDepRepo.EXPECT().SelectByUserIDAndDepID(gomock.Any(), gomock.Any(), gomock.Any()),
		userDepRepo.EXPECT().Update(gomock.Any(), gomock.Any()),
		userRepo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:          suite.db,
		userDepRepo: userDepRepo,
		userRepo:    userRepo,
	}
	request := CancelDEPLeaderRequest{
		UserID: "2",
		DepID:  "2",
		Attr:   "leader",
	}
	tree, err := suite.department.CancelDEPLeader(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tree)

}

func (suite *DepartmentSuite) TestCheckDEPIsExist() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	depRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		depRepo.EXPECT().SelectByPIDAndName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: depRepo,
	}
	request := CheckDEPIsExistRequest{

		DepID:   "2",
		DepName: "test",
	}
	tree, err := suite.department.CheckDEPIsExist(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tree)

}

func (suite *DepartmentSuite) TestGetDepsByIDs() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	depRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		depRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: depRepo,
	}
	request := GetDepsByIDsRequest{
		IDs: []string{"1", "2"},
	}
	res, err := suite.department.GetDepsByIDs(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

func (suite *DepartmentSuite) TestGetMaxGrade() {
	ctl := gomock.NewController(suite.t)
	defer ctl.Finish()

	depRepo := mock.NewMockDepartmentRepo(ctl)

	gomock.InOrder(
		depRepo.EXPECT().GetMaxGrade(gomock.Any(), gomock.Any()),
	)
	suite.department = &department{
		DB:      suite.db,
		depRepo: depRepo,
	}
	request := GetMaxGradeRequest{}
	res, err := suite.department.GetMaxGrade(suite.Ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}
