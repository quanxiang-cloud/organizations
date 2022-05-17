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
	"github.com/quanxiang-cloud/organizations/internal/models/octopus"
)

type userTableColumnsRepo struct {
}

func (u *userTableColumnsRepo) SelectByID(ctx context.Context, db *gorm.DB, id string) (res *octopus.UserTableColumns) {
	db = db.Where("id=?", id)
	affected := db.Find(&res).
		RowsAffected
	if affected == 1 {
		return res
	}
	return nil
}
func (u *userTableColumnsRepo) SelectByIDAndName(ctx context.Context, db *gorm.DB, id, name string) (res *octopus.UserTableColumns) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	sql := ""
	if tenantID == "" {
		sql = sql + "(tenant_id='" + tenantID + "' or tenant_id is null)"
	} else {
		sql = sql + "tenant_id='" + tenantID + "'"
	}
	db = db.Where("id=? and name=? and "+sql, id, name)
	affected := db.Find(&res).
		RowsAffected
	if affected == 1 {
		return res
	}
	return nil
}
func (u *userTableColumnsRepo) Update(ctx context.Context, tx *gorm.DB, req *octopus.UserTableColumns) (err error) {
	err = tx.Model(req).Updates(req).Error
	return err
}

func (u *userTableColumnsRepo) Delete(ctx context.Context, tx *gorm.DB, id string) (err error) {
	err = tx.Where("id=?", id).Delete(&octopus.UserTableColumns{}).Error
	return err
}

func (u *userTableColumnsRepo) Insert(ctx context.Context, tx *gorm.DB, req *octopus.UserTableColumns) (err error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	req.TenantID = tenantID
	err = tx.Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userTableColumnsRepo) GetAll(ctx context.Context, db *gorm.DB, status int, name string) (list []octopus.UserTableColumns, total int64) {
	cls := make([]octopus.UserTableColumns, 0)
	var num int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	db.Model(&octopus.UserTableColumns{}).Count(&num)
	affected := db.Find(&cls).RowsAffected
	if affected > 0 {
		return cls, num
	}
	return nil, 0
}

func (u *userTableColumnsRepo) GetFilter(ctx context.Context, db *gorm.DB, attr int, self bool, id ...string) ([]octopus.UserTableColumns, map[string]string) {
	filter := make(map[string]string)
	useColumns := make([]octopus.UserTableColumns, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	db = db.Where("status = ?", consts.NormalStatus)
	if attr != 0 {
		db = db.Where("attr = ?", attr)
	}
	if !self {
		if len(id) > 0 {
			db = db.Where("id in (?)", id)
		}

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
	useColumns := make([]octopus.UserTableColumns, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	columns := octopus.UseColumns{}
	affected := db.Where("id in (select column_id from " + columns.TableName() + ")").Find(&useColumns).RowsAffected

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

func (u *userTableColumnsRepo) GetByName(ctx context.Context, db *gorm.DB, name string) *octopus.UserTableColumns {
	cls := &octopus.UserTableColumns{}
	var num int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}

	db = db.Where("name = ?", name)

	db.Model(&octopus.UserTableColumns{}).Count(&num)
	affected := db.Find(&cls).RowsAffected
	if affected > 0 {
		return cls
	}
	return nil
}

func (u *userTableColumnsRepo) GetByColumnName(ctx context.Context, db *gorm.DB, columName string) *octopus.UserTableColumns {
	cls := &octopus.UserTableColumns{}
	var num int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}

	db = db.Where("columns_name=?", columName)

	db.Model(&octopus.UserTableColumns{}).Count(&num)
	affected := db.Find(&cls).RowsAffected
	if affected > 0 {
		return cls
	}
	return nil
}

//NewUserTableColumnsRepo 初始化
func NewUserTableColumnsRepo() octopus.UserTableColumnsRepo {
	return new(userTableColumnsRepo)
}
