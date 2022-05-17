package user

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
	"errors"

	"gorm.io/gorm"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/es"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
)

var search *Search

// Search  search for es searvice
type Search struct {
	ctx            context.Context
	db             *gorm.DB
	userRepo       org.UserRepo
	userLeaderRepo org.UserLeaderRelationRepo
	userDepRepo    org.UserDepartmentRelationRepo
	depRepo        org.DepartmentRepo
	user           chan *SearchUser
	dep            chan *SearchDepartment
}

// SearchUser push data
type SearchUser struct {
	User []*org.User
	Ctx  context.Context
	Sig  chan int
}

// SearchDepartment push data
type SearchDepartment struct {
	Ctx context.Context
	Sig chan int
}

// NewSearch new
func NewSearch(db *gorm.DB) {
	search = &Search{
		ctx:            context.Background(),
		db:             db,
		userRepo:       mysql2.NewUserRepo(),
		userDepRepo:    mysql2.NewUserDepartmentRelationRepo(),
		userLeaderRepo: mysql2.NewUserLeaderRelationRepo(),
		depRepo:        mysql2.NewDepartmentRepo(),
		user:           make(chan *SearchUser),
		dep:            make(chan *SearchDepartment),
	}
	go search.process(search.ctx)

}

// GetSearch  get search
func GetSearch() *Search {
	if search == nil {
		return nil
	}
	return search
}
func (s *Search) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info("**done**")
			return
		case user := <-s.user:
			s.pushUserToSearch(user)
		case dep := <-s.dep:
			s.pushDepToSearch(dep)
		}
	}
}

// PushUser push data
func (s *Search) PushUser(ctx context.Context, sig chan int, user ...*org.User) {
	if len(user) > 0 {
		u := new(SearchUser)
		if sig != nil {
			u.Sig = sig
		}
		u.User = append(u.User, user...)
		u.Ctx = ctx
		s.user <- u
	}
}

// PushDep push data
func (s *Search) PushDep(ctx context.Context, sig chan int) {
	d := new(SearchDepartment)
	d.Ctx = ctx
	if sig != nil {
		d.Sig = sig
	}

	s.dep <- d

}

func (s *Search) pushDepToSearch(dep *SearchDepartment) {
	list, _ := s.depRepo.PageList(dep.Ctx, s.db, 1, 1, 10000)
	if len(list) > 0 {
		departments := new(es.SearchDepartment)
		departments.Ctx = dep.Ctx
		for k := range list {
			department := v1alpha1.Department{}
			department.ID = list[k].ID
			department.PID = list[k].PID
			department.Name = list[k].Name
			department.TenantID = list[k].TenantID
			departments.Deps = append(departments.Deps, department)
		}
		search := es.GetSearch()
		if dep.Sig != nil {
			departments.Sig = dep.Sig
		}
		search.AddDepartmentSearch(departments)
	}
}

// pushUserToSearch push user info to search server
func (s *Search) pushUserToSearch(user *SearchUser) {
	search := es.GetSearch()

	allDeps, _ := s.depRepo.PageList(user.Ctx, s.db, 1, 1, 10000)
	depMap := make(map[string]*org.Department)
	for k := range allDeps {
		depMap[allDeps[k].ID] = &allDeps[k]
	}
	esData := new(es.SearchUser)
	for _, v := range user.User {

		eu := new(v1alpha1.User)
		esData.Ctx = user.Ctx
		eu.ID = v.ID
		eu.Name = v.Name
		eu.Phone = v.Phone
		eu.Email = v.Email
		eu.CreatedAt = v.CreatedAt
		eu.JobNumber = v.JobNumber
		eu.Avatar = v.Avatar
		eu.TenantID = v.TenantID
		eu.Gender = v.Gender
		eu.Source = v.Source
		eu.SelfEmail = v.SelfEmail
		eu.Position = v.Position
		eu.UseStatus = v.UseStatus
		departmentRelations := s.userDepRepo.SelectByUserIDs(s.db, v.ID)
		//组装部门，从当前到顶层
		for _, v1 := range departmentRelations {
			departments := new(es.SearchDepartment)
			departments.Ctx = user.Ctx
			department := v1alpha1.Department{}
			dep := depMap[v1.DepID]
			if dep != nil {
				department.ID = dep.ID
				department.Name = dep.Name
				department.PID = dep.PID
				departments.Deps = append(departments.Deps, department)

				depss := s.getDepToTop(dep.PID, departments.Deps, depMap)
				eu.Departments = append(eu.Departments, depss)

			}
		}
		//寻找leader，从当前到顶层
		leaderToTop, err := s.getLeaderToTop(user.Ctx, v.ID, v.ID)
		if err == nil && leaderToTop != nil {
			eu.Leaders = append(eu.Leaders, leaderToTop...)
		}
		esData.User = append(esData.User, *eu)

	}
	if user.Sig != nil {
		esData.Sig = user.Sig
	}
	search.AddUserSearch(esData)

}

func (s *Search) getDepToTop(depPID string, deps []v1alpha1.Department, depMap map[string]*org.Department) []v1alpha1.Department {
	dep := depMap[depPID]
	if dep != nil {
		department := v1alpha1.Department{}
		department.ID = dep.ID
		department.Name = dep.Name
		department.PID = dep.PID
		deps = append(deps, department)
		if dep.PID != "" {
			return s.getDepToTop(dep.PID, deps, depMap)
		}

	}
	return deps
}

func (s *Search) getLeaderToTop(ctx context.Context, userID, startUserID string) ([][]v1alpha1.Leader, error) {
	relations := s.userLeaderRepo.SelectByUserIDs(s.db, userID)
	if len(relations) > 0 {
		res := make([][]v1alpha1.Leader, 0)
		for k := range relations {
			if relations[k].LeaderID == startUserID {
				return nil, errors.New("circle leader")
			}
			if relations[k].LeaderID != "" {
				ls := make([]v1alpha1.Leader, 0)
				get := s.userRepo.Get(ctx, s.db, relations[k].LeaderID)
				if get != nil {
					leader := v1alpha1.Leader{}
					leader.ID = get.ID
					leader.Name = get.Name
					leader.Attr = relations[k].Attr
					ls = append(ls, leader)
					array, err := s.getLeaderToTop(ctx, get.ID, startUserID)
					if err != nil {
						return nil, err
					}
					if array != nil {
						for k1 := range array {
							ll := append(ls, array[k1]...)
							res = append(res, ll)
						}
						continue
					}
					res = append(res, ls)
				}
			}

		}
		return res, nil

	}
	return nil, nil

}
