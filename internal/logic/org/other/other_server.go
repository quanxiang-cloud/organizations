package other

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
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/encode2"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	ldap "github.com/quanxiang-cloud/organizations/pkg/ladp"
	"github.com/quanxiang-cloud/organizations/pkg/random2"
	"github.com/quanxiang-cloud/organizations/pkg/systems"
	"github.com/quanxiang-cloud/organizations/pkg/verification"
)

// OthServer interface
type OthServer interface {
	AddUsers(c context.Context, r *AddUsersRequest) (res *AddListResponse, err error)
	AddDepartments(c context.Context, r *AddDepartmentRequest) (res *AddListResponse, err error)
	GetUserByIDs(c context.Context, r *GetUserByIDsRequest) (res *GetUserByIDsResponse, err error)
	GetOneUser(c context.Context, r *GetOneRequest) (res *GetOneResponse, err error)
	GetAllUsers(c context.Context, r *UserAllRequest) (res *UserAllResp, err error)
	GetAllDeps(c context.Context, r *DepAllRequest) (res *DepAllDepsResp, err error)
	OtherGetUsersByDepID(c context.Context, r *GetUsersByDepIDRequest) (res *GetUsersByDepIDResponse, err error)
	PushUserToSearch(c context.Context, sig, total chan int)
	PushDepToSearch(c context.Context, sig chan int)
}

// othersServer
type othersServer struct {
	DB          *gorm.DB
	userRepo    org.UserRepo
	userDepRepo org.UserDepartmentRelationRepo
	depRepo     org.DepartmentRepo
	accountReo  org.AccountRepo
	//message     message.Message
	redisClient    redis.UniversalClient
	columnRepo     org.UserTableColumnsRepo
	ldap           ldap.Ldap
	conf           configs.Config
	userLeaderRepo org.UserLeaderRelationRepo
	search         *user.Search
}

// NewOtherServer 实例
func NewOtherServer(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) OthServer {
	return &othersServer{
		userRepo:    mysql2.NewUserRepo(),
		userDepRepo: mysql2.NewUserDepartmentRelationRepo(),
		depRepo:     mysql2.NewDepartmentRepo(),
		accountReo:  mysql2.NewAccountRepo(),
		DB:          db,
		//message:     message.NewMessage(configs.Config.InternalNet),
		redisClient:    redisClient,
		columnRepo:     mysql2.NewUserTableColumnsRepo(),
		ldap:           ldap.NewLdap(conf.InternalNet),
		conf:           conf,
		userLeaderRepo: mysql2.NewUserLeaderRelationRepo(),
		search:         user.GetSearch(),
	}
}

// AddUsersRequest other server add user request
type AddUsersRequest struct {
	Users []AddUser `json:"users"`
	//1:sync department data,-1:no action
	SyncDEP int `json:"syncDep"`
	//1:update old data,-1:no action
	IsUpdate   int    `json:"isUpdate"`
	SyncID     string `json:"syncID"`
	SyncSource string `json:"syncSource"`
	Profile    header2.Profile
}

// AddUser other server add user to org
type AddUser struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	AccountID string `json:"-"`
	SelfEmail string `json:"selfEmail,omitempty"`
	IDCard    string `json:"idCard,omitempty"`
	Address   string `json:"address,omitempty"`
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int `json:"useStatus,omitempty"`
	//0:null,1:man,2:woman
	Gender    int      `json:"gender,omitempty"`
	CompanyID string   `json:"companyID,omitempty"`
	Position  string   `json:"position,omitempty"`
	Avatar    string   `json:"avatar,omitempty"`
	Remark    string   `json:"remark,omitempty"`
	JobNumber string   `json:"jobNumber,omitempty"`
	DepsID    []string `json:"depsID,omitempty"`
	LeadersID []string `json:"leadersID,omitempty"`
	EntryTime int64    `json:"entryTime,omitempty" `
	Source    string   `json:"source,omitempty" `
	SourceID  string   `json:"sourceID,omitempty" `
}

