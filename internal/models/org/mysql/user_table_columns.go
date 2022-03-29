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

	"gorm.io/gorm"

	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
)

type userTableColumnsRepo struct {
}

func (u *userTableColumnsRepo) SelectByID(ctx context.Context, db *gorm.DB, id string) (res *org.UserTableColumns) {
	db = db.Where("id=?", id)
	affected := db.Find(&res).
		RowsAffected
	if affected == 1 {
		return res
	}
	return nil
}
func (u *userTableColumnsRepo) SelectByIDAndName(ctx context.Context, db *gorm.DB, id, name string) (res *org.UserTableColumns) {
	db = db.Where("id=? and name=?", id)
	affected := db.Find(&res).
		RowsAffected
	if affected == 1 {
		return res
	}
	return nil
}
func (u *userTableColumnsRepo) Update(ctx context.Context, tx *gorm.DB, req *org.UserTableColumns) (err error) {
	err = tx.Model(req).Updates(req).Error
	return err
}

func (u *userTableColumnsRepo) Delete(ctx context.Context, tx *gorm.DB, id string) (err error) {
	err = tx.Where("id=?", id).Delete(&org.UserTableColumns{}).Error
	return err
}

func (u *userTableColumnsRepo) Insert(ctx context.Context, tx *gorm.DB, req *org.UserTableColumns) (err error) {
	err = tx.Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userTableColumnsRepo) GetAll(ctx context.Context, db *gorm.DB, status int) (list []org.UserTableColumns, total int64) {
	users := make([]org.UserTableColumns, 0)
	var num int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id = ?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	db.Model(&org.UserTableColumns{}).Count(&num)
	affected := db.Find(&users).RowsAffected
	if affected > 0 {
		return users, num
	}
	return nil, 0
}

func (u *userTableColumnsRepo) GetFilter(ctx context.Context, db *gorm.DB, status, attr int) ([]org.UserTableColumns, map[string]string) {
	filter := make(map[string]string)
	useColumns := make([]org.UserTableColumns, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id = ?", tenantID)
	if attr != 0 {
		db = db.Where("attr = ?", attr)
	}
	if status == 0 {
		db = db.Where("id in (select column_id from org_use_columns)")
	} else {
		db = db.Where("id in (select column_id from org_use_columns where viewer_status=?)", status)
	}

	affected := db.Find(&useColumns).RowsAffected
	if affected > 0 {
		for _, v := range useColumns {
			filter[v.ColumnsName] = v.Types
		}
		return useColumns, filter
	}
	return nil, nil
}

func (u *userTableColumnsRepo) GetXlsxField(ctx context.Context, db *gorm.DB, status int) map[string]string {
	fields := make(map[string]string)
	useColumns := make([]org.UserTableColumns, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id = ?", tenantID)
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	affected := db.Where("id in (select column_id from org_use_columns)").Find(&useColumns).RowsAffected

	if affected > 0 {

		for _, v := range useColumns {

			if v.Name == "" {
				return nil
			}
			fields[v.Name] = v.ColumnsName
			switch v.ColumnsName {
			case consts.AVATAR, consts.UPDATEDBY, consts.UPDATEDAT, consts.CREATEDBY, consts.CREATEDAT, consts.LEADERID, consts.USESTATUS, consts.PASSWORDSTATUS, consts.COMPANYID:
				delete(fields, v.Name)

			}

		}
		return fields
	}
	return nil
}

//NewUserTableColumnsRepo new
func NewUserTableColumnsRepo() org.UserTableColumnsRepo {
	return new(userTableColumnsRepo)
}
