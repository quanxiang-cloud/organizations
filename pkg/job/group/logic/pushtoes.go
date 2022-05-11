package logic

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
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	newmodels "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"gorm.io/gorm"
)

// PushToES clean data from old version
type PushToES interface {
	PushData() error
}

// Data data
type Data struct {
	DB *gorm.DB

	newAccountRepo    org.AccountRepo
	newDepRepo        org.DepartmentRepo
	newUserRepo       org.UserRepo
	newUserDepRepo    org.UserDepartmentRelationRepo
	newUserLeaderRepo org.UserLeaderRelationRepo
	newUserTenantRepo org.UserTenantRelationRepo
	search            *user.Search
}

// NewPushToES new
func NewPushToES(conf *configs.Config, log logger.AdaptedLogger) (*Data, error) {
	db, err := mysql.New(conf.Mysql, log)
	if err != nil {
		return nil, err
	}
	d := &Data{
		DB:          db,
		newUserRepo: newmodels.NewUserRepo(),
	}
	user.NewSearch(db)
	d.search = user.GetSearch()

	return d, nil

}

// PushData PushData data
func (o *Data) PushData() error {
	affected := o.DB.Exec("update org_department set attr=2 where attr is null").RowsAffected
	logger.Logger.Info("push data dep rows:", affected)
	ctx := context.Background()
	ctx = header2.SetContext(ctx, user.TenantID, "")
	list, _ := o.newUserRepo.PageList(ctx, o.DB, 0, 1, 10000, nil)
	if len(list) > 0 {
		u := make(chan int, 1)
		d := make(chan int, 1)

		o.search.PushUser(ctx, u, list...)
		o.search.PushDep(ctx, d)
		var num = 0
		for {
			if num >= 2 {
				fmt.Println("done")
				break
			}
			select {
			case n := <-u:
				num = num + n
			case m := <-d:
				num = num + m
			}
		}
	}
	return nil
}
