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

type userDepartmentRelationRepo struct {
}

func (u userDepartmentRelationRepo) Add(tx *gorm.DB, rq *org.UserDepartmentRelation) (err error) {
	err = tx.Create(&rq).Error
	return err
}
func (u *userDepartmentRelationRepo) InsertBranch(tx *gorm.DB, req ...org.UserDepartmentRelation) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}
func (u userDepartmentRelationRepo) Update(tx *gorm.DB, rq *org.UserDepartmentRelation) (err error) {
	err = tx.Model(rq).Updates(rq).Error
	return err
}

func (u userDepartmentRelationRepo) DeleteByUserIDs(tx *gorm.DB, userID ...string) (err error) {
	err = tx.Where("user_id in (?)", userID).Delete(org.UserDepartmentRelation{}).Error
	return err
}
func (u userDepartmentRelationRepo) DeleteByDepIDs(tx *gorm.DB, depID ...string) (err error) {
	err = tx.Where("dep_id in (?)", depID).Delete(org.UserDepartmentRelation{}).Error
	return err
}

func (u userDepartmentRelationRepo) SelectByDEPID(db *gorm.DB, depID ...string) []org.UserDepartmentRelation {
	relations := make([]org.UserDepartmentRelation, 0)
	affected := db.Where("dep_id in (?)", depID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}
func (u userDepartmentRelationRepo) SelectByUserIDs(db *gorm.DB, UserID ...string) []org.UserDepartmentRelation {
	relations := make([]org.UserDepartmentRelation, 0)
	affected := db.Where("user_id in (?)", UserID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}

func (u userDepartmentRelationRepo) SelectByUserIDAndDepID(tx *gorm.DB, userID, depID string) *org.UserDepartmentRelation {
	relation := &org.UserDepartmentRelation{}
	affected := tx.Model(org.UserDepartmentRelation{}).Where("user_id=? and dep_id=?", userID, depID).Find(relation).RowsAffected
	if affected == 1 {
		return relation
	}
	return nil
}

func (u userDepartmentRelationRepo) DeleteByUserIDAndDepID(tx *gorm.DB, userID, depID string) error {
	relation := &org.UserDepartmentRelation{}
	return tx.Model(org.UserDepartmentRelation{}).Where("user_id=? and dep_id=?", userID, depID).Delete(relation).Error
}

func (u userDepartmentRelationRepo) DeleteByDepIDAndUserIDs(tx *gorm.DB, depID string, userID ...string) error {
	relation := &org.UserDepartmentRelation{}
	return tx.Model(org.UserDepartmentRelation{}).Where("dep_id=? and user_id in (?)", depID, userID).Delete(relation).Error
}

//NewUserDepartmentRelationRepo new
func NewUserDepartmentRelationRepo() org.UserDepartmentRelationRepo {
	return new(userDepartmentRelationRepo)
}
