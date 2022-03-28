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
	"github.com/quanxiang-cloud/organizations/internal/models/org"
)

type useColumnsRepo struct {
}

func (u *useColumnsRepo) Update(ctx context.Context, tx *gorm.DB, reqs []org.UseColumns) (err error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	err = tx.Exec("delete from org_use_columns where tenant_id=?", tenantID).Error
	if len(reqs) > 0 {
		for k := range reqs {
			err = tx.Create(&reqs[k]).Error
		}

	}
	if err != nil {
		return err
	}
	return nil
}

func (u *useColumnsRepo) SelectAll(ctx context.Context, db *gorm.DB, status int) (res []org.UseColumns) {
	data := make([]org.UseColumns, 0)
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
	return tx.Where("column_id=?", id).Delete(&org.UseColumns{}).Error
}

//NewUseColumnsRepo new
func NewUseColumnsRepo() org.UseColumnsRepo {
	return new(useColumnsRepo)
}
