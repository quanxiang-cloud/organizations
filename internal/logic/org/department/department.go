package department

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

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/page"
)

// Department interface
type Department interface {
	UserSelectByCondition(c context.Context, r *ViewerSearchListRequest) (*page.Page, error)
	PageList(c context.Context, r *AdminSearchListRequest) (*page.Page, error)
	UserSelectByID(c context.Context, r *SearchOneRequest) (*ViewerDepartmentResponse, error)
	AdminSelectByID(c context.Context, r *SearchOneRequest) (*AdminDepartmentResponse, error)
	UserSelectByPID(c context.Context, r *SearchListByPIDRequest) (*page.Page, error)
	AdminSelectByPID(c context.Context, r *SearchListByPIDRequest) (*page.Page, error)
	Add(c context.Context, r *AddRequest) (res *AddResponse, err error)
	Update(c context.Context, r *UpdateRequest) (*UpdateResponse, error)
	Delete(c context.Context, r *DelOneRequest) (*DelOneResponse, error)
	Tree(c context.Context, r *TreeRequest) (*TreeResponse, error)
	GetDepByIDs(c context.Context, r *GetByIDsRequest) (*GetByIDsResponse, error)
	SetDEPLeader(c context.Context, r *SetDEPLeaderRequest) (*SetDEPLeaderResponse, error)
	CancelDEPLeader(c context.Context, r *CancelDEPLeaderRequest) (*CancelDEPLeaderResponse, error)
	CheckDEPIsExist(c context.Context, r *CheckDEPIsExistRequest) (*CheckDEPIsExistResponse, error)
	GetDepsByIDs(c context.Context, r *GetDepsByIDsRequest) (*GetDepsByIDsResponse, error)
	GetMaxGrade(c context.Context, r *GetMaxGradeRequest) (*GetMaxGradeResponse, error)
	// TestDelete just for test
	TestDelete(c context.Context, r *DelOneRequest) error
}

const (
	allStatus  = 0
	firsGrade  = 1
	exist      = 1
	notExist   = -1
	depAttrDEP = 2
	depAttrCOM = 1
)

// department
type department struct {
	DB          *gorm.DB
	userRepo    org.UserRepo
	depRepo     org.DepartmentRepo
	userDepRepo org.UserDepartmentRelationRepo
	search      *user.Search
}

// NewDepartment new
func NewDepartment(db *gorm.DB) Department {
	return &department{
		depRepo:     mysql2.NewDepartmentRepo(),
		userDepRepo: mysql2.NewUserDepartmentRelationRepo(),
		DB:          db,
		search:      user.GetSearch(),
		userRepo:    mysql2.NewUserRepo(),
	}
}