// AddListResponse other server add user or dep to org response
type AddListResponse struct {
	Result map[int]*Result `json:"result"`
}

// Result list add response
type Result struct {
	ID     string `json:"id"`
	Remark string `json:"remark"`
	Attr   int    `json:"attr"` //11 add ok,0fail,12, update ok
}

// AddUsers  other server add user
func (u *othersServer) AddUsers(c context.Context, r *AddUsersRequest) (res *AddListResponse, err error) {
	result, err := u.addUserOrUpdate(c, r.Users, r.IsUpdate, r.Profile)
	if err != nil {
		return nil, err
	}
	res = &AddListResponse{}
	res.Result = result
	return res, nil
}

func (u *othersServer) PushUserToSearch(c context.Context, sig, total chan int) {
	var index = 1
	var size = 300
	for {

		list, _ := u.userRepo.PageList(c, u.DB, 0, index, size, nil)
		if len(list) > 0 {
			u.search.PushUser(c, sig, list...)
			index = index + 1
			continue
		}
		if total != nil {

			total <- index
		}
		break
	}

}

// AddDepartmentRequest other server add  department to org request
type AddDepartmentRequest struct {
	Deps []AddDep `json:"deps"`
	//1:sync department data,-1:no action
	SyncDEP int `json:"syncDep"`
	//1:update old data,-1:no action
	IsUpdate   int    `json:"isUpdate"`
	SyncID     string `json:"syncID"`
	SyncSource string `json:"syncSource"`
	Profile    header2.Profile
}

