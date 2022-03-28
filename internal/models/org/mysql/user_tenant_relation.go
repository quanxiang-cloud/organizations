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
	"context"

	"gorm.io/gorm"

	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
)

type userTenantRelationRepo struct {
}

func (u *userTenantRelationRepo) Add(ctx context.Context, tx *gorm.DB, rq *org.UserTenantRelation) error {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	rq.TenantID = tenantID
	err := tx.Create(&rq).Error
	return err
}
func (u *userTenantRelationRepo) InsertBranch(tx *gorm.DB, req ...org.UserTenantRelation) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userTenantRelationRepo) Update(ctx context.Context, tx *gorm.DB, rq *org.UserTenantRelation) error {
	err := tx.Model(org.UserTenantRelation{}).Updates(rq).Error
	return err
}

func (u *userTenantRelationRepo) DeleteByUserIDs(ctx context.Context, tx *gorm.DB, userID ...string) error {
	err := tx.Where("user_id in (?)", userID).Delete(org.UserTenantRelation{}).Error
	return err
}

func (u *userTenantRelationRepo) DeleteByUserIDAndTenantID(ctx context.Context, db *gorm.DB, userID string) error {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	err := db.Where("user_id=? and tenant_id=?", userID, tenantID).Delete(org.UserTenantRelation{}).Error
	return err
}
func (u *userTenantRelationRepo) SelectByUserIDs(ctx context.Context, db *gorm.DB, UserID ...string) []org.UserTenantRelation {
	relations := make([]org.UserTenantRelation, 0)
	affected := db.Where("user_id in (?)", UserID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}

func (u *userTenantRelationRepo) SelectByUserIDAndTenantID(ctx context.Context, db *gorm.DB, userID string) *org.UserTenantRelation {
	relation := &org.UserTenantRelation{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	affected := db.Where("user_id=? and tenant_id=?", userID, tenantID).Find(relation).RowsAffected
	if affected == 1 {
		return relation
	}
	return nil
}

//NewUserTenantRelationRepo new
func NewUserTenantRelationRepo() org.UserTenantRelationRepo {
	return new(userTenantRelationRepo)
}