// TestDelete just for tst
func (d *department) TestDelete(c context.Context, r *DelOneRequest) error {
	relations := d.userDepRepo.SelectByDEPID(d.DB, r.ID)
	if len(relations) == 0 {
		tx := d.DB.Begin()
		err := d.depRepo.Delete(c, tx, r.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
	return error2.New(code.InvalidDELDEP)
}

// CheckDEPIsExistRequest check department exist
type CheckDEPIsExistRequest struct {
	DepID   string `json:"depID" binding:"required,max=64"`
	DepName string `json:"depName" binding:"required,max=60,excludesall=0x2C!@#$?.%:*&^+><=；;"`
}

// CheckDEPIsExistResponse check response
type CheckDEPIsExistResponse struct {
	IsExist int `json:"isExist"`
}

// CheckDEPIsExist check department exist
func (d *department) CheckDEPIsExist(c context.Context, r *CheckDEPIsExistRequest) (*CheckDEPIsExistResponse, error) {
	res := d.depRepo.SelectByPIDAndName(c, d.DB, r.DepID, r.DepName)
	result := CheckDEPIsExistResponse{}
	if res != nil {
		result.IsExist = exist
	} else {
		result.IsExist = notExist
	}
	return &result, nil
}

// SetDEPLeaderRequest set leader request
type SetDEPLeaderRequest struct {
	DepID  string `json:"depID" binding:"required,max=64"`
	UserID string `json:"userID" binding:"required,max=64"`
	Attr   string `json:"attr" binding:"required"`
}

// SetDEPLeaderResponse set leader response
type SetDEPLeaderResponse struct {
}

// SetDEPLeader set leader
func (d *department) SetDEPLeader(c context.Context, r *SetDEPLeaderRequest) (*SetDEPLeaderResponse, error) {
	tx := d.DB.Begin()
	addData := &org.UserDepartmentRelation{}

	relation := d.userDepRepo.SelectByUserIDAndDepID(d.DB, r.UserID, r.DepID)
	if relation == nil {
		addData.ID = id2.ShortID(0)
		addData.DepID = r.DepID
		addData.UserID = r.UserID
		addData.Attr = r.Attr
		err := d.userDepRepo.Add(tx, addData)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return nil, nil
	}
	if relation != nil && relation.Attr == r.Attr {
		return nil, nil
	}
	relation.Attr = r.Attr
	err := d.userDepRepo.Update(tx, relation)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	user := d.userRepo.Get(c, d.DB, r.UserID)
	d.search.PushUser(c, user)
	return nil, nil

}

// CancelDEPLeaderRequest cancel leader request
type CancelDEPLeaderRequest struct {
	DepID  string `json:"depID" binding:"required,max=64"`
	UserID string `json:"userID" binding:"required"`
	Attr   string `json:"attr"`
}

// CancelDEPLeaderResponse cancel leader response
type CancelDEPLeaderResponse struct {
}

// CancelDEPLeader cancel leader request
func (d *department) CancelDEPLeader(c context.Context, r *CancelDEPLeaderRequest) (*CancelDEPLeaderResponse, error) {

	relation := d.userDepRepo.SelectByUserIDAndDepID(d.DB, r.UserID, r.DepID)
	if relation == nil {
		return nil, nil
	}
	tx := d.DB.Begin()
	if relation.Attr != r.Attr {
		relation.Attr = r.Attr
		err := d.userDepRepo.Update(tx, relation)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		user := d.userRepo.Get(c, d.DB, r.UserID)
		d.search.PushUser(c, user)
	}

	return nil, nil
}

// AddRequest ad request
type AddRequest struct {
	Name      string `json:"name" binding:"required,max=60,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	UseStatus int    `json:"useStatus" binding:"-"`
	PID       string `json:"pid" binding:"max=64"`
	CreatBy   string `json:"-" binding:"-"`
	//1:company,2:department
	Attr int `json:"attr"`
}

// AddResponse add response
type AddResponse struct {
	ID string `json:"id"`
}

// Add add
func (d *department) Add(c context.Context, r *AddRequest) (res *AddResponse, err error) {
	if r.PID == "" {
		supper := d.depRepo.SelectSupper(c, d.DB)
		if supper != nil {
			return nil, error2.New(code.TopDepExist)
		}
	}
	one := d.depRepo.SelectByPIDAndName(c, d.DB, r.PID, r.Name)
	if one != nil {
		return nil, error2.New(code.NameUsed)
	}
	tx := d.DB.Begin()
	id := id2.ShortID(0)
	nowUnix := time.NowUnix()

	insertData := org.Department{
		ID:   id,
		Name: r.Name,

		UseStatus: consts.NormalStatus,
		CreatedAt: nowUnix,
		UpdatedAt: nowUnix,
		CreatedBy: r.CreatBy,
	}
	if r.Attr == 0 {
		insertData.Attr = depAttrDEP
	} else {
		insertData.Attr = r.Attr
	}

	if r.PID != "" {
		insertData.PID = r.PID
		p := d.depRepo.Get(c, d.DB, r.PID)
		insertData.SuperPID = p.SuperPID
		insertData.Grade = p.Grade + 1
	} else {
		insertData.SuperPID = id
		insertData.Grade = firsGrade
	}
	err = d.depRepo.Insert(c, tx, &insertData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	d.search.PushDep(c)
	adminDepartment := AddResponse{}
	adminDepartment.ID = id
	return &adminDepartment, nil
}

// UpdateRequest update request
type UpdateRequest struct {
	ID        string `json:"id" binding:"required,max=64"`
	Name      string `json:"name"`
	UseStatus int    `json:"useStatus" binding:"-"`
	//1:company,2:department
	Attr     int    `json:"attr"`
	PID      string `json:"pid" binding:"max=64"`
	UpdateBy string `json:"updateBy"`
}

// UpdateResponse update response
type UpdateResponse struct {
}

// Update update
func (d *department) Update(c context.Context, r *UpdateRequest) (*UpdateResponse, error) {
	upUinx := time.NowUnix()
	dep := d.depRepo.Get(c, d.DB, r.ID)
	if dep != nil {
		if r.PID == "" || r.PID == dep.PID {
			dep.ID = r.ID
			dep.Name = r.Name
			dep.UseStatus = consts.NormalStatus
			dep.UpdatedAt = upUinx
			dep.UpdatedBy = r.UpdateBy
			dep.Attr = r.Attr
		} else {
			if r.PID == r.ID {
				return nil, error2.New(code.InvalidDEPID)
			}
			if d.checkNewPIDIsChildID(c, r.ID, r.PID) {
				return nil, error2.New(code.CanNotMoveParentToChild)
			}
			one := d.depRepo.SelectByPIDAndName(c, d.DB, r.PID, r.Name)
			if one != nil && one.ID != "" {
				return nil, error2.New(code.NameUsed)
			}
			p := d.depRepo.Get(c, d.DB, r.PID)
			dep.ID = r.ID
			dep.Name = r.Name
			dep.UseStatus = r.UseStatus
			dep.UpdatedAt = upUinx
			dep.PID = r.PID
			dep.Attr = r.Attr
			dep.Grade = p.Grade + 1
		}
		tx := d.DB.Begin()
		err := d.depRepo.Update(c, tx, dep)
		err = d.updateChildGrade(c, dep.ID, dep.Grade, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()

		relations := d.userDepRepo.SelectByDEPID(d.DB, r.ID)
		if len(relations) > 0 {
			userIDs := make([]string, 0, len(relations))
			for k := range relations {
				userIDs = append(userIDs, relations[k].UserID)
			}
			users := d.userRepo.List(c, d.DB, userIDs...)
			d.search.PushUser(c, users...)
		}

		d.search.PushDep(c)
		return nil, nil
	}
	return nil, error2.New(code.InvalidUpdate)

}

// ViewerSearchListRequest user search request
type ViewerSearchListRequest struct {
	PID       string `json:"pid" form:"pid"  binding:"max=64"`
	SuperPID  string `json:"superPID" form:"superPID" binding:"max=64"`
	Name      string `json:"name" form:"name" binding:"max=60,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	CompanyID string `json:"companyID" form:"companyID" binding:"max=64"`
	UseStatus int    `json:"useStatus" form:"useStatus" binding:"-"`
	Page      int    `json:"page" form:"page" binding:"-"`
	Limit     int    `json:"limit" form:"limit" binding:"-"`
}

// ViewerDepartmentResponse user search request
type ViewerDepartmentResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	PID      string `json:"pid"`
	SuperPID string `json:"superID,omitempty"`
	Grade    int    `json:"grade,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

// UserSelectByCondition user select by condition
func (d *department) UserSelectByCondition(c context.Context, r *ViewerSearchListRequest) (*page.Page, error) {
	list, total := d.depRepo.PageList(c, d.DB, consts.NormalStatus, r.Page, r.Limit)
	page := page.Page{}
	if len(list) > 0 {
		res := make([]ViewerDepartmentResponse, 0)
		for k := range list {
			re := ViewerDepartmentResponse{}
			re.ID = list[k].ID
			re.Name = list[k].Name
			re.PID = list[k].PID
			re.SuperPID = list[k].SuperPID
			re.Grade = list[k].Grade
			re.Attr = list[k].Attr
			res = append(res, re)
		}
		page.Data = res
		page.TotalCount = total
		return &page, nil
	}
	return &page, nil

}

// AdminSearchListRequest admin get list request
type AdminSearchListRequest struct {
	PID       string `json:"pid" form:"pid"  binding:"max=64"`
	SuperPID  string `json:"superPID" form:"superPID" binding:"max=64"`
	UseStatus int    `json:"useStatus" form:"useStatus" binding:"-"`
	Page      int    `json:"page" form:"page" binding:"-"`
	Limit     int    `json:"limit" form:"limit" binding:"-"`
}

// AdminDepartmentResponse admin get list response
type AdminDepartmentResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid,omitempty"`
	SuperPID  string `json:"superID,omitempty"`
	CompanyID string `json:"companyID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	CreateAt  int64  `json:"createAt,omitempty"`
	UpdateAt  int64  `json:"updateAt,omitempty"`
	CreateBy  string `json:"creatBy,omitempty"`
	UpdateBy  string `json:"updateBy,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

// PageList page list
func (d *department) PageList(c context.Context, r *AdminSearchListRequest) (*page.Page, error) {
	list, total := d.depRepo.PageList(c, d.DB, r.UseStatus, r.Page, r.Limit)
	page := page.Page{}
	if len(list) > 0 {
		res := make([]AdminDepartmentResponse, 0)
		for k := range list {
			re := AdminDepartmentResponse{}
			re.ID = list[k].ID
			re.Name = list[k].Name
			re.PID = list[k].PID
			re.UseStatus = list[k].UseStatus
			re.SuperPID = list[k].SuperPID
			re.Grade = list[k].Grade
			re.CreateAt = list[k].CreatedAt
			re.UpdateAt = list[k].UpdatedAt
			re.CreateBy = list[k].CreatedBy
			re.UpdateBy = list[k].UpdatedBy
			re.Attr = list[k].Attr
			res = append(res, re)
		}
		page.Data = res
		page.TotalCount = total
		return &page, nil
	}
	return &page, nil

}

// SearchOneRequest get one
type SearchOneRequest struct {
	ID string `json:"id" form:"id" binding:"required,max=90"`
}

// UserSelectByID select by id
func (d *department) UserSelectByID(c context.Context, r *SearchOneRequest) (*ViewerDepartmentResponse, error) {
	data := d.depRepo.Get(c, d.DB, r.ID)
	if data != nil {
		res := ViewerDepartmentResponse{}
		res.ID = data.ID
		res.Name = data.Name
		res.PID = data.PID
		res.SuperPID = data.SuperPID
		res.Grade = data.Grade
		res.Attr = data.Attr
		return &res, nil
	}
	return nil, error2.New(code.DataNotExist)

}

// AdminSelectByID amdin select by id
func (d *department) AdminSelectByID(c context.Context, r *SearchOneRequest) (*AdminDepartmentResponse, error) {
	data := d.depRepo.Get(c, d.DB, r.ID)
	if data != nil {
		res := AdminDepartmentResponse{}
		res.ID = data.ID
		res.Name = data.Name
		res.PID = data.PID
		res.UseStatus = data.UseStatus
		res.SuperPID = data.SuperPID
		res.Grade = data.Grade
		res.CreateAt = data.CreatedAt
		res.UpdateAt = data.UpdatedAt
		res.CreateBy = data.CreatedBy
		res.UpdateBy = data.UpdatedBy
		res.Attr = data.Attr
		return &res, nil
	}
	return nil, error2.New(code.DataNotExist)
}

// SearchListByPIDRequest select by pid
type SearchListByPIDRequest struct {
	PID       string `json:"pid" form:"pid" binding:"required,max=64"`
	Name      string `json:"name" form:"name" binding:"max=60,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	UseStatus int    `json:"useStatus" form:"useStatus" binding:"-"`
	Page      int    `json:"page" form:"page" binding:"-"`   //页码
	Limit     int    `json:"limit" form:"limit" binding:"-"` //每页数量
}

// UserSelectByPID user select by pid
func (d *department) UserSelectByPID(c context.Context, r *SearchListByPIDRequest) (*page.Page, error) {
	list, total := d.depRepo.SelectByPID(c, d.DB, r.PID, consts.NormalStatus, r.Page, r.Limit)
	pageRes := page.Page{}
	if len(list) > 0 {
		res := make([]ViewerDepartmentResponse, 0)
		for k := range list {
			re := ViewerDepartmentResponse{}
			re.ID = list[k].ID
			re.Name = list[k].Name
			re.PID = list[k].PID
			re.SuperPID = list[k].SuperPID
			re.Grade = list[k].Grade
			re.Attr = list[k].Attr
			res = append(res, re)
		}
		pageRes.Data = res
		pageRes.TotalCount = total
		return &pageRes, nil
	}
	return &pageRes, nil

}

// AdminSelectByPID  admin select by pid
func (d *department) AdminSelectByPID(c context.Context, r *SearchListByPIDRequest) (*page.Page, error) {
	list, total := d.depRepo.SelectByPID(c, d.DB, r.PID, r.UseStatus, r.Page, r.Limit)
	pageRes := page.Page{}
	if len(list) > 0 {
		res := make([]AdminDepartmentResponse, 0)
		for k := range list {
			re := AdminDepartmentResponse{}
			re.ID = list[k].ID
			re.Name = list[k].Name
			re.PID = list[k].PID
			re.UseStatus = list[k].UseStatus
			re.SuperPID = list[k].SuperPID
			re.Grade = list[k].Grade
			re.CreateAt = list[k].CreatedAt
			re.UpdateAt = list[k].UpdatedAt
			re.CreateBy = list[k].CreatedBy
			re.UpdateBy = list[k].UpdatedBy
			re.Attr = list[k].Attr
			res = append(res, re)
		}
		pageRes.Data = res
		pageRes.TotalCount = total
		return &pageRes, nil
	}
	return &pageRes, nil

}

// DelOneRequest delete by id
type DelOneRequest struct {
	ID       string `json:"id" binding:"required,max=64"`
	DeleteBy string `json:"deleteBy"`
}

// DelOneResponse delete response
type DelOneResponse struct {
}

// Delete delete
func (d *department) Delete(c context.Context, r *DelOneRequest) (*DelOneResponse, error) {
	res := d.depRepo.Get(c, d.DB, r.ID)
	if res != nil && res.PID == "" {
		return nil, error2.New(code.CanNotDel)
	}
	_, total := d.depRepo.SelectByPID(c, d.DB, r.ID, allStatus, 1, 10000)
	if total > 0 {
		return nil, error2.New(code.CanNotDel)
	}
	tx := d.DB.Begin()
	err := d.userDepRepo.DeleteByDepIDs(d.DB, r.ID)

	res.UseStatus = consts.DelStatus
	err = d.depRepo.Update(c, tx, res)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	d.search.PushDep(c)
	return nil, nil
}

// TreeRequest tree request
type TreeRequest struct {
	UserID string `json:"userID"`
}

// TreeResponse  tree data
type TreeResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid,omitempty"`
	SuperPID  string `json:"superID,omitempty"`
	CompanyID string `json:"companyID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	//1:company,2:department
	Attr  int            `json:"attr"`
	Child []TreeResponse `json:"child"`
}

// Tree get department tree
func (d *department) Tree(c context.Context, r *TreeRequest) (*TreeResponse, error) {
	all, _ := d.depRepo.PageList(c, d.DB, consts.NormalStatus, 1, 10000)
	departments := make([]TreeResponse, 0)
	for k := range all {
		res := TreeResponse{}
		res.ID = all[k].ID
		res.Name = all[k].Name
		res.PID = all[k].PID
		res.UseStatus = all[k].UseStatus
		res.SuperPID = all[k].SuperPID
		res.Grade = all[k].Grade
		res.Attr = all[k].Attr
		departments = append(departments, res)
	}
	trees := d.makeTrees(departments)
	return trees, nil
}

/*
makeRoot pid=""
*/
func (d *department) makeTrees(deps []TreeResponse) *TreeResponse {
	var outs *TreeResponse = nil
	var mps = make(map[string][]TreeResponse)
	for k, v := range deps {
		if v.PID == "" {
			outs = &deps[k]
		} else {
			mps[v.PID] = append(mps[v.PID], v)
		}

	}
	if outs != nil {

		d.makeTree(outs, mps)

	}

	return outs
}

// GetByIDsRequest get by ids request
type GetByIDsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// GetByIDsResponse get by ids response
type GetByIDsResponse struct {
	Deps []AdminDepartmentResponse `json:"deps"`
}

//GetDepByIDs get by ids
func (d *department) GetDepByIDs(c context.Context, r *GetByIDsRequest) (*GetByIDsResponse, error) {
	list := d.depRepo.List(c, d.DB, r.IDs...)
	if len(list) > 0 {
		res := &GetByIDsResponse{}
		for k := range list {
			re := AdminDepartmentResponse{}
			re.ID = list[k].ID
			re.Name = list[k].Name
			re.PID = list[k].PID
			re.UseStatus = list[k].UseStatus
			re.SuperPID = list[k].SuperPID
			re.Grade = list[k].Grade
			re.Attr = list[k].Attr
			res.Deps = append(res.Deps, re)
		}
		return res, nil
	}
	return nil, nil

}

/*
makeTree
*/
func (d *department) makeTree(dep *TreeResponse, mps map[string][]TreeResponse) {
	for k := range mps {
		if k == dep.ID {
			dep.Child = append(dep.Child, mps[k]...)
		}
	}
	if len(dep.Child) > 0 {
		for i := 0; i < len(dep.Child); i++ {
			d.makeTree(&dep.Child[i], mps)
		}
	}
}

func (d *department) checkNewPIDIsChildID(c context.Context, oldPID, newPID string) bool {
	list, _ := d.depRepo.SelectByPID(c, d.DB, oldPID, allStatus, 1, 10000)
	if len(list) > 0 {
		var flag = false
		for k := range list {
			if list[k].ID == newPID {
				flag = true
				break
			}
			flag = d.checkNewPIDIsChildID(c, list[k].ID, newPID)
			if flag {
				return flag
			}
			continue
		}
		return flag
	}
	return false
}

func (d *department) updateChildGrade(c context.Context, pid string, grade int, tx *gorm.DB) error {
	list, _ := d.depRepo.SelectByPID(c, d.DB, pid, allStatus, 1, 10000)
	if len(list) > 0 {
		var err error = nil
		for k := range list {
			res := org.Department{}
			res.ID = list[k].ID
			res.Grade = grade + 1
			err = d.depRepo.Update(c, tx, &res)
			if err != nil {
				return err
			}
			err = d.updateChildGrade(c, list[k].ID, res.Grade, tx)
			if err != nil {
				return err
			}
			continue
		}
		return err
	}
	return nil

}

// GetDepsByIDsRequest request
type GetDepsByIDsRequest struct {
	IDs []string `json:"ids"`
}

// GetDepsByIDsResponse response
type GetDepsByIDsResponse struct {
	Deps []ViewerDepartmentResponse `json:"deps"`
}

// GetDepsByIDs get dep by ids
func (d *department) GetDepsByIDs(c context.Context, r *GetDepsByIDsRequest) (*GetDepsByIDsResponse, error) {
	list := d.depRepo.List(c, d.DB, r.IDs...)
	if len(list) == 0 {
		return nil, nil
	}
	resp := &GetDepsByIDsResponse{}
	for k := range list {
		response := ViewerDepartmentResponse{}
		response.ID = list[k].ID
		response.Name = list[k].Name
		response.PID = list[k].PID
		response.SuperPID = list[k].SuperPID
		response.Grade = list[k].Grade
		response.Attr = list[k].Grade
		resp.Deps = append(resp.Deps, response)
	}
	return resp, nil
}

// GetMaxGradeRequest request
type GetMaxGradeRequest struct {
}

// GetMaxGradeResponse response
type GetMaxGradeResponse struct {
	Grade int64 `json:"grade"`
}

// GetMaxGrade get dep by ids
func (d *department) GetMaxGrade(c context.Context, r *GetMaxGradeRequest) (*GetMaxGradeResponse, error) {
	maxGrade := d.depRepo.GetMaxGrade(c, d.DB)

	return &GetMaxGradeResponse{
		maxGrade,
	}, nil
}
