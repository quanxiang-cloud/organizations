package octopus

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

// Extend alias column table value
type Extend struct {
	ID        string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id,omitempty" comment:"ID"`
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" comment:"删除时间"`
}

// TableName table name
func (e Extend) TableName() string {
	return ""
}

// ExtendRepo db action
type ExtendRepo interface {
	Insert(db, tx *gorm.DB, tableName string, r map[string]interface{}) (err error)
	InsertList(db, tx *gorm.DB, tableName string, r []map[string]interface{}) (err error)
	UpdateByID(db, tx *gorm.DB, tableName string, extend *Extend, r map[string]interface{}) (err error)
	SelectList(db *gorm.DB, tableName string, status, page, limit int) (list []map[string]interface{}, total int64)
	SelectByID(db *gorm.DB, tableName string, id string) (res map[string]interface{})
	SelectByIDs(db *gorm.DB, tableName string, ids []string) (list []map[string]interface{})
}
