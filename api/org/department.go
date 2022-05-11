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
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/pkg/component/publish"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/department"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/other"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// Department api
type Department struct {
	dep    department.Department
	other  other.OthServer
	log    logger.AdaptedLogger
	search *user.Search
	bus    *publish.Bus
}

// NewDepartmentAPI new
func NewDepartmentAPI(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger, bus *publish.Bus) Department {
	return Department{
		dep:    department.NewDepartment(db),
		other:  other.NewOtherServer(conf, db, redisClient, bus),
		log:    log,
		search: user.GetSearch(),
	}
}

//AddDep add
func (d *Department) AddDep(c *gin.Context) {
	r := new(department.AddRequest)
	err := c.ShouldBind(r)
	if err != nil {
		d.log.Error(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	profile := header2.GetProfile(c)
	r.CreatBy = profile.UserID
	res, err := d.dep.Add(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	d.search.PushDep(ginheader.MutateContext(c), nil)
	resp.Format(res, err).Context(c)
	return
}

//UpdateDep update
func (d *Department) UpdateDep(c *gin.Context) {
	r := new(department.UpdateRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	profile := header2.GetProfile(c)
	r.UpdateBy = profile.UserID
	res, err := d.dep.Update(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	if len(res.Users) > 0 {
		d.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	}
	d.search.PushDep(ginheader.MutateContext(c), nil)
	resp.Format(res, err).Context(c)
	return
}

//PageList page select
func (d *Department) PageList(c *gin.Context) {
	r := new(department.AdminSearchListRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.PageList(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//GroupPageList page select
func (d *Department) GroupPageList(c *gin.Context) {
	r := new(department.AdminSearchGroupListRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.GroupPageList(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//SelectDepByConditionUser condition select
func (d *Department) SelectDepByConditionUser(c *gin.Context) {
	r := new(department.ViewerSearchListRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.UserSelectByCondition(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//SelectDepByIDAdmin admin select by id
func (d *Department) SelectDepByIDAdmin(c *gin.Context) {
	r := new(department.SearchOneRequest)
	err := c.ShouldBind(r)
	if err != nil {
		logger.Logger.Errorw(err.Error(), ginlogger.GetRequestID(c))
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.AdminSelectByID(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//SelectDepByIDUser user select by id
func (d *Department) SelectDepByIDUser(c *gin.Context) {
	r := new(department.SearchOneRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.UserSelectByID(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// SelectDepByPIDAdmin adin select by pid
func (d *Department) SelectDepByPIDAdmin(c *gin.Context) {
	r := new(department.SearchListByPIDRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.AdminSelectByPID(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// SelectDepByPIDUser user select by pid
func (d *Department) SelectDepByPIDUser(c *gin.Context) {
	r := new(department.SearchListByPIDRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.UserSelectByPID(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// DeleteDepByID delete by id
func (d *Department) DeleteDepByID(c *gin.Context) {
	r := new(department.DelOneRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := d.dep.Delete(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	d.search.PushDep(ginheader.MutateContext(c), nil)
	resp.Format(res, err).Context(c)
	return
}

// DEPTree get department tree
func (d *Department) DEPTree(c *gin.Context) {
	r := new(department.TreeRequest)

	tree, _ := d.dep.Tree(ginheader.MutateContext(c), r)
	resp.Format(tree, nil).Context(c)
	return
}

// SelectDepByIDs select departments by ids
func (d *Department) SelectDepByIDs(c *gin.Context) {
	r := new(department.GetByIDsRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.GetDepByIDs(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, err).Context(c)
	return
}

// SetDEPLeader set leader of the department
func (d *Department) SetDEPLeader(c *gin.Context) {
	r := new(department.SetDEPLeaderRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.SetDEPLeader(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	if len(res.Users) > 0 {
		d.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	}
	resp.Format(res, err).Context(c)
	return
}

// CancelDEPLeader cacnel leader of the department
func (d *Department) CancelDEPLeader(c *gin.Context) {
	r := new(department.CancelDEPLeaderRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.CancelDEPLeader(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	if len(res.Users) > 0 {
		d.search.PushUser(ginheader.MutateContext(c), nil, res.Users...)
	}
	resp.Format(res, err).Context(c)
	return
}

// CheckDEPIsExist check department exist
func (d *Department) CheckDEPIsExist(c *gin.Context) {
	r := new(department.CheckDEPIsExistRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.CheckDEPIsExist(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

//OtherAll other server get all department data
func (d *Department) OtherAll(c *gin.Context) {
	r := new(other.DepAllRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := d.other.GetAllDeps(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// GetDepsByIDs get deps by ids
func (d *Department) GetDepsByIDs(c *gin.Context) {
	r := new(department.GetDepsByIDsRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.GetDepsByIDs(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// GetMaxGrade get dep max grade
func (d *Department) GetMaxGrade(c *gin.Context) {
	r := new(department.GetMaxGradeRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	res, err := d.dep.GetMaxGrade(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}
