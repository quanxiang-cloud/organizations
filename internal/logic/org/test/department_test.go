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
	"github.com/stretchr/testify/assert"

	"github.com/quanxiang-cloud/organizations/internal/logic/org/department"
	"github.com/quanxiang-cloud/organizations/pkg/random2"
)

func (suite *AllSuite) UserSelectByCondition() {
	rq := department.ViewerSearchListRequest{
		Page:  1,
		Limit: 10,
	}
	res, err := suite.department.UserSelectByCondition(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

func (suite *AllSuite) AdminSelectDEPByCondition() {
	rq := department.AdminSearchListRequest{
		Page:  1,
		Limit: 10,
	}
	res, err := suite.department.PageList(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 需要联动测试
func (suite *AllSuite) UserSelectByID() {
	rq := department.SearchOneRequest{
		ID: suite.depID,
	}
	res, err := suite.department.UserSelectByID(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测
func (suite *AllSuite) AdminSelectByID() {
	rq := department.SearchOneRequest{
		ID: suite.depID,
	}
	res, err := suite.department.AdminSelectByID(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测
func (suite *AllSuite) UserSelectByPID() {
	rq := department.SearchListByPIDRequest{
		PID:   suite.depPID,
		Page:  1,
		Limit: 10,
	}
	res, err := suite.department.UserSelectByPID(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)

}

// 联测5
func (suite *AllSuite) AdminSelectByPID() {
	rq := department.SearchListByPIDRequest{
		PID:   suite.depPID,
		Page:  1,
		Limit: 10,
	}
	res, err := suite.department.AdminSelectByPID(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// 联测1
func (suite *AllSuite) AddDEP() {

	rq := department.AddRequest{
		Name:      random2.RandomString(8, 2),
		PID:       "",
		UseStatus: 1,
		CreatBy:   "testUserID",
	}
	res, err := suite.department.Add(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	suite.depID = res.ID
	suite.depPID = res.ID
	suite.superPID = res.ID
	assert.NotNil(suite.T(), res)

	rqChild := department.AddRequest{
		Name:      random2.RandomString(8, 2),
		PID:       suite.depPID,
		UseStatus: 1,
		CreatBy:   "testUserID",
	}
	res1, err1 := suite.department.Add(suite.ctx, &rqChild)
	assert.Nil(suite.T(), err1)
	assert.NotNil(suite.T(), res1)
	suite.childID = res1.ID
}

func (suite *AllSuite) Tree() {
	request := department.TreeRequest{}
	tree, err := suite.department.Tree(suite.ctx, &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tree)

}

// last
func (suite *AllSuite) SomeAction() {

	up := department.UpdateRequest{
		ID:        suite.childID,
		Name:      random2.RandomString(8, 2),
		PID:       suite.depPID,
		UseStatus: 1,
	}
	res, err := suite.department.Update(suite.ctx, &up)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), res)

	setL := department.SetDEPLeaderRequest{
		DepID:  suite.depID,
		UserID: "testUserID",
	}
	res1, err1 := suite.department.SetDEPLeader(suite.ctx, &setL)
	assert.Nil(suite.T(), err1)
	assert.Nil(suite.T(), res1)

	rq3 := department.CancelDEPLeaderRequest{
		DepID: suite.depID,
	}
	res3, err3 := suite.department.CancelDEPLeader(suite.ctx, &rq3)
	assert.Nil(suite.T(), err3)
	assert.Nil(suite.T(), res3)

	rq2 := department.GetByIDsRequest{
		IDs: []string{suite.depID},
	}
	res2, err2 := suite.department.GetDepByIDs(suite.ctx, &rq2)
	assert.Nil(suite.T(), err2)
	assert.NotNil(suite.T(), res2)

	del := department.DelOneRequest{
		ID: suite.depID,
	}
	resL, errL := suite.department.Delete(suite.ctx, &del)
	assert.NotNil(suite.T(), errL)
	assert.NotNil(suite.T(), resL)

	del1 := department.DelOneRequest{
		ID: suite.depID,
	}
	resL1, errL1 := suite.department.Delete(suite.ctx, &del1)
	assert.Nil(suite.T(), errL1)
	assert.Nil(suite.T(), resL1)

	errLTest := suite.department.TestDelete(suite.ctx, &del)
	assert.Nil(suite.T(), errLTest)

	errL1Test := suite.department.TestDelete(suite.ctx, &del1)
	assert.Nil(suite.T(), errL1Test)
}

// 2
func (suite *AllSuite) CheckDEPIsExist() {
	rq := department.CheckDEPIsExistRequest{
		DepID:   suite.depID,
		DepName: "test1",
	}
	res, err := suite.department.CheckDEPIsExist(suite.ctx, &rq)
	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
}
