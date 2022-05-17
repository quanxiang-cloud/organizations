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

//UseColumns  use column
type UseColumns struct {
	ID       string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	RoleID   string `gorm:"column:role_id;type:varchar(64)" json:"roleID"`
	ColumnID string `gorm:"column:column_id;type:varchar(64)" json:"columnID"`
	TenantID string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"` //租户id
}

//TableName table name
func (UseColumns) TableName() string {
	return "org_use_columns"
}

//UseColumnsRepo interface
type UseColumnsRepo interface {
	Create(ctx context.Context, tx *gorm.DB, reqs []UseColumns) (err error)
	Update(ctx context.Context, tx *gorm.DB, reqs []UseColumns) (err error)
	SelectAll(ctx context.Context, db *gorm.DB, roleID ...string) (res []UseColumns)
	DeleteByID(ctx context.Context, tx *gorm.DB, id ...string) error
}
