package oldmodels

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
import "gorm.io/gorm"

//Account account
type Account struct {
	ID           string `gorm:"column:id;type:varchar(100);primaryKey ;" json:"id"`
	UserName     string `gorm:"column:user_name;type:varchar(100);index:index_user_name" json:"userName"`
	UserID       string `gorm:"column:user_id;type:varchar(64);" json:"userID"`
	Password     string `gorm:"column:password;type:varchar(100);" json:"password"`
	UseStatus    int    `gorm:"column:use_status; " json:"useStatus"`
	CreateTime   int64  `gorm:"column:create_time;type:bigint; " json:"createTime"`
	UpdateTime   int64  `gorm:"column:update_time;type:bigint; " json:"updateTime"`
	CreateBy     string `gorm:"column:create_by;type:varchar(64); " json:"createBy"`
	AuthSourceID string `gorm:"column:auth_source_id;type:varchar(64); " json:"authSourceId"`
}

//TableName 设置表名
func (Account) TableName() string {
	return "user_account"
}

//AccountRepo account
type AccountRepo interface {
	All(db *gorm.DB) []Account
}

type accountRepo struct {
}

func (d *accountRepo) All(db *gorm.DB) (one []Account) {
	list := make([]Account, 0)
	affected := db.Find(&list).RowsAffected
	if affected > 0 {
		return list
	}
	return nil
}

//NewAccountRepo new
func NewAccountRepo() AccountRepo {
	return new(accountRepo)
}
