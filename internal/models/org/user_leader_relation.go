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

// UserLeaderRelation user leader relaiton
type UserLeaderRelation struct {
	ID       string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	UserID   string `gorm:"column:user_id;type:varchar(64);" json:"userID"`
	LeaderID string `gorm:"column:leader_id;type:varchar(64);" json:"leaderID"`
	Attr     string `gorm:"column:attr;" json:"attr"`
}

//TableName table name
func (UserLeaderRelation) TableName() string {
	return "org_user_leader_relation"
}

// UserLeaderRelationRepo interface
type UserLeaderRelationRepo interface {
	Add(tx *gorm.DB, rq *UserLeaderRelation) (err error)
	InsertBranch(tx *gorm.DB, req ...UserLeaderRelation) error
	Update(tx *gorm.DB, rq *UserLeaderRelation) (err error)
	DeleteByUserIDs(tx *gorm.DB, userID ...string) (err error)
	SelectByLeaderID(db *gorm.DB, leaderID string) []UserLeaderRelation
	SelectByUserIDs(db *gorm.DB, userID ...string) []UserLeaderRelation
}
