package org

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
)

// Department department
type Department struct {
	ID        string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	Name      string `gorm:"column:name;type:varchar(64);" json:"name"`
	UseStatus int    `gorm:"column:use_status; " json:"useStatus"`               //1正常，-1真删除，-2禁用
	Attr      int    `gorm:"column:attr; " json:"attr"`                          //1公司，2部门，3组织
	PID       string `gorm:"column:pid;type:varchar(64); " json:"pid"`           //上层ID
	SuperPID  string `gorm:"column:super_pid;type:varchar(64); " json:"superId"` //最顶层父级ID
	Grade     int    `gorm:"column:grade" form:"grade" json:"grade"`             //部门等级
	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" comment:"创建时间"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" comment:"更新时间"`
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" comment:"删除时间"`
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty" comment:"创建者"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty" comment:"修改者"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty" comment:"删除者"` //删除者
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`                           //租户id
}

//TableName table name
func (Department) TableName() string {
	return "org_department"
}

// DepartmentRepo interface
type DepartmentRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, req *Department) (err error)
	InsertBranch(tx *gorm.DB, req ...Department) error
	Update(ctx context.Context, tx *gorm.DB, req *Department) (err error)
	Delete(ctx context.Context, tx *gorm.DB, id ...string) (err error)
	List(ctx context.Context, db *gorm.DB, attr []int, id ...string) (list []Department)
	PageList(ctx context.Context, db *gorm.DB, status, page, limit int, attr []int) (list []Department, total int64)
	Get(ctx context.Context, db *gorm.DB, id string) (res *Department)
	SelectByPID(ctx context.Context, db *gorm.DB, pid string, status, page, limit int) (list []Department, total int64)
	SelectByPIDAndName(ctx context.Context, db *gorm.DB, pid, name string) (one *Department)
	SelectByPIDs(ctx context.Context, db *gorm.DB, status int, pid ...string) (one []Department)
	SelectSupper(ctx context.Context, db *gorm.DB, attr []int) *Department
	Count(ctx context.Context, db *gorm.DB, status int) int64
	GetMaxGrade(ctx context.Context, db *gorm.DB) int64
}
