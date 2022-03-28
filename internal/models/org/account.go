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
	"gorm.io/gorm"
)

// Account account
type Account struct {
	ID       string `gorm:"column:id;type:varchar(100);primaryKey ;" json:"id"`            //userID
	Account  string `gorm:"column:account;type:varchar(100);index:account" json:"account"` //多形态:邮箱、手机、其它
	UserID   string `gorm:"column:user_id;type:varchar(64);" json:"userID"`
	Password string `gorm:"column:password;type:varchar(100);" json:"password"`

	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" comment:"创建时间"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" comment:"更新时间"`
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" comment:"删除时间"`
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty" comment:"创建者"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty" comment:"修改者"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty" comment:"删除者"` //删除者
}

// TableName table name
func (Account) TableName() string {
	return "org_user_account"
}

//AccountRepo interface
type AccountRepo interface {
	Insert(tx *gorm.DB, req *Account) error
	InsertBranch(tx *gorm.DB, req ...Account) error
	SelectByAccount(db *gorm.DB, account string) (res *Account)
	SelectByUserID(db *gorm.DB, id string) []Account
	DeleteByID(db *gorm.DB, id ...string) error
	DeleteByUserID(db *gorm.DB, id ...string) error
	Update(tx *gorm.DB, res *Account) error
	UpdatePasswordByUserID(tx *gorm.DB, res *Account) error
}
