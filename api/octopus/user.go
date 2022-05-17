package octopus

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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/core"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/other"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/user"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// UserAPI api
type UserAPI struct {
	user  user.User
	other other.OthServer
	log   logger.AdaptedLogger
}

// NewUserAPI new
func NewUserAPI(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger) UserAPI {
	return UserAPI{
		user:  user.NewUser(conf, db, redisClient),
		other: other.NewOtherServer(conf, db, redisClient),
		log:   log,
	}
}

// Add add user
func (u *UserAPI) Add(c *gin.Context) {
	data := make(map[string]interface{})
	err := c.ShouldBind(&data)
	if err != nil {
		u.log.Error(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r := new(user.AddUserRequest)
	r.UserInfo = data
	response, err := u.user.Add(ginheader.MutateContext(c), r, c.Request)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	core.DealResponse(c.Writer, response.Response)

}

//Update update user info
func (u *UserAPI) Update(c *gin.Context) {
	data := make(map[string]interface{})
	err := c.ShouldBind(&data)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r := new(user.UpdateUserRequest)
	r.UserInfo = data

	response, err := u.user.Update(ginheader.MutateContext(c), r, c.Request)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	core.DealResponse(c.Writer, response.Response)
}

// GetTemplateFile get file template
func (u *UserAPI) GetTemplateFile(c *gin.Context) {
	r := new(user.GetTemplateFileRequest)
	resp.Format(u.user.Template(ginheader.MutateContext(c), r, c.Request)).Context(c)
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
	response, err := u.user.AdminSelectByID(ginheader.MutateContext(c), r, c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	core.DealResponse(c.Writer, response.Response)
}

// UserUserInfo user gei self info
func (u *UserAPI) UserUserInfo(c *gin.Context) {
	profile := header2.GetProfile(c)
	r := new(user.ViewerSearchOneUserRequest)
	r.ID = profile.UserID

	response, err := u.user.UserSelectByID(ginheader.MutateContext(c), r, c.Request)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	core.DealResponse(c.Writer, response.Response)

}

// UserGetInfo user get other user info
func (u *UserAPI) UserGetInfo(c *gin.Context) {
	r := new(user.ViewerSearchOneUserRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	response, err := u.user.UserSelectByID(ginheader.MutateContext(c), r, c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	core.DealResponse(c.Writer, response.Response)

}

//------------------------inner server---------------------

//OtherServerAddUser other server add user
func (u *UserAPI) OtherServerAddUser(c *gin.Context) {
	r := new(other.AddRequest)
	err := c.ShouldBind(&r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.R = c.Request
	r.W = c.Writer
	res, err := u.other.AddUsers(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}
