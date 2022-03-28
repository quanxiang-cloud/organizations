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

//// 联动部门
//func (suite *AllSuite) AddUser() {
//	rq := user.AddUserRequest{
//		DepIDs:       []string{suite.depID},
//		TenantID: suite.tenantID,
//	}
//	rq.UserInfo[consts.USERNAME]="test123445"
//	rq.UserInfo[consts.PHONE]="13628005221"
//	rq.UserInfo[consts.EMAIL]="123@yunify.com"
//	rq.UserInfo[consts.CREATEBY]="adminUserID"
//	res, err := suite.user.AddUser(suite.ctx, &rq)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), res)
//	suite.userID = res.ID
//	suite.userEmail = rq.Email
//	suite.userPhone = rq.Phone
//	rq1 := req.UpdateUser{
//		ID:       suite.userID,
//		UserName: "test123445",
//		UpdatedBy: "adminUserID",
//	}
//	err1 := suite.user.Update(suite.ctx, &rq1)
//	assert.Nil(suite.T(), err1)
//
//	rq2 := req.UpdateUserAvatar{
//		ID:       suite.userID,
//		Avatar:   "图片",
//		UpdatedBy: "adminUserID",
//	}
//	err2 := suite.user.UpdateAvatar(suite.ctx, &rq2)
//	assert.Nil(suite.T(), err2)
//
//	rq3 := req.UserStatus{
//		ID:        suite.userID,
//		UseStatus: 1,
//		UpdatedBy:  "adminUserID",
//	}
//	err3 := suite.user.UpdateUserStatus(suite.ctx, &rq3)
//	assert.Nil(suite.T(), err3)
//
//	res6 := suite.user.UserSelectByID(suite.ctx, suite.userID)
//	assert.NotNil(suite.T(), res6)
//
//	rq7 := req.ListUserReq{
//		IDs: []string{suite.userID},
//	}
//	res7 := suite.user.UserSelectByIDs(suite.ctx, &rq7)
//	assert.NotNil(suite.T(), res7)
//
//	rq8 := req.SearchOneUser{
//		ID: suite.userID,
//	}
//	res8 := suite.user.AdminSelectByID(suite.ctx, &rq8)
//	assert.NotNil(suite.T(), res8)
//
//}
//
//// last
//func (suite *AllSuite) DELUser() {
//	rq := req.UsersStatus{
//		IDS:       []string{suite.userID},
//		UseStatus: -1,
//	}
//	err := suite.user.UpdateUsersStatus(suite.ctx, &rq)
//	assert.Nil(suite.T(), err)
//}
//
//func (suite *AllSuite) TestAdminSelectByCondition() {
//
//	rq := req.SearchListUser{
//		Page:  1,
//		Limit: 100,
//	}
//	res := suite.user.PageList(suite.ctx, &rq)
//	assert.NotNil(suite.T(), res)
//}
//
//func (suite *AllSuite) OtherServerSelectByCondition() {
//	rq := req.SearchListUser{
//		DepID:                suite.depID,
//		IncludeChildDEPChild: 1,
//		Page:                 1,
//		Limit:                10,
//	}
//	res := suite.user.OtherServerSelectByCondition(suite.ctx, &rq)
//	assert.NotNil(suite.T(), res)
//}
//
//func (suite *AllSuite) AdminChangeUsersDEP() {
//	rq := req.ChangeUsersDEP{
//		UsersID:  []string{suite.userID},
//		OldDepID: suite.depID,
//		NewDepID: suite.childID,
//	}
//	err := suite.user.AdminChangeUsersDEP(suite.ctx, &rq)
//	assert.Nil(suite.T(), err)
//}
//
//func (suite *AllSuite) ImportFile() {
//	ctx := logger.ReentryRequestID(context.Background(), "test-UserSelectByID")
//	//open, _ := os.Open("/Users/vvlgo/Desktop/template.xlsx") //本地测试
//	open, _ := os.Open("/root/static/template.xlsx") //gitlab runner 路径
//	defer open.Close()
//	all, _ := io.ReadAll(open)
//	res, _ := suite.user.ImportFile(ctx, all, suite.userID, "adminUserID")
//	assert.NotNil(suite.T(), res)
//	rq := req.UsersStatus{
//		IDS:       res.Success,
//		UseStatus: -1,
//	}
//	err := suite.user.UpdateUsersStatus(ctx, &rq)
//	assert.Nil(suite.T(), err)
//}
