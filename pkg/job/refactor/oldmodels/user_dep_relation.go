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

// UserDepartmentRelation user department relation
type UserDepartmentRelation struct {
	UserID string `gorm:"column:user_id;type:varchar(64);" json:"userId"`
	DepID  string `gorm:"column:dep_id;type:varchar(64);" json:"depId"`
}

//TableName 设置表名
func (UserDepartmentRelation) TableName() string {
	return "t_user_department_relation"
}

// UserDepartmentRelationRepo interface
type UserDepartmentRelationRepo interface {
	All(db *gorm.DB) []UserDepartmentRelation
}

type userDepRepo struct {
}

func (d *userDepRepo) All(db *gorm.DB) (one []UserDepartmentRelation) {
	list := make([]UserDepartmentRelation, 0)
	affected := db.Find(&list).RowsAffected
	if affected > 0 {
		return list
	}
	return nil
}

//NewUserDepRepo new
func NewUserDepRepo() UserDepartmentRelationRepo {
	return new(userDepRepo)
}