// AddDep other server add department to org
type AddDep struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	//1:normal,-1:del,-2:invalid
	UseStatus int `json:"useStatus"`
	//1:company,2:department
	Attr      int    `json:"attr"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID"`
	Grade     int    `json:"grade"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	Remark    string `json:"remark,omitempty"`
	LeaderID  string `json:"leaderID,omitempty"`
}

// AddDepartments  other server add department
func (u *othersServer) AddDepartments(c context.Context, r *AddDepartmentRequest) (res *AddListResponse, err error) {

	result, err := u.addDEP(c, r.Deps, r.IsUpdate, r.Profile)

	if err != nil {
		return nil, err
	}
	res = &AddListResponse{}
	res.Result = result
	return res, nil
}

func (u *othersServer) PushDepToSearch(c context.Context, sig chan int) {
	u.search.PushDep(c, sig)
}

//GetUserByIDsRequest get user by ids request
type GetUserByIDsRequest struct {
	IDs     []string `json:"ids"`
	Profile header2.Profile
}

// GetUserByIDsResponse get user by ids response
type GetUserByIDsResponse struct {
	Users []user.ViewerSearchOneUserResponse `json:"users"`
}

//GetUserByIDs get user info by ids
func (u *othersServer) GetUserByIDs(c context.Context, rq *GetUserByIDsRequest) (*GetUserByIDsResponse, error) {
	if len(rq.IDs) > 0 {
		res := &GetUserByIDsResponse{}
		users := GetUserLeader(c, u.userRepo, u.userLeaderRepo, u.DB, rq.IDs...)
		res.Users = users

		relations := u.userDepRepo.SelectByUserIDs(u.DB, rq.IDs...)
		userDep := make(map[string][]string)
		for k := range relations {
			userDep[relations[k].UserID] = append(userDep[relations[k].UserID], relations[k].DepID)
		}
		depList, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
		depMap := make(map[string]*org.Department)
		for k := range depList {
			depMap[depList[k].ID] = &depList[k]
		}

		for k := range res.Users {

			if v, ok := userDep[res.Users[k].ID]; ok {
				for k1 := range v {
					responses := make([]user.DepOneResponse, 0)
					res.Users[k].Dep = append(res.Users[k].Dep, user.FindDepToTop(depMap, v[k1], responses))
				}

			}
		}

		return res, nil
	}
	return nil, error2.New(code.InvalidParams)
}

const (
	addOk    = 11
	fail     = 0
	updateOk = 12
	standBy  = -1
)

func (u *othersServer) addUserOrUpdate(c context.Context, reqData []AddUser, isUpdate int, profile header2.Profile) (result map[int]*Result, err error) {
	result = make(map[int]*Result)
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)
A:
	for k := range reqData {
		switch reqData[k].UseStatus {
		case consts.NormalStatus:
			reqData[k].UseStatus = consts.NormalStatus
		case consts.ActiveStatus:
			reqData[k].UseStatus = consts.ActiveStatus
		case consts.UnWork:
			reqData[k].UseStatus = consts.UnWork
		case consts.DelStatus:
			reqData[k].UseStatus = consts.DelStatus
		default:
			reqData[k].UseStatus = consts.UnNormalStatus
		}
		if reqData[k].Phone != "" {
			if !verification.CheckPhone(reqData[k].Phone) {
				reqData[k].Remark = code.CodeTable[code.InvalidPhone]
				result[k] = &Result{
					Attr:   fail,
					Remark: code.CodeTable[code.InvalidPhone],
				}
				continue
			}
		}
		if reqData[k].Email != "" {
			if !verification.CheckEmail(reqData[k].Email) {
				reqData[k].Remark = code.CodeTable[code.InvalidEmail]
				result[k] = &Result{
					Attr:   fail,
					Remark: code.CodeTable[code.InvalidEmail],
				}
				continue
			}
		} else {
			reqData[k].Remark = code.CodeTable[code.EmailRequired]
			result[k] = &Result{
				Attr:   fail,
				Remark: code.CodeTable[code.EmailRequired],
			}
			continue
		}
		if reqData[k].SelfEmail != "" {
			if !verification.CheckEmail(reqData[k].SelfEmail) {
				reqData[k].Remark = code.CodeTable[code.InvalidEmail]
				result[k] = &Result{
					Attr:   fail,
					Remark: code.CodeTable[code.EmailRequired],
				}
				continue
			}
		}

		nowUnix := time.NowUnix()
		u2 := &org.User{
			Name:      reqData[k].Name,
			Phone:     reqData[k].Phone,
			Email:     reqData[k].Email,
			SelfEmail: reqData[k].SelfEmail,
			IDCard:    reqData[k].IDCard,
			Address:   reqData[k].Address,
			Position:  reqData[k].Position,
			Avatar:    reqData[k].Avatar,
			JobNumber: reqData[k].JobNumber,
			Gender:    reqData[k].Gender,
			Source:    reqData[k].Source,
			UpdatedAt: nowUnix,
			UpdatedBy: profile.UserID,
			UseStatus: reqData[k].UseStatus,
		}

		if u2.UseStatus == consts.ActiveStatus {
			u2.PasswordStatus = consts.ResetPasswordStatus
		} else {
			u2.PasswordStatus = consts.NormalStatus
		}
		oldAccount := u.accountReo.SelectByAccount(u.DB, reqData[k].Email)
		tx := u.DB.Begin()
		var userID = ""
		result[k] = &Result{}
		if oldAccount == nil {
			userID = id2.HexUUID(true)
			reqData[k].ID = userID
			u2.ID = userID
			u2.CreatedBy = profile.UserID
			u2.CreatedAt = nowUnix

			account := &org.Account{}
			account.ID = id2.ShortID(0)
			account.UserID = userID
			account.Account = reqData[k].Email
			pwd := random2.RandomString(int(info.PwdMinLen), info.PwdType)
			if u.conf.Model == "debug" {
				pwd = "654321a.."
				account.Password = encode2.MD5Encode(pwd)
			} else {
				account.Password = encode2.MD5Encode(pwd)
			}
			account.CreatedBy = profile.UserID
			account.CreatedAt = nowUnix
			account.UpdatedAt = nowUnix
			err := u.insertUser(c, tx, u2, account)
			if err != nil {
				tx.Rollback()
				result[k].Attr = fail
				return nil, err
			}
			result[k].Attr = addOk
			result[k].ID = userID
		} else {
			userID = oldAccount.UserID
			old := u.userRepo.Get(c, u.DB, userID)
			if old != nil {
				if isUpdate == 1 {
					result[k].ID = userID
					oldAccount.Account = reqData[k].Email
					oldAccount.UpdatedBy = profile.UserID
					oldAccount.UpdatedAt = nowUnix

					u2.ID = userID

					err := u.updateUser(c, tx, u2, oldAccount)
					if err != nil {
						tx.Rollback()
						result[k].Attr = fail
						continue
					}
					result[k].Attr = updateOk
				}

			}
		}
		for k1 := range reqData[k].LeadersID {
			err := user.CheckLeader(c, u.DB, u.userLeaderRepo, userID, reqData[k].LeadersID[k1])
			if err != nil {
				tx.Rollback()
				result[k].Attr = fail
				continue A
			}
		}
		err = u.dealUserDepartmentRelation(c, tx, userID, reqData[k].DepsID...)
		if err != nil {
			tx.Rollback()
			result[k] = &Result{
				Attr: fail,
			}
			return nil, err
		}
		err = u.dealUserLeaderRelation(c, tx, userID, reqData[k].LeadersID...)
		if err != nil {
			tx.Rollback()
			result[k] = &Result{
				Attr: fail,
			}
			return nil, err
		}
		tx.Commit()
	}
	return result, nil
}

func (u *othersServer) insertUser(c context.Context, tx *gorm.DB, u2 *org.User, account *org.Account) error {
	err := u.userRepo.Insert(c, tx, u2)
	err = u.accountReo.Insert(tx, account)
	return err
}

func (u *othersServer) updateUser(c context.Context, tx *gorm.DB, u2 *org.User, account *org.Account) error {
	err := u.userRepo.UpdateByID(c, tx, u2)
	err = u.accountReo.Update(tx, account)
	return err
}

func (u *othersServer) dealUserDepartmentRelation(c context.Context, tx *gorm.DB, userID string, depsID ...string) error {
	if len(depsID) > 0 {
		err := u.userDepRepo.DeleteByUserIDs(tx, userID)
		if err != nil {
			return err
		}
		for _, v := range depsID {
			relation := org.UserDepartmentRelation{
				ID:     id2.ShortID(0),
				UserID: userID,
				DepID:  v,
			}
			err := u.userDepRepo.Add(tx, &relation)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (u *othersServer) dealUserLeaderRelation(c context.Context, tx *gorm.DB, userID string, leadersID ...string) error {
	if len(leadersID) > 0 {
		err := u.userLeaderRepo.DeleteByUserIDs(tx, userID)
		if err != nil {
			return err
		}
		for _, v := range leadersID {
			if v == "" || v == "0" {
				continue
			}
			relation := org.UserLeaderRelation{
				ID:       id2.ShortID(0),
				UserID:   userID,
				LeaderID: v,
			}
			err := u.userLeaderRepo.Add(tx, &relation)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (u *othersServer) addDEP(c context.Context, reqData []AddDep, isUpdate int, profile header2.Profile) (map[int]*Result, error) {
	supper := u.depRepo.SelectSupper(c, u.DB)
	var supperID = ""
	if supper != nil {
		supperID = supper.ID
		for k := range reqData {
			if reqData[k].PID == "" && reqData[k].ID != supperID {
				reqData[k].PID = supperID
				break
			}
		}
	} else {
		for k := range reqData {
			if reqData[k].PID == "" {
				supperID = reqData[k].ID
				break
			}
		}
	}

	oldDeps, _ := u.depRepo.PageList(c, u.DB, 0, 1, 100000)

	res, err := u.insertOrUpdateDep(c, oldDeps, reqData, supperID, profile)
	if err != nil {
		return nil, err
	}
	all, _ := u.depRepo.PageList(c, u.DB, 1, 1, 100000)
	makeDepGrade(supperID, all, consts.FirsGrade)
	for k := range all {
		u.depRepo.Update(c, u.DB, &all[k])
	}
	return res, nil
}

func (u *othersServer) insertOrUpdateDep(c context.Context, oldDeps []org.Department, reqData []AddDep, supperID string, profile header2.Profile) (map[int]*Result, error) {
	unix := time.NowUnix()
	result := make(map[int]*Result)
	if len(oldDeps) > 0 {
		depMap := make(map[string]*org.Department)
		for k := range oldDeps {
			depMap[oldDeps[k].ID] = &oldDeps[k]
		}
		for k := range reqData {
			if v, ok := depMap[reqData[k].ID]; ok && v != nil {
				v.Name = reqData[k].Name
				v.PID = reqData[k].PID
				v.SuperPID = supperID
				v.UseStatus = reqData[k].UseStatus
				v.UpdatedAt = unix
				v.Grade = reqData[k].Grade
				err := u.depRepo.Update(c, u.DB, v)
				if err != nil {
					result[k] = &Result{
						ID:   reqData[k].ID,
						Attr: fail,
					}
					continue
				}
				result[k] = &Result{
					ID:   reqData[k].ID,
					Attr: updateOk,
				}
			} else {
				d := &org.Department{}
				d.ID = id2.ShortID(0)
				d.PID = reqData[k].PID
				d.SuperPID = supperID
				d.Name = reqData[k].Name
				d.UseStatus = reqData[k].UseStatus
				d.CreatedBy = profile.UserID
				d.CreatedAt = unix
				d.UpdatedAt = unix
				d.Grade = reqData[k].Grade
				err := u.depRepo.Insert(c, u.DB, d)
				if err != nil {
					result[k] = &Result{
						Attr: fail,
					}
					continue
				}
				result[k] = &Result{
					ID:   d.ID,
					Attr: addOk,
				}
			}
		}
		return result, nil
	}
	for k := range reqData {
		d := &org.Department{}
		d.ID = id2.ShortID(0)
		d.PID = reqData[k].PID
		d.SuperPID = supperID
		d.Name = reqData[k].Name
		d.UseStatus = reqData[k].UseStatus
		d.CreatedBy = profile.UserID
		d.CreatedAt = unix
		d.UpdatedAt = unix
		d.Grade = reqData[k].Grade
		err := u.depRepo.Insert(c, u.DB, d)
		if err != nil {
			result[k] = &Result{
				Attr: fail,
			}
			continue
		}
		result[k] = &Result{
			ID:   d.ID,
			Attr: addOk,
		}
	}
	return result, nil

}

func makeDepGrade(pid string, list []org.Department, grade int) {
	for k := range list {
		if list[k].PID == pid {
			list[k].Grade = grade + 1
			makeDepGrade(list[k].ID, list, list[k].Grade)
		}

	}
}

// GetOneRequest get one request
type GetOneRequest struct {
	ID      string `json:"id"`
	Profile header2.Profile
}

// GetOneResponse get one response
type GetOneResponse struct {
	User user.ViewerSearchOneUserResponse `json:"user"`
}

// GetOneUser get one user info
func (u *othersServer) GetOneUser(c context.Context, r *GetOneRequest) (*GetOneResponse, error) {

	userBtye := u.redisClient.Get(c, consts.RedisTokenUserInfo+r.ID).Val()
	res := GetOneResponse{}
	resUser := user.ViewerSearchOneUserResponse{}
	if userBtye != "" {
		err := json.Unmarshal([]byte(userBtye), &resUser)
		if err != nil {
			return nil, err
		}
		res.User = resUser
		return &res, nil
	}
	userOne := u.userRepo.Get(c, u.DB, r.ID)
	if userOne != nil {

		resUser.ID = userOne.ID
		resUser.Name = userOne.Name
		resUser.Email = userOne.Email
		resUser.Status = userOne.UseStatus

		relations := u.userDepRepo.SelectByUserIDs(u.DB, r.ID)
		if len(relations) > 0 {
			depIDs := make([]string, 0, len(relations))
			for _, v := range relations {
				depIDs = append(depIDs, v.DepID)
			}
			departments := u.depRepo.List(c, u.DB, depIDs...)
			depList, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
			depMap := make(map[string]*org.Department)
			for k := range depList {
				depMap[depList[k].ID] = &depList[k]
			}

			if len(departments) > 0 {
				for _, v := range departments {
					responses := make([]user.DepOneResponse, 0)
					resUser.Dep = append(resUser.Dep, user.FindDepToTop(depMap, v.ID, responses))
				}

			}
		}

		return &res, nil
	}
	return nil, nil
}

// DepAllRequest get department request
type DepAllRequest struct {
}

// DepAllDepsResp get department response
type DepAllDepsResp struct {
	Deps []user.DepOneResponse `json:"deps"`
}

// GetAllDeps  get all department
func (u *othersServer) GetAllDeps(c context.Context, r *DepAllRequest) (res *DepAllDepsResp, err error) {

	list, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 1000)
	resData := &DepAllDepsResp{}
	for k := range list {
		add := user.DepOneResponse{}
		add.Name = list[k].Name
		add.ID = list[k].ID
		add.PID = list[k].PID
		add.SuperPID = list[k].SuperPID
		resData.Deps = append(resData.Deps, add)
	}
	return resData, nil
}

// UserAllRequest get all user request
type UserAllRequest struct {
	Profile header2.Profile
}

// UserAllResp get all user response
type UserAllResp struct {
	All []user.ViewerSearchOneUserResponse `json:"all"`
}

// GetAllUsers get all user
func (u *othersServer) GetAllUsers(c context.Context, r *UserAllRequest) (res *UserAllResp, err error) {

	userList, _ := u.userRepo.PageList(c, u.DB, 0, 1, 10000, nil)
	if len(userList) > 0 {
		userIDs := make([]string, 0, len(userList))
		for k := range userList {
			userIDs = append(userIDs, userList[k].ID)
		}
		relations := u.userDepRepo.SelectByUserIDs(u.DB, userIDs...)
		userDep := make(map[string][]string)
		for k := range relations {
			userDep[relations[k].UserID] = append(userDep[relations[k].UserID], relations[k].DepID)
		}
		departments, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 1000)
		depMap := make(map[string]*org.Department)
		for k := range departments {
			depMap[departments[k].ID] = &departments[k]
		}
		allUser := new(UserAllResp)
		for k := range userList {
			oneUser := user.ViewerSearchOneUserResponse{}
			oneUser.ID = userList[k].ID
			oneUser.Name = userList[k].Name
			oneUser.Email = userList[k].Email
			oneUser.Status = userList[k].UseStatus
			oneUser.JobNumber = userList[k].JobNumber
			if v, ok := userDep[oneUser.ID]; ok {
				for k1 := range v {
					department := depMap[v[k1]]
					if department != nil {
						departmentResp := user.DepOneResponse{
							ID:        department.ID,
							Name:      department.Name,
							UseStatus: department.UseStatus,
							PID:       department.PID,
							SuperPID:  department.SuperPID,
							Grade:     department.Grade,
							Attr:      department.Attr,
						}
						oneUser.Dep = append(oneUser.Dep, []user.DepOneResponse{departmentResp})
					}
				}

			}
			allUser.All = append(allUser.All, oneUser)
		}
		return allUser, nil
	}
	return nil, nil

}

// GetUsersByDepIDRequest get users by ids request
type GetUsersByDepIDRequest struct {
	DepID string `json:"depID"`
	//1:include
	IsIncludeChild int `json:"isIncludeChild"`
}

// GetUsersByDepIDResponse get users by ids response
type GetUsersByDepIDResponse struct {
	Users []user.ViewerSearchOneUserResponse `json:"users"`
}

const includeChildDep = 1

// OtherGetUsersByDepID get users by dep id
func (u *othersServer) OtherGetUsersByDepID(c context.Context, r *GetUsersByDepIDRequest) (res *GetUsersByDepIDResponse, err error) {
	relation := u.userDepRepo.SelectByDEPID(u.DB, r.DepID)
	if len(relation) == 0 && r.IsIncludeChild != includeChildDep {
		return nil, nil
	}
	userIDs := make([]string, 0)
	for k := range relation {
		userIDs = append(userIDs, relation[k].UserID)
	}
	if r.IsIncludeChild == includeChildDep {
		list, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
		depMap := make(map[string][]org.Department)
		for k := range list {
			depMap[list[k].PID] = append(depMap[list[k].PID], list[k])
		}
		depIDS := u.getChildDep(c, r.DepID, depMap)
		if len(depIDS) > 0 {
			relations := u.userDepRepo.SelectByDEPID(u.DB, depIDS...)
			for k := range relations {
				userIDs = append(userIDs, relations[k].UserID)
			}
		}
	}

	response := &GetUsersByDepIDResponse{}
	users := GetUserLeader(c, u.userRepo, u.userLeaderRepo, u.DB, userIDs...)
	response.Users = users
	return response, nil
}

func (u *othersServer) getChildDep(c context.Context, pid string, depMap map[string][]org.Department) []string {
	depIDs := make([]string, 0)

	for k := range depMap[pid] {
		depIDs = append(depIDs, depMap[pid][k].ID)
		depIDs = append(depIDs, u.getChildDep(c, depMap[pid][k].ID, depMap)...)
	}

	if len(depIDs) > 0 {
		return depIDs
	}
	return nil
}

// GetUserLeader get user leader
func GetUserLeader(c context.Context, userRepo org.UserRepo, userLeaderRepo org.UserLeaderRelationRepo, db *gorm.DB, userIDs ...string) []user.ViewerSearchOneUserResponse {
	users := userRepo.List(c, db, userIDs...)
	leaderRelations := userLeaderRepo.SelectByUserIDs(db, userIDs...)
	ud := make(map[string][]string)
	leaderIDs := make([]string, 0)
	for k := range leaderRelations {
		ud[leaderRelations[k].UserID] = append(ud[leaderRelations[k].UserID], leaderRelations[k].LeaderID)
		leaderIDs = append(leaderIDs, leaderRelations[k].LeaderID)
	}
	leaderMap := make(map[string]*org.User)
	leaderList := userRepo.List(c, db, leaderIDs...)
	for k := range leaderList {
		leaderMap[leaderList[k].ID] = leaderList[k]
	}
	responses := make([]user.ViewerSearchOneUserResponse, 0)
	for k := range users {
		resp := user.ViewerSearchOneUserResponse{}
		resp.ID = users[k].ID
		resp.Name = users[k].Name
		resp.Email = users[k].Email
		resp.JobNumber = users[k].JobNumber
		resp.UseStatus = users[k].UseStatus
		resp.Position = users[k].Position

		for k1 := range ud[users[k].ID] {
			leaders := make([]user.Leader, 0)
			leader := user.Leader{}
			if v, ok := leaderMap[ud[users[k].ID][k1]]; ok && v != nil {
				leader.ID = leaderMap[ud[users[k].ID][k1]].ID
				leader.Name = leaderMap[ud[users[k].ID][k1]].Name
				leader.Email = leaderMap[ud[users[k].ID][k1]].Email
				leader.Phone = leaderMap[ud[users[k].ID][k1]].Phone
				leader.UseStatus = leaderMap[ud[users[k].ID][k1]].UseStatus
				leader.Position = leaderMap[ud[users[k].ID][k1]].Position
				leaders = append(leaders, leader)
				resp.Leader = append(resp.Leader, leaders)
			}

		}
		responses = append(responses, resp)
	}
	return responses
}
