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
	page2 "github.com/quanxiang-cloud/organizations/pkg/page"
)

type departmentRepo struct {
}

func (d *departmentRepo) Insert(ctx context.Context, tx *gorm.DB, req *org.Department) (err error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	req.TenantID = tenantID
	err = tx.Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *departmentRepo) InsertBranch(tx *gorm.DB, req ...org.Department) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *departmentRepo) Update(ctx context.Context, tx *gorm.DB, req *org.Department) (err error) {
	err = tx.Model(req).Updates(req).Error
	return err
}

func (d *departmentRepo) Delete(ctx context.Context, tx *gorm.DB, id ...string) (err error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	tx = tx.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		tx = tx.Or("tenant_id is null")
	}
	err = tx.Where("id in(?)", id).Delete(&org.Department{}).Error
	return err
}

func (d *departmentRepo) List(ctx context.Context, db *gorm.DB, id ...string) (list []org.Department) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	db = db.Where("id in (?)", id)

	db = db.Order("updated_at desc")
	departments := make([]org.Department, 0)

	affected := db.Find(&departments).RowsAffected
	if affected > 0 {
		return departments
	}
	return nil
}

func (d *departmentRepo) PageList(ctx context.Context, db *gorm.DB, status, page, limit int) (list []org.Department, total int64) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	if status != 0 {
		db = db.Where("use_status=?", status)
	}
	db = db.Order("updated_at desc")
	departments := make([]org.Department, 0)
	var num int64
	db.Model(&org.Department{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&departments).RowsAffected
	if affected > 0 {
		return departments, num
	}

	return nil, 0
}

func (d *departmentRepo) Get(ctx context.Context, db *gorm.DB, id string) (res *org.Department) {
	one := org.Department{}
	db = db.Where("id=?", id)
	affected := db.Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}

func (d *departmentRepo) SelectByPID(ctx context.Context, db *gorm.DB, pid string, status, page, limit int) (list []org.Department, total int64) {
	departments := make([]org.Department, 0)
	var num int64
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	db = db.Where("pid=?", pid)
	if status != 0 {
		db = db.Where("use_status=?", status)
	}
	db.Model(&org.Department{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&departments).RowsAffected
	if affected > 0 {
		return departments, num
	}
	return nil, 0
}

func (d *departmentRepo) SelectByPIDAndName(ctx context.Context, db *gorm.DB, superPID, name string) (one *org.Department) {
	res := org.Department{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	db = db.Where("pid=? and name=? and use_status=1", superPID, name)
	affected := db.Find(&res).RowsAffected
	if affected > 0 {
		return &res
	}
	return nil
}

func (d *departmentRepo) SelectSupper(ctx context.Context, db *gorm.DB) *org.Department {
	res := org.Department{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	db = db.Where("(pid='' or pid is null) and use_status=1")
	affected := db.Find(&res).RowsAffected
	if affected == 1 {
		return &res
	}
	return nil
}

func (d *departmentRepo) Count(ctx context.Context, db *gorm.DB, status int) (total int64) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	if status != 0 {
		db = db.Where("use_status=?", status)
	}
	var num int64
	db.Model(&org.Department{}).Count(&num)
	return num
}

func (d *departmentRepo) GetMaxGrade(ctx context.Context, db *gorm.DB) int64 {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	db = db.Where("tenant_id=?", tenantID)
	if tenantID == "" {
		db = db.Or("tenant_id is null")
	}
	db = db.Select("max(grade) as grade")

	var num int64
	db.Model(&org.Department{}).Find(&num)
	return num
}

//NewDepartmentRepo 初始化
func NewDepartmentRepo() org.DepartmentRepo {
	return new(departmentRepo)
}
