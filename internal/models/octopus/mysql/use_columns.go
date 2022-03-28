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
	"github.com/quanxiang-cloud/organizations/internal/models/octopus"
)

type useColumnsRepo struct {
}

func (u *useColumnsRepo) Update(ctx context.Context, tx *gorm.DB, reqs []octopus.UseColumns) (err error) {
	columns := octopus.UseColumns{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	err = tx.Exec("delete from "+columns.TableName()+" where tenant_id=?", tenantID).Error
	if len(reqs) > 0 {
		for k := range reqs {
			reqs[k].TenantID = tenantID
			err = tx.Create(&reqs[k]).Error
		}

	}
	if err != nil {
		return err
	}
	return nil
}

func (u *useColumnsRepo) SelectAll(ctx context.Context, db *gorm.DB, status int) (res []octopus.UseColumns) {
	data := make([]octopus.UseColumns, 0)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if status != 0 {
		db = db.Where("viewer_status=?", status)
	}
	affected := db.Find(&data).RowsAffected
	if affected > 0 {
		return data
	}
	return nil
}

func (u *useColumnsRepo) DeleteByID(ctx context.Context, tx *gorm.DB, id string) (err error) {
	return tx.Where("column_id=?", id).Delete(&octopus.UseColumns{}).Error
}

//NewUseColumnsRepo 初始化
func NewUseColumnsRepo() octopus.UseColumnsRepo {
	return new(useColumnsRepo)
}
