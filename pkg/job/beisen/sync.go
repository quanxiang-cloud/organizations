package main

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
	"flag"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/es"

	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"github.com/quanxiang-cloud/organizations/pkg/job/beisen/logic"
)

var (
	configPATH = flag.String("config", "configs/config.yml", "-config=配置文件地址")
	isUpdate   = flag.Int("isUpdate", 1, "-isUpdate=是否跟新已存在数据，1更新")
	syncDep    = flag.Int("syncDep", 1, "-syncDep=是否同步组织数据,1同步")
	tenantID   = flag.String("tenantID", "", "-tenantID=租户id")
	requestURL = flag.String("requestURL", "", "-requestURL=数据获取地址http://ip/api/v1/sk/sync/get")
)

func main() {
	flag.Parse()
	conf, err := configs.NewConfig(*configPATH)
	if err != nil {
		panic(err)
	}
	adaptedLogger := logger.Logger
	if err != nil {
		panic(err)
	}
	db, err := mysql.New(conf.Mysql, adaptedLogger)
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}
	client, err := redis2.NewClient(conf.Redis)
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}
	es.New(&conf.Elastic, adaptedLogger)
	sync := &logic.SyncRequest{}
	//任务参数 kube里面参数放在commond里面
	//sync.SyncDEP = 1
	//sync.IsUpdate = 1
	//sync.RequestURL = "http://127.0.0.1:8010/api/v1/sk/sync/refactor/get"
	sync.TenantID = ""
	if *isUpdate != 1 {
		sync.IsUpdate = 0
	} else {
		sync.IsUpdate = *isUpdate
	}
	if *syncDep != 1 {
		sync.SyncDEP = 0
	} else {
		sync.SyncDEP = *syncDep
	}
	if *requestURL == "" {
		logger.Logger.Warn("请求数据地址不能为空")
		return
	}
	sync.RequestURL = *requestURL
	fmt.Println(sync)
	ctx := context.Background()
	ctx = header2.SetContext(ctx, user.TenantID, sync.TenantID)
	newSync := logic.NewSync(*conf, db, client)
	_, err = newSync.SyncData(ctx, sync)
	if err != nil {
		panic(err)
	}
}
