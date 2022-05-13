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

//UserTableColumns user table column
type UserTableColumns struct {
	ID          string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	Name        string `gorm:"column:name;type:varchar(64)" json:"name"`
	ColumnsName string `gorm:"column:columns_name;type:varchar(64);index:columns_name" json:"columnsName"`
	Types       string `gorm:"column:types;type:varchar(64);" json:"types"`
	Len         int    `gorm:"column:len;" json:"len"`
	PointLen    int    `gorm:"column:point_len;" json:"pointLen"`
	//1:system attr,2:user alias
	Attr int `gorm:"column:attr;" json:"attr"`
	//1:normal,-1:delete
	Status    int    `gorm:"column:status;" json:"status"`
	Format    string `gorm:"column:format;type:varchar(64); " json:"format"`
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" comment:"创建时间"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" comment:"更新时间"`
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" comment:"删除时间"`
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty" comment:"创建者"`
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty" comment:"修改者"`
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty" comment:"删除者"`
}

//TableName table name
func (UserTableColumns) TableName() string {
	return "org_user_table_columns"
}

// UserTableColumnsRepo interface
type UserTableColumnsRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, req *UserTableColumns) (err error)
	Update(ctx context.Context, tx *gorm.DB, req *UserTableColumns) (err error)
	Delete(ctx context.Context, tx *gorm.DB, id string) (err error)
	GetAll(ctx context.Context, db *gorm.DB, status int, name string) (list []UserTableColumns, total int64)
	SelectByID(ctx context.Context, db *gorm.DB, id string) (res *UserTableColumns)
	SelectByIDAndName(ctx context.Context, db *gorm.DB, id, name string) (res *UserTableColumns)
	GetFilter(ctx context.Context, db *gorm.DB, status, attr int) ([]UserTableColumns, map[string]string)
	GetXlsxField(ctx context.Context, db *gorm.DB, status int) map[string]string
	GetByName(ctx context.Context, db *gorm.DB, name string) *UserTableColumns
	GetByColumnName(ctx context.Context, db *gorm.DB, columName string) *UserTableColumns
}
