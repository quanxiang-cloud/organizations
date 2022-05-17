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
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/columns"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// Columns api
type Columns struct {
	columns columns.Columns
	log     logger.AdaptedLogger
}

// Open open
func (co *Columns) Open(c *gin.Context) {
	r := new(columns.OpenColumnRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := co.columns.Open(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// Set set will use columns
func (co *Columns) Set(c *gin.Context) {
	r := new(columns.SetUseColumnsRequest)
	profile := header2.GetProfile(c)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.CreatedBy = profile.UserID
	res, err := co.columns.Set(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// GetAll get all columns
func (co *Columns) GetAll(c *gin.Context) {
	r := new(columns.GetAllColumnsRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	all, err := co.columns.GetAll(ginheader.MutateContext(c), r)
	resp.Format(all, err).Context(c)
	return
}

// GetByRoleID get by role
func (co *Columns) GetByRoleID(c *gin.Context) {
	r := new(columns.GetColumnsByRoleRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	all, err := co.columns.GetByRoleID(ginheader.MutateContext(c), r)
	resp.Format(all, err).Context(c)
	return
}

// Update update name
func (co *Columns) Update(c *gin.Context) {
	r := new(columns.UpdateColumnRequest)
	profile := header2.GetProfile(c)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.UpdatedBy = profile.UserID

	res, err := co.columns.Update(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// Add create
func (co *Columns) Add(c *gin.Context) {
	r := new(columns.AddColumnRequest)
	profile := header2.GetProfile(c)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	r.CreatedBy = profile.UserID
	res, err := co.columns.Add(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// Drop del
func (co *Columns) Drop(c *gin.Context) {
	r := new(columns.DropColumnRequest)
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := co.columns.Drop(ginheader.MutateContext(c), r)
	resp.Format(res, err).Context(c)
	return
}

// NewColumnsAPI new
func NewColumnsAPI(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient, log logger.AdaptedLogger) Columns {
	return Columns{
		columns: columns.NewColumns(conf, db, redisClient),
		log:     log,
	}
}
