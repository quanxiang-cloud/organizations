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
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	client2 "github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/core"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/probe"
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

var httpClient http.Client

// NewRouter 开启路由
func NewRouter(ctx context.Context, c configs.Config, log logger.AdaptedLogger) (*Router, error) {
	httpClient = client2.New(c.InternalNet)

	db, err := mysql.New(c.Mysql, log)
	if err != nil {
		return nil, err
	}

	engine, err := newRouter(c)
	if err != nil {
		return nil, err
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

	//启动操作记录

	verification.RegisterValidation()
	userAPI := NewUserAPI(c, db, nil, log)

	manage := v1.Group("/m")
	manageUser := manage.Group("/user")
	{
		manageUser.POST("/add", userAPI.Add)
		manageUser.PUT("/update", userAPI.Update)
		manageUser.POST("/list", forward(c, httpClient))
		manageUser.GET("/template", userAPI.GetTemplateFile)
		manageUser.GET("/info", userAPI.AdminUserInfo)
		manageUser.PUT("/change/dep", forward(c, httpClient))
		manageUser.GET("/index/count", forward(c, httpClient))
		manageUser.PUT("/group/set", forward(c, httpClient))

	}

	manageAccount := manage.Group("/account")
	{
		manageAccount.POST("/admin/reset", forward(c, httpClient))
	}

	//---------------------------用户端用户信息-----------------------
	viewer := v1.Group("/h")
	viewerAccount := viewer.Group("/account")
	{

		viewerAccount.GET("/login/code", forward(c, httpClient))
		viewerAccount.GET("/reset/code", forward(c, httpClient))
		viewerAccount.GET("/forget/code", forward(c, httpClient))
		viewerAccount.GET("/register/code", forward(c, httpClient))
		viewerAccount.POST("/check", forward(c, httpClient))
		viewerAccount.POST("/user/reset", forward(c, httpClient))
		viewerAccount.POST("/user/forget", forward(c, httpClient))
		viewerAccount.POST("/user/first/reset", forward(c, httpClient))
	}
	viewerUser := viewer.Group("/user")
	{

		viewerUser.GET("/info", userAPI.UserUserInfo)
		viewerUser.GET("/id", userAPI.UserGetInfo)
		viewerUser.PUT("/update/avatar", forward(c, httpClient))
		viewerUser.POST("/register", forward(c, httpClient))
		viewerUser.POST("/ids", forward(c, httpClient))
	}

	manageDep := manage.Group("/dep")
	{

		manageDep.POST("/add", forward(c, httpClient))
		manageDep.PUT("/update", forward(c, httpClient))
		manageDep.GET("/tree", forward(c, httpClient))
		manageDep.GET("/list", forward(c, httpClient))
		manageDep.GET("/info", forward(c, httpClient))
		manageDep.GET("/pid", forward(c, httpClient))
		manageDep.PUT("/set/leader", forward(c, httpClient))
		manageDep.PUT("/cancel/leader", forward(c, httpClient))
		manageDep.GET("/check", forward(c, httpClient))
		manageDep.GET("/group/list", forward(c, httpClient))
		manageDep.DELETE("/group/:id", forward(c, httpClient))

	}

	viewerDep := viewer.Group("/dep")
	{

		viewerDep.GET("/list", forward(c, httpClient))
		viewerDep.GET("/info", forward(c, httpClient))
		viewerDep.GET("/pid", forward(c, httpClient))
		viewerDep.POST("/ids", forward(c, httpClient)) //3
	}

	columnAPI := NewColumnsAPI(c, db, nil, log)
	manageColumn := manage.Group("/column")
	{
		manageColumn.POST("/open", columnAPI.Open)
		manageColumn.POST("/add", columnAPI.Add)
		manageColumn.DELETE("/del", columnAPI.Drop)
		manageColumn.GET("/all", columnAPI.GetAll)
		manageColumn.GET("/all/role", columnAPI.GetByRoleID)
		manageColumn.PUT("/set", columnAPI.Set)
		manageColumn.PUT("/update/name", columnAPI.Update)
	}

	//inner server sdk api
	oth := v1.Group("/o")
	otherUser := oth.Group("/user")
	{
		otherUser.POST("/info", forward(c, httpClient))
		otherUser.POST("/update/status", forward(c, httpClient))
		otherUser.POST("/updates/status", forward(c, httpClient))
		otherUser.POST("/user/add", userAPI.OtherServerAddUser)
		otherUser.POST("/department/add", forward(c, httpClient))
		otherUser.POST("/ids", forward(c, httpClient))
		otherUser.POST("/dep/id", forward(c, httpClient))
	}
	otherDep := oth.Group("/dep")
	{
		otherDep.POST("/ids", forward(c, httpClient))
		otherDep.POST("/del", forward(c, httpClient))
		otherDep.GET("/max/grade", forward(c, httpClient))
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

func forward(conf configs.Config, client http.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Logger.Info("octopus forward org, ", c.Request.URL.Host, c.Request.URL.Path)
		request := c.Request.Clone(c.Request.Context())
		parse, _ := url.ParseRequestURI(conf.OrgHost)
		request.URL = parse
		request.Host = conf.OrgHost
		request.URL.Path = c.Request.URL.Path
		request.RequestURI = ""

		request.URL.RawQuery = c.Request.URL.RawQuery

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(body) > 0 {
				itoa := strconv.Itoa(len(body))
				request.Header.Set("Content-Length", itoa)
				request.ContentLength = int64(len(body))
				request.Body = io.NopCloser(bytes.NewReader(body))
			}

		}

		response, err := client.Do(request)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		core.DealResponse(c.Writer, response)
		return

	}
}
