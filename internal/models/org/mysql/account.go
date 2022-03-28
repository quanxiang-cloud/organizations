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

type accountRepo struct {
}

// NewAccountRepo new
func NewAccountRepo() org.AccountRepo {
	return new(accountRepo)
}
func (u *accountRepo) Insert(tx *gorm.DB, req *org.Account) (err error) {
	err = tx.Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *accountRepo) InsertBranch(tx *gorm.DB, req ...org.Account) (err error) {
	err = tx.CreateInBatches(req, len(req)).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *accountRepo) SelectByAccount(db *gorm.DB, account string) *org.Account {
	res := new(org.Account)
	db = db.Where("account=?", account)
	affected := db.Find(&res).
		RowsAffected
	if affected == 1 {
		return res
	}
	return nil
}

func (u *accountRepo) SelectByUserID(db *gorm.DB, id string) []org.Account {
	res := make([]org.Account, 0)
	db = db.Where("user_id=?", id)
	affected := db.Find(&res).
		RowsAffected
	if affected > 0 {
		return res
	}
	return nil
}

func (u *accountRepo) DeleteByID(tx *gorm.DB, id ...string) error {
	err := tx.Where("id in (?)", id).Delete(&org.Account{}).Error
	return err
}

func (u *accountRepo) DeleteByUserID(tx *gorm.DB, id ...string) error {
	err := tx.Where("user_id in (?)", id).Delete(&org.Account{}).Error
	return err
}

func (u *accountRepo) Update(tx *gorm.DB, res *org.Account) error {
	err := tx.Model(res).Updates(res).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *accountRepo) UpdatePasswordByUserID(tx *gorm.DB, res *org.Account) error {
	err := tx.Model(res).Where("user_id=?", res.UserID).Update("password", res.Password).Error
	if err != nil {
		return err
	}
	return nil
}
