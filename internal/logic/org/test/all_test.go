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
	"context"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"net/http"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/account"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/department"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

type AllSuite struct {
	suite.Suite
	ctx         context.Context
	user        user.User
	userID      string
	department  department.Department
	depID       string
	account     account.Account
	code        string
	depPID      string
	childID     string
	superPID    string
	leaderID    string
	userEmail   string
	userPhone   string
	tenantID    string
	header      http.Header
	conf        configs.Config
	redisClient redis.UniversalClient
	log         logger.AdaptedLogger
}

func TestAll(t *testing.T) {
	suite.Run(t, new(AllSuite))
}

const (
	requestID = "Request-Id"
)

func (suite *AllSuite) SetupTest() {
	conf, err := configs.NewConfig("../../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), conf)
	ctx := context.Background()
	var name interface{} = requestID
	suite.ctx = context.WithValue(ctx, name, "Test-All")

	//err = logger.New(&configs.Config.Log)
	//assert.Nil(suite.T(), err)

	db, err := mysql.New(conf.Mysql, suite.log)
	assert.Nil(suite.T(), err)

	redisClient, err := redis2.NewClient(conf.Redis)
	assert.Nil(suite.T(), err)
	suite.redisClient = redisClient
	suite.conf = *conf
	header := http.Header{}
	header.Set(requestID, "1000")
	suite.header = header
	suite.department = department.NewDepartment(db)
	suite.user = user.NewUser(*conf, db, redisClient)
	suite.account = account.NewAccount(*conf, db, redisClient)

	suite.leaderID = "testLeaserID"
	suite.tenantID = "1000"
	suite.log = logger.Logger
}

func (suite *AllSuite) TestAll() {
	suite.AddDEP()

	//先删除人再删除部门
	defer suite.SomeAction()
	//defer suite.DELUser()
	// 部门测试
	suite.CheckDEPIsExist()
	suite.AdminSelectByPID()
	suite.UserSelectByPID()
	suite.AdminSelectByID()
	suite.UserSelectByID()
	suite.UserSelectByCondition()
	suite.AdminSelectDEPByCondition()

	// 人员测试 调整逻辑后重写
	//suite.AddUser()
	//suite.AdminChangeUsersDEP()
	//suite.ImportFile()
	//suite.OtherServerSelectByCondition()

	//帐户测试
	suite.AdminUpdatePassword()
	suite.ForgetUpdatePassword()
	suite.UpdatePassword()
	suite.CheckPassword()
	suite.GetCode()

}
