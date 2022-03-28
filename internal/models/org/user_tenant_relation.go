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

// UserTenantRelation user tenant relation
type UserTenantRelation struct {
	ID       string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	UserID   string `gorm:"column:user_id;type:varchar(64);" json:"userID"`
	TenantID string `gorm:"column:tenant_id;type:varchar(64);" json:"tenantID"`
	Status   int    `gorm:"column:status;" json:"status"`
}

//TableName table name
func (UserTenantRelation) TableName() string {
	return "org_user_tenant_relation"
}

// UserTenantRelationRepo interface
type UserTenantRelationRepo interface {
	Add(ctx context.Context, tx *gorm.DB, rq *UserTenantRelation) error
	InsertBranch(tx *gorm.DB, req ...UserTenantRelation) error
	Update(ctx context.Context, tx *gorm.DB, rq *UserTenantRelation) error
	DeleteByUserIDs(ctx context.Context, tx *gorm.DB, userID ...string) error
	DeleteByUserIDAndTenantID(ctx context.Context, tx *gorm.DB, userID string) error
	SelectByUserIDs(ctx context.Context, db *gorm.DB, userID ...string) []UserTenantRelation
	SelectByUserIDAndTenantID(ctx context.Context, db *gorm.DB, userID string) *UserTenantRelation
}
