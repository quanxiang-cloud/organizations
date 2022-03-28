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

// Department department
type Department struct {
	ID                 string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	DepartmentName     string `gorm:"column:department_name;type:varchar(64);" json:"departmentName"`
	DepartmentLeaderID string `gorm:"column:department_leader_id;type:varchar(64); " json:"departmentLeaderId"`
	UseStatus          int    `gorm:"column:use_status; " json:"useStatus"`
	PID                string `gorm:"column:pid;type:varchar(64); " json:"pid"`
	SuperPID           string `gorm:"column:super_pid;type:varchar(64); " json:"superId"`
	CompanyID          string `gorm:"column:company_id;type:varchar(64); " json:"companyId"`
	Grade              int    `gorm:"column:grade" form:"grade" json:"grade"`
	CreateTime         int64  `gorm:"column:create_time;type:bigint; " json:"createTime"`
	UpdateTime         int64  `gorm:"column:update_time;type:bigint; " json:"updateTime"`
	CreatBy            string `gorm:"column:creat_by;type:varchar(64); " json:"creatBy"`
}

//TableName 设置表名
func (Department) TableName() string {
	return "t_department"
}

// DepartmentRepo repo
type DepartmentRepo interface {
	All(db *gorm.DB) (list []Department)
}

type departmentRepo struct {
}

// All all
func (d *departmentRepo) All(db *gorm.DB) (one []Department) {
	list := make([]Department, 0)
	affected := db.Find(&list).RowsAffected
	if affected > 0 {
		return list
	}
	return nil
}

//NewDepartmentRepo new
func NewDepartmentRepo() DepartmentRepo {
	return new(departmentRepo)
}
