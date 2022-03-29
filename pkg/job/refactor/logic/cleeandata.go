package logic

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
	"fmt"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	newmodels "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	"github.com/quanxiang-cloud/organizations/pkg/job/refactor/oldmodels"
	"gorm.io/gorm"
)

// CleanDataV1 clean data from old version
type CleanDataV1 interface {
	CleanV1() error
}

// Data data
type Data struct {
	DB             *gorm.DB
	oldAccountRepo oldmodels.AccountRepo
	oldDepRepo     oldmodels.DepartmentRepo
	oldUserRepo    oldmodels.UserRepo
	oldUserDepRepo oldmodels.UserDepartmentRelationRepo

	newAccountRepo    org.AccountRepo
	newDepRepo        org.DepartmentRepo
	newUserRepo       org.UserRepo
	newUserDepRepo    org.UserDepartmentRelationRepo
	newUserLeaderRepo org.UserLeaderRelationRepo
	newUserTenantRepo org.UserTenantRelationRepo
	search            *user.Search
}

// NewCleanV1 new
func NewCleanV1(conf *configs.Config, log logger.AdaptedLogger) (*Data, error) {
	db, err := mysql.New(conf.Mysql, log)
	if err != nil {
		return nil, err
	}
	d := &Data{
		DB:             db,
		oldAccountRepo: oldmodels.NewAccountRepo(),
		oldDepRepo:     oldmodels.NewDepartmentRepo(),
		oldUserRepo:    oldmodels.NewUserRepo(),
		oldUserDepRepo: oldmodels.NewUserDepRepo(),

		newDepRepo:        newmodels.NewDepartmentRepo(),
		newAccountRepo:    newmodels.NewAccountRepo(),
		newUserRepo:       newmodels.NewUserRepo(),
		newUserDepRepo:    newmodels.NewUserDepartmentRelationRepo(),
		newUserLeaderRepo: newmodels.NewUserLeaderRelationRepo(),
		newUserTenantRepo: newmodels.NewUserTenantRelationRepo(),
	}
	user.NewSearch(db, d.newUserRepo, d.newUserLeaderRepo, d.newUserDepRepo, d.newDepRepo)
	d.search = user.GetSearch()

	return d, nil

}

// CleanV1 clean func
func (o *Data) CleanV1() error {
	var err error = nil
	tx := o.DB.Begin()
	err = o.cleanAccount(tx)
	err = o.cleanUser(tx)
	err = o.cleanDep(tx)
	if err != nil {
		logger.Logger.Info(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}

	ctx := context.Background()
	ctx = header2.SetContext(ctx, user.TenantID, "")
	list, _ := o.newUserRepo.PageList(ctx, o.DB, 0, 1, 10000, nil)
	if len(list) > 0 {
		u := make(chan int, 1)
		d := make(chan int, 1)

		o.search.PushUser(ctx, u, list...)
		o.search.PushDep(ctx, d)
		var num = 0
		for {
			if num >= 2 {
				fmt.Println("done")
				break
			}
			select {
			case n := <-u:
				num = num + n
			case m := <-d:
				num = num + m
			}
		}
	}
	return nil
}

func (o *Data) cleanAccount(tx *gorm.DB) error {
	allOldAccouts := o.oldAccountRepo.All(o.DB)
	newAccounts := make([]org.Account, 0)
	for k := range allOldAccouts {
		account := org.Account{}
		account.ID = allOldAccouts[k].ID
		account.UserID = allOldAccouts[k].UserID
		account.Password = allOldAccouts[k].Password
		account.Account = allOldAccouts[k].UserName
		account.CreatedAt = allOldAccouts[k].CreateTime
		account.CreatedBy = allOldAccouts[k].CreateBy

		newAccounts = append(newAccounts, account)
	}
	if len(newAccounts) > 0 {
		err := o.newAccountRepo.InsertBranch(tx, newAccounts...)
		return err
	}
	return nil
}

func (o *Data) cleanDep(tx *gorm.DB) error {
	allOldDeps := o.oldDepRepo.All(o.DB)
	departments := make([]org.Department, 0)
	for k := range allOldDeps {
		department := org.Department{}
		department.ID = allOldDeps[k].ID
		department.Name = allOldDeps[k].DepartmentName
		department.UseStatus = allOldDeps[k].UseStatus
		department.PID = allOldDeps[k].PID
		department.SuperPID = allOldDeps[k].SuperPID
		department.Grade = allOldDeps[k].Grade
		department.CreatedAt = allOldDeps[k].CreateTime
		department.CreatedBy = allOldDeps[k].CreatBy
		departments = append(departments, department)

	}
	if len(departments) > 0 {
		return o.newDepRepo.InsertBranch(tx, departments...)
	}
	return nil
}

func (o *Data) cleanUser(tx *gorm.DB) error {
	allOldUsers := o.oldUserRepo.All(o.DB)

	users := make([]*org.User, 0)
	leaderRelations := make([]org.UserLeaderRelation, 0)
	depRelations := make([]org.UserDepartmentRelation, 0)

	for k := range allOldUsers {
		if allOldUsers[k].LeaderID != "" {
			leaderRelation := org.UserLeaderRelation{}
			leaderRelation.ID = id2.ShortID(0)
			leaderRelation.LeaderID = allOldUsers[k].LeaderID
			leaderRelation.UserID = allOldUsers[k].ID
			leaderRelations = append(leaderRelations, leaderRelation)
		}
		user := &org.User{}
		user.ID = allOldUsers[k].ID
		user.Name = allOldUsers[k].UserName
		user.Phone = allOldUsers[k].Phone
		user.Email = allOldUsers[k].Email
		user.UseStatus = allOldUsers[k].UseStatus
		user.Position = allOldUsers[k].Position
		user.CreatedAt = allOldUsers[k].CreateTime
		user.CreatedBy = allOldUsers[k].CreatBy
		user.Avatar = allOldUsers[k].Avatar
		user.PasswordStatus = allOldUsers[k].PasswordStatus

		users = append(users, user)
	}
	oldDepartmentRelations := o.oldUserDepRepo.All(o.DB)
	for k := range oldDepartmentRelations {
		departmentRelation := org.UserDepartmentRelation{}
		departmentRelation.ID = id2.ShortID(0)
		departmentRelation.DepID = oldDepartmentRelations[k].DepID
		departmentRelation.UserID = oldDepartmentRelations[k].UserID
		depRelations = append(depRelations, departmentRelation)
	}
	var err error = nil
	if len(users) > 0 {
		err = o.newUserRepo.InsertBranch(tx, users...)
	}
	if len(leaderRelations) > 0 {
		err = o.newUserLeaderRepo.InsertBranch(tx, leaderRelations...)
	}
	if len(depRelations) > 0 {
		err = o.newUserDepRepo.InsertBranch(tx, depRelations...)
	}
	return err
}
