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
import "gorm.io/gorm"

// UserDepartmentRelation user department relation
type UserDepartmentRelation struct {
	ID     string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	UserID string `gorm:"column:user_id;type:varchar(64);" json:"userID"`
	DepID  string `gorm:"column:dep_id;type:varchar(64);" json:"depID"`
	Attr   string `gorm:"column:attr;type:varchar(64);" json:"attr"`
}

// TableName tbale name
func (UserDepartmentRelation) TableName() string {
	return "org_user_department_relation"
}

// UserDepartmentRelationRepo interface
type UserDepartmentRelationRepo interface {
	Add(tx *gorm.DB, rq *UserDepartmentRelation) (err error)
	InsertBranch(tx *gorm.DB, req ...UserDepartmentRelation) error
	Update(tx *gorm.DB, rq *UserDepartmentRelation) (err error)
	DeleteByUserIDs(tx *gorm.DB, userID ...string) (err error)
	DeleteByDepIDs(tx *gorm.DB, depID ...string) (err error)
	SelectByDEPID(db *gorm.DB, depID ...string) []UserDepartmentRelation
	SelectByUserIDs(db *gorm.DB, userID ...string) []UserDepartmentRelation
	SelectByUserIDAndDepID(db *gorm.DB, userID, depID string) *UserDepartmentRelation
	DeleteByUserIDAndDepID(db *gorm.DB, userID, depID string) error
}
