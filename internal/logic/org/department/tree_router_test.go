package department

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
	"fmt"
	"github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"golang.org/x/net/context"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

func TestTreeRouter(t *testing.T) {
	conf, err := configs.NewConfig("../../../configs/config.yml")
	if err != nil {
		panic(err)
	}
	log := logger.Logger
	dbConn, err := mysql2.New(conf.Mysql, log)
	if err != nil {
		panic(err)
	}
	depRepo := mysql.NewDepartmentRepo()
	ctx := context.Background()
	list, _ := depRepo.PageList(ctx, dbConn, 1, 1, 1000)
	router := NewDepartmentRouter()
	router.AddRoute(list)
	d := "Qingcloud/dep"
	d2 := "Qingcloud/超级部门"
	d3 := "Qingcloud/研发/超级部门"

	node1 := router.GetRoute(d)
	fmt.Println(node1)
	node2 := router.GetRoute(d2)
	fmt.Println(node2)
	node3 := router.GetRoute(d3)
	fmt.Println(node3)
}
