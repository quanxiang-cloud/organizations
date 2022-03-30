package mysql

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

	"github.com/quanxiang-cloud/organizations/internal/models/org"
)

type userLeaderRelationRepo struct {
}

func (u *userLeaderRelationRepo) Add(tx *gorm.DB, rq *org.UserLeaderRelation) (err error) {
	err = tx.Create(&rq).Error
	return err
}
func (u *userLeaderRelationRepo) InsertBranch(tx *gorm.DB, req ...org.UserLeaderRelation) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userLeaderRelationRepo) Update(tx *gorm.DB, rq *org.UserLeaderRelation) (err error) {
	err = tx.Model(org.UserLeaderRelation{}).Updates(rq).Error
	return err
}

func (u *userLeaderRelationRepo) DeleteByUserIDs(tx *gorm.DB, userID ...string) (err error) {
	err = tx.Where("user_id in (?)", userID).Delete(org.UserLeaderRelation{}).Error
	return err
}

func (u *userLeaderRelationRepo) SelectByLeaderID(db *gorm.DB, leaderID ...string) []org.UserLeaderRelation {
	relations := make([]org.UserLeaderRelation, 0)
	affected := db.Where("leader_id in(?)", leaderID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}
func (u *userLeaderRelationRepo) SelectByUserIDs(db *gorm.DB, UserID ...string) []org.UserLeaderRelation {
	relations := make([]org.UserLeaderRelation, 0)
	affected := db.Where("user_id in (?)", UserID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}

//NewUserLeaderRelationRepo new
func NewUserLeaderRelationRepo() org.UserLeaderRelationRepo {
	return new(userLeaderRelationRepo)
}
