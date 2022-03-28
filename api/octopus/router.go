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
	"context"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"

	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/probe"
	"github.com/quanxiang-cloud/organizations/pkg/util"
	"github.com/quanxiang-cloud/organizations/pkg/verification"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router 路由
type Router struct {
	c configs.Config

	engine *gin.Engine
}

// NewRouter 开启路由
func NewRouter(ctx context.Context, c configs.Config, log logger.AdaptedLogger) (*Router, error) {
	db, err := mysql.New(c.Mysql, log)
	if err != nil {
		return nil, err
	}

	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	{
		probe := probe.New(util.LoggerFromContext(ctx))
		engine.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		engine.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})

	}
	v1 := engine.Group("/api/v1/org")

	//启动操作记录

	verification.RegisterValidation()
	userAPI := NewUserAPI(c, db, nil, log)

	manage := v1.Group("/m")
	manageUser := manage.Group("/user")
	{
		manageUser.POST("/add", userAPI.Add)
		manageUser.PUT("/update", userAPI.Update)
		manageUser.POST("/list", redirect)
		manageUser.GET("/template", userAPI.GetTemplateFile)
		manageUser.GET("/info", userAPI.AdminUserInfo)
		manageUser.PUT("/change/dep", redirect)
		manageUser.GET("/index/count", redirect)

	}

	manageAccount := manage.Group("/account")
	{
		manageAccount.POST("/admin/reset", redirect)
	}

	//---------------------------用户端用户信息-----------------------
	viewer := v1.Group("/h")
	viewerAccount := viewer.Group("/account")
	{

		viewerAccount.GET("/login/code", redirect)
		viewerAccount.GET("/reset/code", redirect)
		viewerAccount.GET("/forget/code", redirect)
		viewerAccount.GET("/register/code", redirect)
		viewerAccount.POST("/check", redirect)
		viewerAccount.POST("/user/reset", redirect)
		viewerAccount.POST("/user/forget", redirect)
		viewerAccount.POST("/user/first/reset", redirect)
	}
	viewerUser := viewer.Group("/user")
	{

		viewerUser.GET("/info", userAPI.UserUserInfo)
		viewerUser.GET("/id", userAPI.UserGetInfo)
		viewerUser.PUT("/update/avatar", redirect)
		viewerUser.POST("/register", redirect)
		viewerUser.POST("/ids", redirect)
	}

	manageDep := manage.Group("/dep")
	{

		manageDep.POST("/add", redirect)
		manageDep.PUT("/update", redirect)
		manageDep.GET("/tree", redirect)
		manageDep.GET("/list", redirect)
		manageDep.GET("/info", redirect)
		manageDep.GET("/pid", redirect)
		manageDep.PUT("/set/leader", redirect)
		manageDep.PUT("/cancel/leader", redirect)
		manageDep.GET("/check", redirect)

	}

	viewerDep := viewer.Group("/dep")
	{

		viewerDep.GET("/list", redirect)
		viewerDep.GET("/info", redirect)
		viewerDep.GET("/pid", redirect)
		viewerDep.POST("/ids", redirect) //3
	}

	columnAPI := NewColumnsAPI(c, db, nil, log)
	manageColumn := manage.Group("/column")
	{
		manageColumn.POST("/open", columnAPI.Open)
		manageColumn.POST("/add", columnAPI.Add)
		manageColumn.GET("/all", columnAPI.GetAll)
		manageColumn.PUT("/set", columnAPI.Set)
		manageColumn.PUT("/update/name", columnAPI.Update)
	}

	//inner server sdk api
	oth := v1.Group("/o")
	otherUser := oth.Group("/user")
	{
		otherUser.POST("/token", redirect)
		otherUser.POST("/update/status", redirect)
		otherUser.POST("/updates/status", redirect)
		otherUser.POST("/user/add", userAPI.OtherServerAddUser)
		otherUser.POST("/department/add", redirect)
		otherUser.POST("/ids", redirect)
		otherUser.POST("/dep/id", redirect)
	}
	otherDep := oth.Group("/dep")
	{
		otherDep.POST("/ids", redirect)
		otherDep.POST("/del", redirect)
		otherDep.GET("/max/grade", redirect)
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
	engine.Use(ginlogger.LoggerFunc(), ginlogger.RecoveryFunc())
	return engine, nil
}

// Run run
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close close
func (r *Router) Close() {
}

func redirect(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "http://org"+c.Request.URL.Path)
}
