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
	"github.com/quanxiang-cloud/organizations/internal/logic/common"
	"github.com/quanxiang-cloud/organizations/pkg/component/publish"
	"gorm.io/gorm"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/other"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// UserAPI api
type UserAPI struct {
	user        user.User
	other       other.OthServer
	log         logger.AdaptedLogger
	conf        configs.Config
	redisClient redis.UniversalClient
	search      *user.Search
	bus         *publish.Bus
}

// NewUserAPI new
func NewUserAPI(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger, bus *publish.Bus) UserAPI {
	user.NewSearch(db)
	return UserAPI{
		user:        user.NewUser(conf, db, redisClient),
		other:       other.NewOtherServer(conf, db, redisClient, bus),
		log:         log,
		conf:        conf,
		redisClient: redisClient,
		search:      user.GetSearch(),
		bus:         bus,
	}
}

// Add add
func (u *UserAPI) Add(c *gin.Context) {
	r := new(user.AddUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		u.log.Error(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	r.Password = user.CreatePassword(ginheader.MutateContext(c), u.conf, u.redisClient)
	res, err := u.user.Add(ginheader.MutateContext(c), r)
	if err != nil {
		if strings.Contains(err.Error(), "PRIMARY") {
			resp.Format(nil, error2.New(code.AccountExist)).Context(c)
			return
		}
		resp.Format(nil, err).Context(c)
		return
	}
	//push data to search
	u.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	resp.Format(res, nil).Context(c)
	return
}

// Update update
func (u *UserAPI) Update(c *gin.Context) {
	r := new(user.UpdateUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.Update(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	//if res.UpdateUser != nil {
	//	u.search.PushUser(ginheader.MutateContext(c), nil, res.UpdateUser)
	//}
	//if len(res.Users) > 0 {
	//	u.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	//
	//}
	//if len(res.Spec) > 0 {
	//	common.SendToDapr(ginheader.MutateContext(c), u.bus, res.Spec...)
	//}

	resp.Format(res, nil).Context(c)
	return
}

// UpdateAvatar update avatar
func (u *UserAPI) UpdateAvatar(c *gin.Context) {
	r := new(user.UpdateUserAvatarRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.UpdateAvatar(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	u.search.PushUser(ginheader.MutateContext(c), nil, res.UpdateUser)
	resp.Format(res, nil).Context(c)
	return
}

//PageList page select
func (u *UserAPI) PageList(c *gin.Context) {
	r := new(user.SearchListUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.PageList(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//UpdateStatus update one user use status
func (u *UserAPI) UpdateStatus(c *gin.Context) {
	r := new(user.StatusRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := u.user.UpdateUserStatus(ginheader.MutateContext(c), r)
	if res.User != nil {
		u.search.PushUser(ginheader.MutateContext(c), nil, res.User)

	}
	if len(res.Spec) > 0 {
		common.SendToDapr(ginheader.MutateContext(c), u.bus, res.Spec...)
	}
	resp.Format(res, err).Context(c)
	return
}

//UpdateUsersStatus update list user status
func (u *UserAPI) UpdateUsersStatus(c *gin.Context) {
	r := new(user.ListStatusRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.UpdatedBy = profile.UserID
	res, err := u.user.UpdateUsersStatus(ginheader.MutateContext(c), r)
	if err == nil && len(res.Users) > 0 {
		u.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	}
	resp.Format(res, err).Context(c)
	return
}

// GetTemplateFile get file template
func (u *UserAPI) GetTemplateFile(c *gin.Context) {
	r := new(user.GetTemplateFileRequest)
	err := c.ShouldBind(&r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.Template(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)
	return

}

// AdminUserInfo admin get user info
func (u *UserAPI) AdminUserInfo(c *gin.Context) {
	r := new(user.SearchOneUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := u.user.AdminSelectByID(ginheader.MutateContext(c), r, c.Request)
	resp.Format(res, err).Context(c)
	return
}

// UserUserInfo user get self info
func (u *UserAPI) UserUserInfo(c *gin.Context) {
	profile := header2.GetProfile(c)
	r := new(user.ViewerSearchOneUserRequest)
	r.ID = profile.UserID
	res, err := u.user.UserSelectByID(ginheader.MutateContext(c), r, c.Request)
	resp.Format(res, err).Context(c)
	return
}

// UserGetInfo user get other user info
func (u *UserAPI) UserGetInfo(c *gin.Context) {
	r := new(user.ViewerSearchOneUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := u.user.UserSelectByID(ginheader.MutateContext(c), r, c.Request)
	resp.Format(res, err).Context(c)
	return
}

// AdminChangeUsersDEP admin change user dep
func (u *UserAPI) AdminChangeUsersDEP(c *gin.Context) {
	r := new(user.ChangeUsersDEPRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := u.user.AdminChangeUsersDEP(ginheader.MutateContext(c), r)
	if err == nil && len(res.Users) > 0 {
		u.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	}
	if len(res.Spec) > 0 {
		common.SendToDapr(ginheader.MutateContext(c), u.bus, res.Spec...)
	}
	resp.Format(res, err).Context(c)
	return
}

// OthGetOneUser inner seerver get one user info
func (u *UserAPI) OthGetOneUser(c *gin.Context) {
	r := new(user.TokenUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.OthGetOneUser(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//------------------------other server---------------------

// OtherServerAddUser other server add user
func (u *UserAPI) OtherServerAddUser(c *gin.Context) {
	r := new(other.AddUsersRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.Profile = profile
	res, err := u.other.AddUsers(ginheader.MutateContext(c), r)
	u.other.PushUserToSearch(ginheader.MutateContext(c), nil, nil)
	resp.Format(res, err).Context(c)
	return
}

// OtherServerAddDepartment other server add department
func (u *UserAPI) OtherServerAddDepartment(c *gin.Context) {
	r := new(other.AddDepartmentRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.Profile = profile
	res, err := u.other.AddDepartments(ginheader.MutateContext(c), r)
	u.other.PushDepToSearch(ginheader.MutateContext(c), nil)
	resp.Format(res, err).Context(c)
	return
}

// OtherGetUserByIDs other server get by ids
func (u *UserAPI) OtherGetUserByIDs(c *gin.Context) {
	r := new(other.GetUserByIDsRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)

	r.Profile = profile
	res, err := u.other.GetUserByIDs(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// IndexCount index html get count data
func (u *UserAPI) IndexCount(c *gin.Context) {
	r := new(user.IndexCountRequest)

	resp.Format(u.user.IndexCount(c, r)).Context(c)
	return
}

//OtherUserAll 其它服务全部正常用户的信息
func (u *UserAPI) OtherUserAll(c *gin.Context) {
	r := new(other.UserAllRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.Profile = profile
	resp.Format(u.other.GetAllUsers(ginheader.MutateContext(c), r)).Context(c)
	return
}

// Register register
func (u *UserAPI) Register(c *gin.Context) {
	r := new(user.RegisterRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.Profile = header2.GetProfile(c)
	r.Header = r.Header.Clone()
	res, err := u.user.Register(ginheader.MutateContext(c), r)
	if err == nil && res.User != nil {
		u.search.PushUser(ginheader.MutateContext(c), nil, res.User)
	}
	resp.Format(res, err).Context(c)
	return
}

// OtherGetUsersByDepID 其它服务通过部门id获取下属人员
func (u *UserAPI) OtherGetUsersByDepID(c *gin.Context) {
	r := new(other.GetUsersByDepIDRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := u.other.OtherGetUsersByDepID(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// GetUsersByIDs get users by ids
func (u *UserAPI) GetUsersByIDs(c *gin.Context) {
	r := new(user.GetUsersByIDsRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.GetUsersByIDs(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
}

// UserGroupSet UserGroupSet
func (u *UserAPI) UserGroupSet(c *gin.Context) {
	r := new(user.GroupUserSetRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := u.user.GroupUserSet(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}

	if len(res.Users) > 0 {
		u.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)

	}
	//if len(res.Spec) > 0 {
	//	common.SendToDapr(ginheader.MutateContext(c), u.bus, res.Spec...)
	//}

	resp.Format(res, nil).Context(c)
	return
}
