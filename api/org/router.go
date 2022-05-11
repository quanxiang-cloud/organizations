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
	"context"
	"github.com/quanxiang-cloud/organizations/pkg/component/publish"

	"github.com/gin-gonic/gin"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/es"
	"github.com/quanxiang-cloud/organizations/pkg/probe"
	"github.com/quanxiang-cloud/organizations/pkg/verification"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router router
type Router struct {
	c configs.Config

	engine *gin.Engine
	bus    *publish.Bus
}

// NewRouter new router
func NewRouter(ctx context.Context, c configs.Config, log logger.AdaptedLogger, bus *publish.Bus) (*Router, error) {
	db, err := mysql.New(c.Mysql, log)
	if err != nil {
		return nil, err
	}
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis.NewClient(c.Redis)
	if err != nil {
		panic(err)
	}
	{
		probe := probe.New(log)
		engine.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		engine.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})

	}

	v1 := engine.Group("/api/v1/org")

	es.New(&c.Elastic, log)

	verification.RegisterValidation()
	userAPI := NewUserAPI(c, db, redisClient, log, bus)

	manage := v1.Group("/m")
	manageUser := manage.Group("/user")
	{

		manageUser.POST("/add", userAPI.Add)
		manageUser.PUT("/update", userAPI.Update)
		manageUser.POST("/list", userAPI.PageList)
		manageUser.GET("/template", userAPI.GetTemplateFile)
		manageUser.GET("/info", userAPI.AdminUserInfo)
		manageUser.PUT("/change/dep", userAPI.AdminChangeUsersDEP)
		manageUser.GET("/index/count", userAPI.IndexCount)
		manageUser.PUT("/group/set", userAPI.UserGroupSet)

	}
	accountAPI := NewAccountAPI(c, db, redisClient, log)
	manageAccount := manage.Group("/account")
	{
		manageAccount.POST("/admin/reset", accountAPI.AdminResetPassword)
	}

	viewer := v1.Group("/h")
	viewerAccount := viewer.Group("/account")
	{

		viewerAccount.GET("/login/code", accountAPI.LoginGetCode)
		viewerAccount.GET("/reset/code", accountAPI.ResetPasswordGetCode)
		viewerAccount.GET("/forget/code", accountAPI.ForgetCode)
		viewerAccount.GET("/register/code", accountAPI.RegisterCode)
		viewerAccount.POST("/check", accountAPI.CheckPWD)
		viewerAccount.POST("/user/reset", accountAPI.UserResetPassword)
		viewerAccount.POST("/user/forget", accountAPI.UserForgetResetPassword)
		viewerAccount.POST("/user/first/reset", accountAPI.UserFirstResetPassword)
	}
	viewerUser := viewer.Group("/user")
	{

		viewerUser.GET("/info", userAPI.UserUserInfo)
		viewerUser.GET("/id", userAPI.UserGetInfo)
		viewerUser.PUT("/update/avatar", userAPI.UpdateAvatar)
		viewerUser.POST("/register", userAPI.Register)
		viewerUser.POST("/ids", userAPI.GetUsersByIDs)
	}

	depAPI := NewDepartmentAPI(c, db, redisClient, log, bus)
	manageDep := manage.Group("/dep")
	{

		manageDep.POST("/add", depAPI.AddDep)
		manageDep.PUT("/update", depAPI.UpdateDep)
		manageDep.GET("/tree", depAPI.DEPTree)
		manageDep.GET("/list", depAPI.PageList)
		manageDep.GET("/group/list", depAPI.GroupPageList)
		manageDep.GET("/info", depAPI.SelectDepByIDAdmin)
		manageDep.GET("/pid", depAPI.SelectDepByPIDAdmin)
		manageDep.PUT("/set/leader", depAPI.SetDEPLeader)
		manageDep.PUT("/cancel/leader", depAPI.CancelDEPLeader)
		manageDep.GET("/check", depAPI.CheckDEPIsExist)

	}

	viewerDep := viewer.Group("/dep")
	{

		viewerDep.GET("/list", depAPI.SelectDepByConditionUser)
		viewerDep.GET("/info", depAPI.SelectDepByIDUser)
		viewerDep.GET("/pid", depAPI.SelectDepByPIDUser)
		viewerDep.POST("/ids", depAPI.GetDepsByIDs)
	}

	columnAPI := NewColumnsAPI(c, db, redisClient, log)
	manageColumn := manage.Group("/column")
	{
		manageColumn.POST("/open", columnAPI.Open)
		manageColumn.GET("/all", columnAPI.GetAll)
		manageColumn.PUT("/set", columnAPI.Set)
		manageColumn.PUT("/update/name", columnAPI.Update)
	}

	oth := v1.Group("/o")
	otherUser := oth.Group("/user")
	{

		otherUser.POST("/info", userAPI.OthGetOneUser)
		otherUser.POST("/update/status", userAPI.UpdateStatus)
		otherUser.POST("/updates/status", userAPI.UpdateUsersStatus)
		otherUser.POST("/user/add", userAPI.OtherServerAddUser)
		otherUser.POST("/department/add", userAPI.OtherServerAddDepartment)
		otherUser.POST("/ids", userAPI.OtherGetUserByIDs)
		otherUser.POST("/dep/id", userAPI.OtherGetUsersByDepID)
	}
	otherDep := oth.Group("/dep")
	{
		otherDep.POST("/ids", depAPI.SelectDepByIDs)
		otherDep.POST("/del", depAPI.DeleteDepByID)
		otherDep.GET("/max/grade", depAPI.GetMaxGrade)
	}
	if err != nil {
		panic(err)
	}

	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c configs.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	engine.Use(ginlogger.LoggerFunc(), ginlogger.LoggerFunc())
	return engine, nil
}

// Run run
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close close
func (r *Router) Close() {
	r.bus.Close()
}
