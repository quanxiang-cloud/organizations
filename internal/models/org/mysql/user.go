package mysql

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
	"reflect"
	"strings"

	"gorm.io/gorm"

	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	page2 "github.com/quanxiang-cloud/organizations/pkg/page"
)

type userRepo struct {
}

//NewUserRepo new
func NewUserRepo() org.UserRepo {
	return new(userRepo)
}

func (u *userRepo) Insert(ctx context.Context, tx *gorm.DB, r *org.User) (err error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	r.TenantID = tenantID
	err = tx.Create(r).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) InsertBranch(tx *gorm.DB, req ...*org.User) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) UpdateByID(ctx context.Context, tx *gorm.DB, r *org.User) (err error) {
	err = tx.Model(r).Updates(r).Error
	return err
}

func (u *userRepo) PageList(ctx context.Context, db *gorm.DB, status, page, limit int, userIDs []string) (list []*org.User, total int64) {
	if len(userIDs) > 0 {
		db = db.Where("id in (?)", userIDs)
	}
	if status != 0 {
		db = db.Where("use_status=?", status)
	}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	db = db.Order("updated_at desc")
	users := make([]*org.User, 0)
	var num int64
	db.Model(&org.User{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&users).RowsAffected
	if affected > 0 {
		return users, num
	}

	return nil, 0
}

func (u *userRepo) Get(ctx context.Context, db *gorm.DB, id string) (res *org.User) {
	user := new(org.User)
	affected := db.Model(&org.User{}).Where("id=?", id).Find(&user).RowsAffected
	if affected == 1 {
		return user
	}
	return nil
}

func (u *userRepo) List(ctx context.Context, db *gorm.DB, id ...string) (res []*org.User) {
	users := make([]*org.User, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()

	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	affected := db.Model(&org.User{}).Where("id in (?)", id).Find(&users).RowsAffected
	if affected > 0 {
		return users
	}
	return nil
}

func (u *userRepo) SelectByEmailOrPhone(ctx context.Context, db *gorm.DB, info string) (res *org.User) {
	user := org.User{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	affected := db.Model(&org.User{}).Where("email=? or phone=?", info, info).Find(&user).RowsAffected
	if affected == 1 {
		return &user
	}
	return nil
}

//GetColumns get columns from db
func (u *userRepo) GetColumns(ctx context.Context, db *gorm.DB, user *org.User, schema string) (res []org.Columns) {

	columns := make([]org.MysqlUserColumn, 0)
	affected := db.Raw("select * from information_schema.columns where TABLE_NAME=? and TABLE_SCHEMA =?", user.TableName(), schema).Scan(&columns).RowsAffected
	if affected > 0 {
		v := reflect.TypeOf(user).Elem()
		m := make(map[string]string)
		for i := 0; i < v.NumField(); i++ {
			get := v.Field(i).Tag.Get("gorm")
			split := strings.Split(get, ";")
			i2 := strings.Split(split[0], ":")
			m[i2[1]] = v.Field(i).Tag.Get("comment")
		}
		data := make([]org.Columns, 0)
		for k := range columns {
			columns[k].Name = m[columns[k].ColumnName]
			data = append(data, &columns[k])
		}

		return data
	}
	return nil
}

func (u *userRepo) Count(ctx context.Context, db *gorm.DB, status, activeStatus int) (totalUser, activeUserNum int64) {
	var num1 int64
	var num2 int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	if status != 0 {
		db.Model(&org.User{}).Where("use_status=?", status).Count(&num1)
	}
	if activeStatus != 0 {
		db.Model(&org.User{}).Where("use_status=? and id in (select user_id from org_user_department_relation)", activeStatus).Count(&num2)
	}
	return num1, num2
}
