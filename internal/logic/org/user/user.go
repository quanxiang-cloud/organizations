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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/department"
	"github.com/quanxiang-cloud/organizations/pkg/goalie"
	"github.com/tealeg/xlsx"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/encode2"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
	ldap "github.com/quanxiang-cloud/organizations/pkg/ladp"
	"github.com/quanxiang-cloud/organizations/pkg/landlord"
	"github.com/quanxiang-cloud/organizations/pkg/message"
	"github.com/quanxiang-cloud/organizations/pkg/page"
	"github.com/quanxiang-cloud/organizations/pkg/random2"
	"github.com/quanxiang-cloud/organizations/pkg/systems"
	"github.com/quanxiang-cloud/organizations/pkg/verification"
)

// User interface
type User interface {
	Add(c context.Context, r *AddUserRequest) (res *AddUserResponse, err error)
	Update(c context.Context, r *UpdateUserRequest) (*UpdateUserResponse, error)
	UpdateAvatar(c context.Context, r *UpdateUserAvatarRequest) (*UpdateUserAvatarResponse, error)
	PageList(c context.Context, r *SearchListUserRequest) (*page.Page, error)
	AdminSelectByID(c context.Context, r *SearchOneUserRequest) (*SearchOneUserResponse, error)
	UserSelectByID(c context.Context, r *ViewerSearchOneUserRequest) (*ViewerSearchOneUserResponse, error)
	UpdateUserStatus(c context.Context, r *StatusRequest) (*StatusResponse, error)
	UpdateUsersStatus(c context.Context, r *ListStatusRequest) (*ListStatusResponse, error)
	AdminChangeUsersDEP(c context.Context, r *ChangeUsersDEPRequest) (*ChangeUsersDEPResponse, error)
	OthGetOneUser(c context.Context, r *TokenUserRequest) (*TokenUserResponse, error)
	IndexCount(c context.Context, r *IndexCountRequest) (*IndexCountResponse, error)
	Register(c context.Context, r *RegisterRequest) (*RegisterResponse, error)
	GetUsersByIDs(c context.Context, r *GetUsersByIDsRequest) (*GetUsersByIDsResponse, error)
	ImportFile(c context.Context, file []byte, profile header2.Profile, r *ImportFileRequest) (*ImportFileResponse, error)
	Template(c context.Context, r *GetTemplateFileRequest) (*GetTemplateFileResponse, error)
}

type user struct {
	DB             *gorm.DB
	userRepo       org.UserRepo
	userDepRepo    org.UserDepartmentRelationRepo
	userLeaderRepo org.UserLeaderRelationRepo
	depRepo        org.DepartmentRepo
	accountReo     org.AccountRepo
	redisClient    redis.UniversalClient
	message        message.Message
	ldap           ldap.Ldap
	conf           configs.Config
	columnRepo     org.UserTableColumnsRepo
	userTenantRepo org.UserTenantRelationRepo
	landlord       landlord.Landlord
	goalie         goalie.Goalie
}

// NewUser new
func NewUser(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) User {

	return &user{
		userRepo:       mysql2.NewUserRepo(),
		userDepRepo:    mysql2.NewUserDepartmentRelationRepo(),
		userLeaderRepo: mysql2.NewUserLeaderRelationRepo(),
		depRepo:        mysql2.NewDepartmentRepo(),
		accountReo:     mysql2.NewAccountRepo(),
		DB:             db,
		redisClient:    redisClient,

		message: message.NewMessage(conf.InternalNet),
		ldap:    ldap.NewLdap(conf.InternalNet),

		conf:           conf,
		columnRepo:     mysql2.NewUserTableColumnsRepo(),
		userTenantRepo: mysql2.NewUserTenantRelationRepo(),
		landlord:       landlord.NewLandlord(conf.InternalNet),
		goalie:         goalie.NewGoalie(conf.InternalNet),
	}
}

// AddUserRequest add user request
type AddUserRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	//0:null,1:man,2:woman
	Gender      int             `json:"gender,omitempty" `
	Source      string          `json:"source,omitempty" `
	CreateBy    string          `json:"createdBy,omitempty" `
	Dep         []DepRequest    `json:"dep"`
	Leader      []LeaderRequest `json:"leader" `
	SendMessage SendMessage     `json:"sendMessage"`
	Header      http.Header
	Profile     header2.Profile
	Password    string
}

// LeaderRequest leader struct
type LeaderRequest struct {
	UserID string `json:"userID"`
	Attr   string `json:"attr"`
}

// DepRequest department struct
type DepRequest struct {
	DepID string `json:"depID"`
	Attr  string `json:"attr"`
}

// send message channel
const (
	NO int = iota
	SENDEMAIL
	SENDPHONE
)

// SendChannelMap send message channel map
type SendChannelMap map[int]string

// SendChannel some channel
var SendChannel = SendChannelMap{
	SENDEMAIL: "email",
	SENDPHONE: "phone",
}

// SendMessage send message
type SendMessage struct {
	UserID      string `json:"userID"`
	SendChannel int    `json:"sendChannel"`
	SendTo      string `json:"sendTo"`
}

// AddUserResponse admin user response
type AddUserResponse struct {
	ID       string      `json:"id"`
	Password string      `json:"password,omitempty"`
	Users    []*org.User `json:"-"`
}

// Add  add user
func (u *user) Add(c context.Context, r *AddUserRequest) (res *AddUserResponse, err error) {

	id := id2.HexUUID(true)
	nowUnix := time2.NowUnix()
	if r.Phone != "" {
		if !verification.CheckPhone(r.Phone) {
			return nil, error2.New(code.InvalidPhone)
		}
	}
	if r.Email != "" {
		if !verification.CheckEmail(r.Email) {
			return nil, error2.New(code.InvalidEmail)
		}
	} else {
		return nil, error2.New(code.EmailRequired)
	}
	if r.SelfEmail != "" {
		if !verification.CheckEmail(r.SelfEmail) {
			return nil, error2.New(code.InvalidEmail)
		}
	}
	for _, v := range r.Leader {
		err := CheckLeader(c, u.DB, u.userLeaderRepo, v.UserID, r.ID)
		if err != nil {
			return nil, error2.New(code.ErrCircleData)
		}
	}

	old := u.accountReo.SelectByAccount(u.DB, r.Email)
	if old != nil {
		return nil, error2.New(code.AccountExist)
	}
	addData := &org.User{}
	addData.ID = id
	addData.Name = r.Name
	addData.Phone = r.Phone
	addData.Email = r.Email
	addData.SelfEmail = r.SelfEmail
	addData.IDCard = r.IDCard
	addData.Address = r.Address
	addData.UseStatus = r.UseStatus
	addData.Position = r.Position
	addData.Avatar = r.Avatar
	addData.JobNumber = r.JobNumber
	addData.Gender = r.Gender
	addData.CreatedBy = r.CreateBy
	addData.UpdatedBy = r.CreateBy
	addData.Source = r.Source

	addData.CreatedAt = nowUnix
	addData.UpdatedAt = nowUnix
	if r.UseStatus == 0 {
		addData.UseStatus = consts.NormalStatus
	} else {
		addData.UseStatus = r.UseStatus
	}

	addData.PasswordStatus = consts.ResetPasswordStatus
	tx := u.DB.Begin()
	err = u.userRepo.Insert(c, tx, addData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, v := range r.Dep {
		if v.DepID != "" {
			relation := org.UserDepartmentRelation{
				ID:     id2.ShortID(0),
				UserID: id,
				DepID:  v.DepID,
				Attr:   v.Attr,
			}
			err := u.userDepRepo.Add(tx, &relation)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

	}

	for _, v := range r.Leader {
		if v.UserID != "" {
			relation := org.UserLeaderRelation{
				ID:       id2.ShortID(0),
				UserID:   id,
				LeaderID: v.UserID,
				Attr:     v.Attr,
			}
			err := u.userLeaderRepo.Add(tx, &relation)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	account := org.Account{}
	account.Account = r.Email
	account.ID = id2.ShortID(0)
	account.Password = encode2.MD5Encode(r.Password)
	account.UserID = id

	account.CreatedBy = r.Profile.UserID
	account.CreatedAt = nowUnix
	account.UpdatedAt = nowUnix

	err = u.accountReo.Insert(tx, &account)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	adminUser := AddUserResponse{
		ID: id,
	}
	if u.conf.POC {
		adminUser.Password = r.Password
	}
	userTenantRelation := &org.UserTenantRelation{}
	userTenantRelation.ID = id2.ShortID(0)
	userTenantRelation.UserID = id
	err = u.userTenantRepo.Add(c, tx, userTenantRelation)
	if err != nil {
		logger.Logger.Error(err)
	}
	tx.Commit()
	//send message
	if r.SendMessage.SendChannel != NO {
		m := make(map[string]string)
		m[id] = r.Password
		SendAccountAndPWDOrCode(c, u.message, "", r.SendMessage.SendTo, u.conf.MessageTemplate.NewPWD, r.Password, r.SendMessage.SendChannel)
	}
	adminUser.Users = append(adminUser.Users, addData)
	return &adminUser, err
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	//0:null,1:man,2:woman
	Gender   int             `json:"gender,omitempty" `
	Source   string          `json:"source,omitempty" `
	UpdateBy string          `json:"updatedBy,omitempty" `
	Dep      []DepRequest    `json:"dep"`
	Leader   []LeaderRequest `json:"leader" `
}

// UpdateUserResponse update response
type UpdateUserResponse struct {
	ID         string      `json:"id"`
	UpdateUser *org.User   `json:"-"`
	Users      []*org.User `json:"-"`
}

// Update update base info
func (u *user) Update(c context.Context, r *UpdateUserRequest) (*UpdateUserResponse, error) {

	if r.Phone != "" {
		if !verification.CheckPhone(r.Phone) {
			return nil, error2.New(code.InvalidPhone)
		}
	}
	if r.Email != "" {
		if !verification.CheckEmail(r.Email) {
			return nil, error2.New(code.InvalidEmail)
		}
	} else {
		return nil, error2.New(code.EmailRequired)
	}
	if r.SelfEmail != "" {
		if !verification.CheckEmail(r.SelfEmail) {
			return nil, error2.New(code.InvalidEmail)
		}
	}
	for _, v := range r.Leader {
		err := CheckLeader(c, u.DB, u.userLeaderRepo, v.UserID, r.ID)
		if err != nil {
			return nil, error2.New(code.ErrCircleData)
		}
	}
	oldUser := u.userRepo.Get(c, u.DB, r.ID)
	updateData := &org.User{}
	updateData.ID = r.ID
	updateData.Name = r.Name
	updateData.Phone = r.Phone
	updateData.Email = r.Email
	updateData.SelfEmail = r.SelfEmail
	updateData.IDCard = r.IDCard
	updateData.Address = r.Address
	updateData.UseStatus = r.UseStatus
	updateData.Position = r.Position
	updateData.Avatar = r.Avatar
	updateData.JobNumber = r.JobNumber
	updateData.Gender = r.Gender
	updateData.UpdatedBy = r.UpdateBy
	updateData.Source = r.Source
	unix := time2.NowUnix()
	updateData.UpdatedAt = unix
	tx := u.DB.Begin()
	updateAccount := &org.Account{}
	if oldUser.Email != r.Email {
		oldAccount := u.accountReo.SelectByAccount(u.DB, oldUser.Email)
		updateAccount.ID = oldAccount.ID
		updateAccount.Account = r.Email
		updateAccount.UpdatedBy = r.UpdateBy
		updateAccount.UpdatedAt = unix
		err := u.accountReo.Update(u.DB, updateAccount)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := u.userRepo.UpdateByID(c, tx, updateData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(r.Dep) > 0 {
		err = u.userDepRepo.DeleteByUserIDs(tx, r.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, v := range r.Dep {
			relation := org.UserDepartmentRelation{
				ID:     id2.ShortID(0),
				UserID: r.ID,
				DepID:  v.DepID,
				Attr:   v.Attr,
			}
			err := u.userDepRepo.Add(tx, &relation)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if len(r.Leader) > 0 {

		err = u.userLeaderRepo.DeleteByUserIDs(tx, r.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, v := range r.Leader {
			relation := org.UserLeaderRelation{
				ID:       id2.ShortID(0),
				UserID:   r.ID,
				LeaderID: v.UserID,
				Attr:     v.Attr,
			}
			err = u.userLeaderRepo.Add(tx, &relation)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	tx.Commit()
	response := &UpdateUserResponse{ID: r.ID}
	newOld := u.userRepo.Get(c, u.DB, r.ID)
	response.UpdateUser = newOld

	if len(r.Leader) > 0 {
		users := findChild(c, u, r.ID)
		if len(users) > 0 {
			response.Users = append(response.Users, users...)
		}
	}
	return response, nil
}
func findChild(c context.Context, u *user, leaderID ...string) []*org.User {
	userIDs := getChildUser(c, u, leaderID...)
	userIDMap := make(map[string]string)
	for k := range userIDs {
		userIDMap[userIDs[k]] = userIDs[k]
	}
	ids := make([]string, 0)
	for _, v := range userIDMap {
		ids = append(ids, v)
	}
	if len(ids) > 0 {
		return u.userRepo.List(c, u.DB, ids...)
	}
	return nil
}

func getChildUser(c context.Context, u *user, leaderID ...string) []string {
	leaderRelations := u.userLeaderRepo.SelectByLeaderID(u.DB, leaderID...)
	ids := make([]string, 0)
	for k := range leaderRelations {
		ids = append(ids, leaderRelations[k].UserID)
	}
	if len(ids) > 0 {
		ids = append(ids, getChildUser(c, u, ids...)...)
	}
	return ids
}

// UpdateUserAvatarRequest update avatar request
type UpdateUserAvatarRequest struct {
	ID       string `json:"id" binding:"required,max=64"`
	Avatar   string `json:"avatar"`
	UpdateBy string `json:"-"`
}

// UpdateUserAvatarResponse update avatar response
type UpdateUserAvatarResponse struct {
	ID         string    `json:"id" binding:"required,max=64"`
	Avatar     string    `json:"avatar"`
	UpdateBy   string    `json:"-"`
	UpdateUser *org.User `json:"-"`
}

// UpdateAvatar update avatar
func (u *user) UpdateAvatar(c context.Context, r *UpdateUserAvatarRequest) (*UpdateUserAvatarResponse, error) {
	nowUnix := time2.NowUnix()
	old := u.userRepo.Get(c, u.DB, r.ID)
	if old.Avatar != r.Avatar {
		old.Avatar = r.Avatar
		old.UpdatedAt = nowUnix
		old.UpdatedBy = r.UpdateBy
	}

	tx := u.DB.Begin()
	err := u.userRepo.UpdateByID(c, tx, old)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	response := &UpdateUserAvatarResponse{}
	response.UpdateUser = old
	return response, nil
}

// SearchListUserRequest 查询集合条件请求体
type SearchListUserRequest struct {
	DepID  string   `json:"depID" form:"depID"`
	DepIDs []string `json:"depIDs" form:"depIDs"`
	//1:include
	IncludeChildDEPChild int `json:"includeChildDEPChild" form:"includeChildDEPChild" `
	Page                 int `json:"page" form:"page" `
	Limit                int `json:"limit" form:"limit" `
}

// SearchListUserResponse response
type SearchListUserResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"use_status,omitempty" `
	CompanyID string `json:"companyID,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	//0:null,1:man,2:woman
	Gender     int                        `json:"gender,omitempty" `
	Source     string                     `json:"source,omitempty" `
	Attr       string                     `json:"attr,omitempty" `
	DEP        []ManageDepartmentResponse `json:"dep,omitempty"`
	LeaderName []UpdateUserResponse       `json:"leaderName,omitempty"`
}

// ManageDepartmentResponse response
type ManageDepartmentResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	LeaderID  string `json:"leaderID,omitempty"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid,omitempty"`
	SuperPID  string `json:"superID,omitempty"`
	CompanyID string `json:"companyID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
	CreatBy   string `json:"creatBy,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

// PageList page list
func (u *user) PageList(c context.Context, r *SearchListUserRequest) (*page.Page, error) {
	pageRes := &page.Page{}
	userList, total := u.getUsersPageList(c, r)
	if len(userList) > 0 {
		listUserResponses := make([]SearchListUserResponse, 0, len(userList))
		for k := range userList {
			response := SearchListUserResponse{}
			response.ID = userList[k].ID
			response.Name = userList[k].Name
			response.Phone = userList[k].Phone
			response.Email = userList[k].Email
			response.SelfEmail = userList[k].SelfEmail
			response.IDCard = userList[k].IDCard
			response.Address = userList[k].Address
			response.UseStatus = userList[k].UseStatus
			response.Position = userList[k].Position
			response.Avatar = userList[k].Avatar
			response.JobNumber = userList[k].JobNumber
			response.Gender = userList[k].Gender
			response.Source = userList[k].Source
			listUserResponses = append(listUserResponses, response)
		}

		pageRes.Data = listUserResponses
		pageRes.TotalCount = total
	}
	return pageRes, nil
}

func (u *user) getUsersPageList(c context.Context, r *SearchListUserRequest) ([]*org.User, int64) {
	depIDs := make([]string, 0)
	if len(r.DepIDs) > 0 {
		depIDs = append(depIDs, r.DepIDs...)
	} else {
		if r.DepID != "" {
			if r.IncludeChildDEPChild != 1 {
				depIDs = append(depIDs, r.DepID)
			} else {
				depIDs = u.getChildDep(c, r.DepID, depIDs, consts.NormalStatus)
			}
		}
	}
	var userIDs = make([]string, 0)
	if len(depIDs) > 0 {
		relations := u.userDepRepo.SelectByDEPID(u.DB, depIDs...)
		for k := range relations {
			userIDs = append(userIDs, relations[k].UserID)
		}
	}

	list, total := u.userRepo.PageList(c, u.DB, consts.NormalStatus, r.Page, r.Limit, userIDs)

	return list, total
}

// DepOneResponse response
type DepOneResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	//1:company,2:department
	Attr  int              `json:"attr,omitempty"`
	Child []DepOneResponse `json:"child,omitempty"`
}

// SearchOneUserRequest get by id
type SearchOneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// SearchOneUserResponse 管理员可见字段
type SearchOneUserResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	//0:null,1:man,2:woman
	Gender int                `json:"gender,omitempty" `
	Source string             `json:"source,omitempty" `
	DEP    [][]DepOneResponse `json:"dep,omitempty"`
	Leader [][]Leader         `json:"leader,omitempty"`
}

// AdminSelectByID 管理员根据ID查询人员
func (u *user) AdminSelectByID(c context.Context, r *SearchOneUserRequest) (*SearchOneUserResponse, error) {
	old := u.userRepo.Get(c, u.DB, r.ID)

	if old != nil {
		_, filter := u.columnRepo.GetFilter(c, u.DB, consts.FieldAdminStatus, consts.AllAttr)
		if filter != nil {
			Filter(old, filter, OUT)
		}
		resUser := SearchOneUserResponse{}
		resUser.ID = old.ID
		resUser.Name = old.Name
		resUser.Phone = old.Phone
		resUser.Email = old.Email
		resUser.SelfEmail = old.SelfEmail
		resUser.UseStatus = old.UseStatus
		resUser.Position = old.Position
		resUser.Avatar = old.Avatar
		resUser.JobNumber = old.JobNumber
		resUser.Gender = old.Gender
		//从当前部门一直找到顶层部门组装
		list, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
		depMap := make(map[string]*org.Department)
		for k := range list {
			depMap[list[k].ID] = &list[k]
		}
		departmentRelations := u.userDepRepo.SelectByUserIDs(u.DB, old.ID)
		depIDs := make([]string, 0)
		for k := range departmentRelations {
			depIDs = append(depIDs, departmentRelations[k].DepID)
		}
		departments := u.depRepo.List(c, u.DB, depIDs...)

		if len(departments) > 0 {
			for k := range departments {
				responses := make([]DepOneResponse, 0)
				resUser.DEP = append(resUser.DEP, FindDepToTop(depMap, departments[k].ID, responses))
			}

		}
		leader, err := makeLeaderToTop(c, u, old.ID, old.ID)
		if err == nil && leader != nil {
			resUser.Leader = append(resUser.Leader, leader...)
		}

		return &resUser, nil
	}
	return nil, nil
}

// ViewerSearchOneUserRequest 查询一个
type ViewerSearchOneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// ViewerSearchOneUserResponse 用户可见字段
type ViewerSearchOneUserResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	//0:null,1:man,2:woman
	Gender int                `json:"gender,omitempty" `
	Source string             `json:"source,omitempty" `
	Dep    [][]DepOneResponse `json:"deps,omitempty"`
	Leader [][]Leader         `json:"leaders,omitempty"`
	// 0x1111 right first 0:need reset password
	Status int `json:"status"`
}

// UserSelectByID user get by id
func (u *user) UserSelectByID(c context.Context, r *ViewerSearchOneUserRequest) (*ViewerSearchOneUserResponse, error) {
	old := u.userRepo.Get(c, u.DB, r.ID)

	if old != nil {
		_, filter := u.columnRepo.GetFilter(c, u.DB, consts.FieldViewerStatus, consts.AllAttr)
		if filter != nil {
			Filter(old, filter, OUT)
		}
		resUser := ViewerSearchOneUserResponse{}
		resUser.ID = old.ID
		resUser.Name = old.Name
		resUser.Phone = old.Phone
		resUser.Email = old.Email
		resUser.SelfEmail = old.SelfEmail
		resUser.UseStatus = old.UseStatus
		resUser.Position = old.Position
		resUser.Avatar = old.Avatar
		resUser.JobNumber = old.JobNumber
		resUser.Gender = old.Gender
		if old.PasswordStatus&consts.NormalStatus == 0 {
			resUser.Status = (resUser.Status >> 1) << 1
		} else {
			resUser.Status = ((resUser.Status >> 1) << 1) + 1
		}
		list, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
		depMap := make(map[string]*org.Department)
		for k := range list {
			depMap[list[k].ID] = &list[k]
		}
		departmentRelations := u.userDepRepo.SelectByUserIDs(u.DB, old.ID)
		depIDs := make([]string, 0)
		for k := range departmentRelations {
			depIDs = append(depIDs, departmentRelations[k].DepID)
		}
		departments := u.depRepo.List(c, u.DB, depIDs...)
		if len(departments) > 0 {
			for k := range departments {
				responses := make([]DepOneResponse, 0)
				resUser.Dep = append(resUser.Dep, FindDepToTop(depMap, departments[k].ID, responses))
			}

		}
		leader, err := makeLeaderToTop(c, u, old.ID, old.ID)
		if err == nil && leader != nil {
			resUser.Leader = append(resUser.Leader, leader...)
		}
		return &resUser, nil
	}
	return nil, error2.New(code.DataNotExist)

}

// SearchUserByIDsRequest get by ids
type SearchUserByIDsRequest struct {
	IDs []string `json:"ids" form:"ids"  binding:"required"`
}

// SearchUserByIDsResponse get by ids response
type SearchUserByIDsResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	IDCard    string `json:"idCard,omitempty" `
	Address   string `json:"address,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	//0:null,1:man,2:woman
	Gender int    `json:"gender,omitempty" `
	Source string `json:"source,omitempty" `
	// 0x1111 right first 0:need reset password
	Status int              `json:"status"`
	DEP    []DepOneResponse `json:"dep,omitempty"`
}

// UserSelectByIDs user get by ids
func (u *user) UserSelectByIDs(c context.Context, r *SearchUserByIDsRequest) ([]SearchUserByIDsResponse, error) {
	list := u.userRepo.List(c, u.DB, r.IDs...)
	res := make([]SearchUserByIDsResponse, 0)
	if len(list) > 0 {
		_, filter := u.columnRepo.GetFilter(c, u.DB, consts.FieldViewerStatus, consts.AllAttr)
		if filter != nil {
			Filter(&list, filter, OUT)
		}
		for k := range list {
			userResponse := SearchUserByIDsResponse{}
			userResponse.ID = list[k].ID
			userResponse.Name = list[k].Name
			userResponse.Phone = list[k].Phone
			userResponse.Email = list[k].Email
			userResponse.SelfEmail = list[k].SelfEmail
			userResponse.UseStatus = list[k].UseStatus
			userResponse.Position = list[k].Position
			userResponse.Avatar = list[k].Avatar
			userResponse.JobNumber = list[k].JobNumber
			userResponse.Gender = list[k].Gender

			res = append(res, userResponse)
		}
		return res, nil
	}
	return nil, error2.New(code.DataNotExist)
}

// StatusRequest update one status request
type StatusRequest struct {
	ID string `json:"id" binding:"required" binding:"max=64"`
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus" binding:"required,max=64"`
	UpdatedBy string `json:"updatedBy"`
	Profile   header2.Profile
}

// StatusResponse response
type StatusResponse struct {
	User *org.User `json:"-"`
}

// UpdateUserStatus update one user status
func (u *user) UpdateUserStatus(c context.Context, r *StatusRequest) (*StatusResponse, error) {
	old := u.userRepo.Get(c, u.DB, r.ID)
	if old == nil {
		return nil, error2.New(code.DataNotExist)
	}
	if old.ID == r.Profile.UserID {
		return nil, error2.New(code.CanNotModifyYourself)
	}
	account := org.Account{}
	nowUnix := time2.NowUnix()
	tx := u.DB.Begin()
	old.UseStatus = r.UseStatus
	old.UpdatedAt = nowUnix
	old.UpdatedBy = r.Profile.UserID

	if old.UseStatus != consts.ActiveStatus && r.UseStatus == consts.ActiveStatus {
		return nil, error2.New(code.ErrHasBeActive)
	}

	if r.UseStatus == consts.ActiveStatus {
		old.UseStatus = consts.NormalStatus
	}

	err := u.userRepo.UpdateByID(c, u.DB, old)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	account.ID = r.ID
	pwd := ""
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)
	if r.UseStatus == consts.ActiveStatus {
		pwd = random2.RandomString(int(info.PwdMinLen), info.PwdType)
		account.Password = encode2.MD5Encode(pwd)
	}

	err = u.accountReo.Update(u.DB, &account)
	if r.UseStatus == consts.DelStatus {
		err = u.userDepRepo.DeleteByUserIDs(tx, r.ID)
		err = u.userRepo.UpdateByID(c, tx, old)
		err = u.accountReo.DeleteByUserID(tx, r.ID)
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	if r.UseStatus == consts.DelStatus {
		delRequest := &goalie.OthDelRequest{
			IDs:   []string{r.ID},
			DelBy: r.UpdatedBy,
		}
		_, err := u.goalie.DelOwner(c, delRequest)
		if err != nil {
			logger.Logger.Error("del user role from goalie err", err)
		}
	}

	return &StatusResponse{User: old}, nil
}

// ListStatusRequest update list user status request
type ListStatusRequest struct {
	IDS []string `json:"ids" binding:"required"`
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus" binding:"required"`
	UpdatedBy string `json:"updatedBy"`
}

// ListStatusResponse update list user status response
type ListStatusResponse struct {
	Users []*org.User `json:"-"`
}

// UpdateUsersStatus update list user status
func (u *user) UpdateUsersStatus(c context.Context, r *ListStatusRequest) (*ListStatusResponse, error) {
	if r.UseStatus == consts.DelStatus {
		return nil, error2.New(code.BatchDeleteNotSupported)
	}
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)
	pwds := make(map[string]string)
	users := make([]*org.User, 0)
	ids := make([]string, 0)
	for _, v := range r.IDS {
		old := u.userRepo.Get(c, u.DB, v)
		if old == nil {
			continue
		}
		if v == r.UpdatedBy {
			continue
		}
		account := org.Account{}
		nowUnix := time2.NowUnix()
		tx := u.DB.Begin()
		old.UseStatus = r.UseStatus
		old.UpdatedAt = nowUnix
		old.UpdatedBy = r.UpdatedBy

		if old.UseStatus != consts.ActiveStatus && r.UseStatus == consts.ActiveStatus {
			return nil, error2.New(code.ErrHasBeActive)
		}

		if r.UseStatus == consts.ActiveStatus {
			old.UseStatus = consts.NormalStatus
		}

		err := u.userRepo.UpdateByID(c, u.DB, old)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		account.ID = v
		account.UpdatedAt = nowUnix
		if r.UseStatus == consts.ActiveStatus {
			pwd := random2.RandomString(int(info.PwdMinLen), info.PwdType)
			account.Password = encode2.MD5Encode(pwd)
			pwds[account.ID] = pwd
		}
		err = u.accountReo.Update(u.DB, &account)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		users = append(users, old)
		ids = append(ids, v)

	}
	if r.UseStatus == consts.DelStatus {
		delRequest := &goalie.OthDelRequest{
			IDs:   ids,
			DelBy: r.UpdatedBy,
		}
		_, err := u.goalie.DelOwner(c, delRequest)
		if err != nil {
			logger.Logger.Error("del user role from goalie err", err)
		}
	}
	response := &ListStatusResponse{}
	if len(users) > 0 {
		response.Users = append(response.Users, users...)
	}
	go func() {
		if r.UseStatus == consts.ActiveStatus {
			users := u.userRepo.List(c, u.DB, r.IDS...)
			for k := range users {
				pwd := pwds[users[k].ID]
				SendAccountAndPWDOrCode(c, u.message, "", users[k].SelfEmail, u.conf.MessageTemplate.NewPWD, pwd, 1)
			}
		}
	}()

	return response, nil
}

// ChangeUsersDEPRequest change user dep request
type ChangeUsersDEPRequest struct {
	UsersID  []string `json:"usersID"  binding:"required"`
	OldDepID string   `json:"oldDepID"  binding:"required,max=64"`
	NewDepID string   `json:"newDepID"  binding:"required,max=64"`
}

// ChangeUsersDEPResponse change user dep response
type ChangeUsersDEPResponse struct {
	Users []*org.User `json:"-"`
}

// AdminChangeUsersDEP change list user dep
func (u *user) AdminChangeUsersDEP(c context.Context, rq *ChangeUsersDEPRequest) (*ChangeUsersDEPResponse, error) {
	tx := u.DB.Begin()
	for _, v := range rq.UsersID {
		oldRelation := u.userDepRepo.SelectByUserIDAndDepID(u.DB, v, rq.OldDepID)
		if oldRelation != nil {
			oldRelation.DepID = rq.NewDepID
			err := u.userDepRepo.Update(tx, oldRelation)
			if err != nil {
				tx.Rollback()
				return nil, error2.New(code.ChangeDepErr)
			}
			u.redisClient.Del(c, consts.RedisTokenUserInfo+v)
		}

	}
	tx.Commit()
	users := u.userRepo.List(c, u.DB, rq.UsersID...)
	response := &ChangeUsersDEPResponse{}
	if len(users) > 0 {
		response.Users = append(response.Users, users...)
	}
	return response, nil
}

func (u *user) getChildDep(c context.Context, depID string, depIDs []string, status int) []string {
	depIDs = append(depIDs, depID)
	list, _ := u.depRepo.SelectByPID(c, u.DB, depID, status, 1, 10000)
	if len(list) > 0 {
		for k := range list {
			depIDs = u.getChildDep(c, list[k].ID, depIDs, status)
		}

		return depIDs
	}
	return depIDs

}

// TokenUserRequest get one user request
type TokenUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// TokenUserResponse get one user response
type TokenUserResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int                `json:"use_status,omitempty" `
	TenantID  string             `json:"tenantID,omitempty" `
	Position  string             `json:"position,omitempty" `
	Avatar    string             `json:"avatar,omitempty" `
	JobNumber string             `json:"jobNumber,omitempty" `
	Status    int                `json:"status"`
	DEP       [][]DepOneResponse `json:"deps,omitempty"`
	Leader    [][]Leader         `json:"leaders,omitempty"`
}

// Leader leader
type Leader struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
}

// TenantID tenant id
const TenantID = "Tenant-Id"

// OthGetOneUser othGetOneUser
func (u *user) OthGetOneUser(c context.Context, rq *TokenUserRequest) (*TokenUserResponse, error) {
	userBtye := u.redisClient.Get(c, consts.RedisTokenUserInfo+rq.ID).Val()
	userUser := TokenUserResponse{}
	if userBtye != "" {
		err := json.Unmarshal([]byte(userBtye), &userUser)
		if err != nil {
			return nil, err
		}
		return &userUser, nil
	}
	old := u.userRepo.Get(c, u.DB, rq.ID)
	if old != nil {
		userUser.ID = old.ID
		userUser.Name = old.Name
		userUser.Phone = old.Phone
		userUser.Email = old.Email
		userUser.SelfEmail = old.SelfEmail
		userUser.UseStatus = old.UseStatus
		userUser.TenantID = old.TenantID
		userUser.Position = old.Position
		userUser.Avatar = old.Avatar
		userUser.JobNumber = old.JobNumber

		c = header2.SetContext(c, TenantID, old.TenantID)

		list, _ := u.depRepo.PageList(c, u.DB, consts.NormalStatus, 1, 10000)
		depMap := make(map[string]*org.Department)
		for k := range list {
			depMap[list[k].ID] = &list[k]
		}
		departmentRelations := u.userDepRepo.SelectByUserIDs(u.DB, old.ID)
		depIDs := make([]string, 0)
		for k := range departmentRelations {
			depIDs = append(depIDs, departmentRelations[k].DepID)
		}
		departments := u.depRepo.List(c, u.DB, depIDs...)

		if len(departments) > 0 {
			for k := range departments {
				responses := make([]DepOneResponse, 0)
				userUser.DEP = append(userUser.DEP, FindDepToTop(depMap, departments[k].ID, responses))
			}

		}
		leader, err := makeLeaderToTop(c, u, old.ID, old.ID)
		if err == nil && leader != nil {
			userUser.Leader = append(userUser.Leader, leader...)
		}

		marshal, err := json.Marshal(userUser)
		if err != nil {
			return nil, err
		}
		u.redisClient.SetEX(c, consts.RedisTokenUserInfo+rq.ID, marshal, consts.RedisTokenUserInfoEx*time.Minute)
		return &userUser, nil
	}
	return nil, error2.New(code.DataNotExist)
}

func makeLeaderToTop(c context.Context, u *user, userID, startUserID string) ([][]Leader, error) {
	relations := u.userLeaderRepo.SelectByUserIDs(u.DB, userID)
	if len(relations) > 0 {
		res := make([][]Leader, 0)
		for k := range relations {
			if relations[k].LeaderID == startUserID {
				return nil, errors.New("circle leader")
			}
			if relations[k].LeaderID != "" {
				ls := make([]Leader, 0)
				get := u.userRepo.Get(c, u.DB, relations[k].LeaderID)
				if get == nil {
					return res, nil
				}
				leader := Leader{}
				leader.ID = get.ID
				leader.Name = get.Name
				leader.Email = get.Email
				ls = append(ls, leader)
				array, err := makeLeaderToTop(c, u, get.ID, startUserID)
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
		return res, nil
	}
	return nil, nil

}

// CheckLeader check relation circle
func CheckLeader(c context.Context, db *gorm.DB, ur org.UserLeaderRelationRepo, userID, startUserID string) error {
	if userID == startUserID {
		return errors.New("circle leader")
	}
	relations := ur.SelectByUserIDs(db, userID)
	if len(relations) > 0 {
		for k := range relations {
			if relations[k].LeaderID != "" {
				err := CheckLeader(c, db, ur, relations[k].LeaderID, startUserID)
				if err != nil {
					return err
				}
			}

		}
		return nil
	}
	return nil

}

// FindDepToTop find department from here to top
func FindDepToTop(AllDep map[string]*org.Department, depID string, res []DepOneResponse) []DepOneResponse {

	if v, ok := AllDep[depID]; v != nil && ok {
		online := DepOneResponse{}
		online.ID = v.ID
		online.Name = v.Name
		online.PID = v.PID
		online.SuperPID = v.SuperPID
		online.Grade = v.Grade
		online.Attr = v.Attr
		res = append(res, online)
		if v.PID != "" {
			return FindDepToTop(AllDep, v.PID, res)
		}
		return res
	}
	return res
}

// GetTemplateFileRequest temp file
type GetTemplateFileRequest struct {
}

// GetTemplateFileResponse temp file
type GetTemplateFileResponse struct {
	Data        []byte            `json:"data"`
	FileName    string            `json:"fileName"`
	ExcelColumn map[string]string `json:"excelColumn"`
}

// Template get xlsx template
func (u *user) Template(c context.Context, r *GetTemplateFileRequest) (*GetTemplateFileResponse, error) {
	xlsxFields := u.columnRepo.GetXlsxField(c, u.DB, consts.FieldAdminStatus)
	if xlsxFields == nil || len(xlsxFields) == 0 {
		return nil, error2.New(code.FieldNameIsNull)
	}
	newFile := xlsx.NewFile()
	sheet, err := newFile.AddSheet("sheet1")
	if err != nil {
		return nil, err
	}
	row := sheet.AddRow()
	s := make([]string, 0)
	for k, v := range xlsxFields {
		if v != consts.ID {
			s = append(s, k)
		}
	}
	s = append(s, consts.OwnerDepName)
	sort.Strings(s)
	for k := range s {
		cell := row.AddCell()
		cell.SetValue(s[k])
	}
	row2 := sheet.AddRow()
	for k := range s {
		addCell := row2.AddCell()
		addCell.SetValue("demo(导入请删除此行)")
		if s[k] == consts.OwnerDepName {
			addCell.SetValue("/部门demo1/子部门demo2")
		}
	}
	buffer := new(bytes.Buffer)
	newFile.Write(buffer)
	res := &GetTemplateFileResponse{}
	res.Data = buffer.Bytes()
	res.FileName = u.conf.TemplateName
	return res, nil
}

// IndexCountRequest count
type IndexCountRequest struct {
}

// IndexCountResponse count
type IndexCountResponse struct {
	UserTotal     int64 `json:"userTotal"`
	DepNumber     int64 `json:"depNumber"`
	ActiveUserNum int64 `json:"activeUserNum"`
}

// IndexCount count
func (u *user) IndexCount(c context.Context, r *IndexCountRequest) (*IndexCountResponse, error) {
	totalUser, activeUserNum := u.userRepo.Count(c, u.DB, consts.NormalStatus, consts.ActiveStatus)
	depNum := u.depRepo.Count(c, u.DB, consts.NormalStatus)
	indexCount := &IndexCountResponse{
		UserTotal:     totalUser,
		ActiveUserNum: activeUserNum,
		DepNumber:     depNum,
	}
	return indexCount, nil
}

// SendAccountAndPWDOrCode sendType 第一位发邮件，第二位发手机
func SendAccountAndPWDOrCode(c context.Context, messageClient message.Message, email, selfEmail, messageTemple, data string, sendType int) {
	emailReq := make([]*message.CreateReq, 0)
	//send email
	if sendType&1 == 1 {
		mesReq := new(message.CreateReq)
		if selfEmail != "" {
			mesReq.Email = &message.Email{
				To: []string{selfEmail},
				Content: &message.Content{
					TemplateID: messageTemple,
					KeyAndValue: map[string]string{
						"code":    data,
						"account": email,
					},
				},
			}
		} else {
			mesReq.Email = &message.Email{
				To: []string{email},
				Content: &message.Content{
					TemplateID: messageTemple,
					KeyAndValue: map[string]string{
						"code": data,
					},
				},
			}
		}

		emailReq = append(emailReq, mesReq)
	}
	//send phone
	if sendType>>1&1 == 1 {
		mesReq := new(message.CreateReq)
		mesReq.Phone = &message.Phone{}
	}
	go func() {
		err := messageClient.SendMessage(c, emailReq)
		if err != nil {
			logger.Logger.Error(err)
		}
	}()

}

// AddUserToLdap add user to ldap
func AddUserToLdap(ctx context.Context, header http.Header, ldapClient ldap.Ldap, id, email, name, password string, jobNumber int, gidNumber string) error {
	addReq := &ldap.UserAddReq{}
	addReq.ID = id
	addReq.UserPassword = password
	addReq.Mail = email
	split := strings.Split(email, "@")
	addReq.UserName = name
	addReq.FirstName = split[0]
	addReq.LastName = split[0]
	addReq.UIDNumber = jobNumber
	addReq.GidNumber = gidNumber
	_, err := ldapClient.AddUser(ctx, header, addReq)
	return err
}

// RegisterRequest register request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Header   http.Header
	Profile  header2.Profile
}

// RegisterResponse register response
type RegisterResponse struct {
	User *org.User `json:"-"`
}

// Register register
func (u *user) Register(c context.Context, r *RegisterRequest) (*RegisterResponse, error) {
	val := u.redisClient.Get(c, u.conf.VerificationCode.RegisterCode+":"+r.Email).Val()
	if val == "" {
		return nil, error2.New(code.ExpireVerificationCode)
	}
	if val != r.Code {
		return nil, error2.New(code.InvalidVerificationCode)
	}
	id := id2.HexUUID(true)
	nowUnix := time2.NowUnix()
	if r.Email != "" {
		if !verification.CheckEmail(r.Email) {
			return nil, error2.New(code.InvalidEmail)
		}
	} else {
		return nil, error2.New(code.EmailRequired)
	}

	addData := &org.User{}
	addData.ID = id
	addData.Name = r.Name

	addData.Email = r.Email

	old := u.accountReo.SelectByAccount(u.DB, addData.Email)
	if old != nil {
		return nil, error2.New(code.AccountExist)
	}
	addData.CreatedAt = nowUnix
	addData.UpdatedAt = nowUnix
	//todo tenantID is tenant server response

	addData.UseStatus = consts.NormalStatus
	addData.PasswordStatus = consts.NormalStatus
	tx := u.DB.Begin()
	registerRequest := &landlord.RegisterRequest{}
	registerRequest.OwnerID = id
	registerRequest.Name = r.Name
	registerResponse, err := u.landlord.Register(c, r.Header, registerRequest)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	c = header2.SetContext(c, TenantID, registerResponse.ID)

	err = u.userRepo.Insert(c, tx, addData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	account := org.Account{}
	account.Account = r.Email
	account.ID = id2.ShortID(0)
	account.UserID = id
	account.Password = encode2.MD5Encode(r.Password)
	account.CreatedBy = id
	account.CreatedAt = nowUnix
	account.UpdatedAt = nowUnix
	err = u.accountReo.Insert(tx, &account)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	userTenantRelation := &org.UserTenantRelation{}
	userTenantRelation.ID = id2.ShortID(0)
	userTenantRelation.UserID = id
	err = u.userTenantRepo.Add(c, tx, userTenantRelation)
	if err != nil {
		logger.Logger.Error(err)
	}

	tx.Commit()

	return &RegisterResponse{User: addData}, nil
}

// CreatePassword create password
func CreatePassword(c context.Context, conf configs.Config, redisClient redis.UniversalClient) string {
	info := systems.GetSecurityInfo(c, conf, redisClient)
	pwd := random2.RandomString(int(info.PwdMinLen), info.PwdType)
	if conf.Model == "debug" {
		pwd = "654321a.."
		return pwd
	}
	return pwd
}

// GetUsersByIDsRequest request
type GetUsersByIDsRequest struct {
	IDs []string `json:"ids"`
}

// GetUsersByIDsResponse response
type GetUsersByIDsResponse struct {
	Users []ViewerSearchOneUserResponse `json:"users"`
}

// GetUsersByIDs get users by user's ids
func (u *user) GetUsersByIDs(c context.Context, r *GetUsersByIDsRequest) (*GetUsersByIDsResponse, error) {
	list := u.userRepo.List(c, u.DB, r.IDs...)
	if len(list) == 0 {
		return nil, nil
	}
	response := &GetUsersByIDsResponse{}
	for k := range list {
		userResponse := ViewerSearchOneUserResponse{}
		userResponse.ID = list[k].ID
		userResponse.Name = list[k].Name
		userResponse.Phone = list[k].Phone
		userResponse.Email = list[k].Email
		userResponse.SelfEmail = list[k].SelfEmail
		userResponse.UseStatus = list[k].UseStatus
		userResponse.Position = list[k].Position
		userResponse.Avatar = list[k].Avatar
		userResponse.JobNumber = list[k].JobNumber
		userResponse.Gender = list[k].Gender
		response.Users = append(response.Users, userResponse)
	}
	relations := u.userDepRepo.SelectByUserIDs(u.DB, r.IDs...)
	ud := make(map[string][]string)
	depIDs := make([]string, 0)
	for k := range relations {
		ud[relations[k].UserID] = append(ud[relations[k].UserID], relations[k].DepID)
		depIDs = append(depIDs, relations[k].DepID)
	}
	departments := u.depRepo.List(c, u.DB, depIDs...)
	depMap := make(map[string]org.Department)
	for k := range departments {
		depMap[departments[k].ID] = departments[k]
	}
	for k := range response.Users {
		for k1 := range ud[response.Users[k].ID] {
			depOneResponses := make([]DepOneResponse, 0)
			oneResponse := DepOneResponse{}
			oneResponse.ID = depMap[ud[response.Users[k].ID][k1]].ID
			oneResponse.Name = depMap[ud[response.Users[k].ID][k1]].Name
			oneResponse.PID = depMap[ud[response.Users[k].ID][k1]].PID
			oneResponse.Attr = depMap[ud[response.Users[k].ID][k1]].Attr
			oneResponse.UseStatus = depMap[ud[response.Users[k].ID][k1]].UseStatus
			oneResponse.Grade = depMap[ud[response.Users[k].ID][k1]].Grade
			depOneResponses = append(depOneResponses, oneResponse)
			response.Users[k].Dep = append(response.Users[k].Dep, depOneResponses)
		}
	}
	return response, nil
}

// ImportFileRequest 上传文件
type ImportFileRequest struct {
	UseStatus int    `json:"useStatus" form:"useStatus" binding:"required,max=64"` //状态：1正常，-2禁用，-1删除，2激活==1 （与账号库相同）
	IsUpdate  int    `json:"isUpdate" form:"isUpdate" `                            //1更新旧数据，-1不更新只插入新数据
	TenantID  string `json:"tenantID"`
}

// ImportFileResponse 文件导入结果
type ImportFileResponse struct {
	AddSuccessTotal    int                      `json:"addSuccessTotal"`
	AddData            []map[string]interface{} `json:"addData"`
	UpdateSuccessTotal int                      `json:"updateSuccessTotal"`
	UpdateData         []map[string]interface{} `json:"updateData"`
	FailTotal          int                      `json:"failTotal"`
	FailUsers          []map[string]interface{} `json:"failUsers"`
	Users              []*org.User              `json:"-"`
}

// ImportFile 上传文件
func (u *user) ImportFile(c context.Context, file []byte, profile header2.Profile, r *ImportFileRequest) (*ImportFileResponse, error) {
	fail := make([]map[string]interface{}, 0)
	//1、开始解析excel文件
	suc1, err := u.makeDataFromExcl(c, file, r.TenantID)
	if err != nil {
		return nil, err
	}
	//2、第二次遍历suc1中的数据，找出邮箱或者手机重复的，剩下的才能进行插入操作
	suc2, fails := u.screenUserData(c, suc1, r)
	fail = append(fail, fails...)
	//3、开始执行插入操作
	suc, fail2, updates, userList := u.insertList(c, suc2, profile.UserID, r)
	fail = append(fail, fail2...)
	updateSucs := make([]map[string]interface{}, 0)
	//4、执行跟新操作
	var updateTotal = 0
	if r.IsUpdate == isUpdate {
		updateSuc, fail3, users := u.updateList(c, updates, r.TenantID)
		fail = append(fail, fail3...)
		updateSucs = append(updateSucs, updateSuc...)
		updateTotal = len(updateSuc)
		userList = append(userList, users...)
	}

	result := ImportFileResponse{
		AddSuccessTotal:    len(suc),
		AddData:            suc,
		UpdateSuccessTotal: updateTotal,
		UpdateData:         updateSucs,
		FailTotal:          len(fail),
		FailUsers:          fail,
		Users:              userList,
	}

	return &result, nil
}

// 1、从excl组装数据
func (u *user) makeDataFromExcl(ctx context.Context, file []byte, tenantID string) ([]map[string]interface{}, error) {
	xlFile, _ := xlsx.OpenBinary(file)
	if xlFile == nil {
		return nil, error2.New(code.InvalidFile)
	}
	suc1 := make([]map[string]interface{}, 0)
	xlsxFields := u.columnRepo.GetXlsxField(ctx, u.DB, consts.FieldAdminStatus)
	departments, _ := u.depRepo.PageList(ctx, u.DB, consts.NormalStatus, 1, 10000)
	depMap := make(map[string]string)
	for _, v := range departments {
		depMap[v.Name] = v.ID
	}
	sheet := xlFile.Sheets[0]
	for k, row := range sheet.Rows {
		fmt.Println("k===", k)
		cells0 := sheet.Rows[0].Cells
		if k >= 1 {
			s := make(map[string]interface{})
			for k1, v1 := range row.Cells {
				fmt.Println("k1===", k1)
				if cells0[k1].Value == consts.OwnerDepName {
					s[consts.DEPNAME] = v1.Value
				} else {
					s[xlsxFields[cells0[k1].Value]] = v1.Value
				}
			}
			if len(s) > 0 {
				suc1 = append(suc1, s)
			}

		}
	}
	return suc1, nil
}

// 2、对数据进行组装，判断出需新增和更新的数据
func (u *user) screenUserData(ctx context.Context, suc1 []map[string]interface{}, r *ImportFileRequest) (suc, fails []map[string]interface{}) {
	m := make(map[string]int)
	fail := make([]map[string]interface{}, 0)
	suc2 := make([]map[string]interface{}, 0)
	suc3 := make([]map[string]interface{}, 0)
	for k := range suc1 {
		if suc1[k][consts.EMAIL].(string) == "" || suc1[k][consts.SELFEMAIL].(string) == "" {
			suc1[k][consts.REMARK] = consts.NamePhoneEmailNotNull
			fail = append(fail, suc1[k])
			continue
		}
		suc2 = append(suc2, suc1[k])
		m[suc1[k][consts.EMAIL].(string)] = m[suc1[k][consts.EMAIL].(string)] + 1
		m[suc1[k][consts.PHONE].(string)] = m[suc1[k][consts.PHONE].(string)] + 1
	}

	//获取部门前缀树
	depRouter := department.NewDepartmentRouter()
	list, _ := u.depRepo.PageList(ctx, u.DB, consts.NormalStatus, 1, 10000)
	depRouter.AddRoute(list)
A:
	for k := range suc2 {
		if m[suc2[k][consts.EMAIL].(string)] <= 1 && m[suc2[k][consts.PHONE].(string)] <= 1 {
			for k1 := range suc2[k] {
				switch k1 {
				case consts.EMAIL, consts.NAME, consts.SELFEMAIL:
					if v1, ok := suc2[k][k1]; !ok || v1.(string) == "" {
						suc2[k][consts.REMARK] = consts.NamePhoneEmailNotNull
						fail = append(fail, suc2[k])
						continue A
					}
					switch k1 {
					case consts.EMAIL, consts.SELFEMAIL:
						if !verification.CheckEmail(suc2[k][k1].(string)) {
							suc2[k][consts.REMARK] = consts.NotEmail
							fail = append(fail, suc2[k])
							continue A
						}
					case consts.PHONE:
						if !verification.CheckPhone(suc2[k][k1].(string)) {
							suc2[k][consts.REMARK] = consts.NotPhone
							fail = append(fail, suc2[k])
							continue A
						}
					}
				case consts.DEPNAME:
					node := depRouter.GetRoute(suc2[k][k1].(string))
					if node == nil {
						suc2[k][consts.REMARK] = consts.NotDepartment
						fail = append(fail, suc2[k])
						continue A
					}
					suc2[k][consts.DEPID] = node.DepID
				}
			}
			suc3 = append(suc3, suc2[k])
			continue
		}
		suc2[k][consts.REMARK] = consts.EmailPhoneRepeat
		fail = append(fail, suc2[k])
	}
	return suc3, fail
}

const (
	isUpdate = 1
)

func (u *user) insertList(ctx context.Context, suc2 []map[string]interface{}, createBy string, r *ImportFileRequest) (suc, fails, updates []map[string]interface{}, userList []*org.User) {
	var i = 0
	fail := make([]map[string]interface{}, 0)
	update := make([]map[string]interface{}, 0)
	su := make([]map[string]interface{}, 0)
	users := make([]*org.User, 0)
	//获取过滤字段
	//_, filters := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldAdminStatus, consts.AllAttr)
	for k := range suc2 {
		res := u.userRepo.SelectByEmailOrPhone(ctx, u.DB, suc2[k][consts.EMAIL].(string))
		if res != nil {
			suc2[k][consts.ID] = res.ID
			if r.IsUpdate == isUpdate {
				suc2[k][consts.USESTATUS] = r.UseStatus
				update = append(update, suc2[k])
			} else {
				suc2[k][consts.REMARK] = consts.EmailPhoneExist
				fail = append(fail, suc2[k])
			}
			continue
		}

		tx := u.DB.Begin()
		id := id2.HexUUID(true)
		nowUnix := time2.NowUnix()
		suc2[k][consts.ID] = id
		suc2[k][consts.USESTATUS] = r.UseStatus

		depID := suc2[k][consts.DEPID].(string)
		//delete(suc2[k], consts.DEPNAME)
		//delete(suc2[k], consts.DEPID)
		u2 := &org.User{}
		marshal, err := json.Marshal(suc2[k])
		if err != nil {
			fail = append(fail, suc2[k])
			continue
		}
		err = json.Unmarshal(marshal, u2)
		if err != nil {
			fail = append(fail, suc2[k])
			continue
		}
		u2.CreatedAt = nowUnix
		u2.UpdatedAt = nowUnix
		u2.CreatedBy = createBy
		u2.UpdatedBy = createBy
		err = u.userRepo.Insert(ctx, tx, u2)
		//Filter(suc2, filters, IN)
		if err != nil {
			suc2[k][consts.REMARK] = consts.EmailPhoneExist
			delete(suc2[k], consts.ID)
			fail = append(fail, suc2[k])
			tx.Rollback()
			continue
		}
		password := CreatePassword(ctx, u.conf, u.redisClient)
		account := org.Account{
			ID:        id2.HexUUID(true),
			Password:  encode2.MD5Encode(password),
			Account:   u2.Email,
			CreatedAt: nowUnix,
			UpdatedAt: nowUnix,
			CreatedBy: createBy,
			UserID:    u2.ID,
		}
		err = u.accountReo.Insert(tx, &account)
		if err != nil {
			suc2[k][consts.REMARK] = consts.EmailPhoneExist
			delete(suc2[k], consts.ID)
			fail = append(fail, suc2[k])
			tx.Rollback()
			continue
		}

		err = u.updateUserDepRelation(tx, id, depID)
		if err != nil {
			suc2[k][consts.REMARK] = consts.RelationDepartmentFail
			delete(suc2[k], consts.ID)
			fail = append(fail, suc2[k])
			tx.Rollback()
			continue
		}
		su = append(su, suc2[k])
		tx.Commit()
		i = i + 1
		users = append(users, u2)
	}
	return su, fail, update, users
}

func (u *user) updateList(ctx context.Context, list []map[string]interface{}, tenantID string) (updateSucs, fails []map[string]interface{}, userList []*org.User) {
	fail := make([]map[string]interface{}, 0)
	updateSuc := make([]map[string]interface{}, 0)
	users := make([]*org.User, 0)
	//获取过滤字段
	//_, filters := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldAdminStatus, consts.AllAttr)
	for k := range list {
		tx := u.DB.Begin()
		depID := list[k][consts.DEPID].(string)
		//delete(list[k], consts.DEPNAME)
		delete(list[k], consts.DEPID)
		u2 := &org.User{}
		marshal, err := json.Marshal(list[k])
		if err != nil {
			fail = append(fail, list[k])
			continue
		}
		err = json.Unmarshal(marshal, u2)
		if err != nil {
			fail = append(fail, list[k])
			continue
		}
		err = u.userRepo.UpdateByID(ctx, tx, u2)
		if err != nil {
			fail = append(fail, list[k])
			tx.Rollback()
			return
		}
		//Filter(list[k], filters, IN)
		id := list[k][consts.ID].(string)
		err = u.updateUserDepRelation(tx, list[k][consts.ID].(string), depID)
		if err != nil {
			list[k][consts.REMARK] = consts.RelationDepartmentFail
			delete(list[k], consts.ID)
			fail = append(fail, list[k])
			tx.Rollback()
			continue
		}
		updateSuc = append(updateSuc, list[k])
		tx.Commit()
		us := u.userRepo.Get(ctx, u.DB, id)
		users = append(users, us)
	}
	return updateSuc, fail, users
}

func (u *user) updateUserDepRelation(tx *gorm.DB, userID, depID string) error {
	err := u.userDepRepo.DeleteByUserIDs(tx, userID)
	if err != nil {
		return err
	}
	relation := org.UserDepartmentRelation{
		ID:     id2.ShortID(0),
		UserID: userID,
		DepID:  depID,
	}
	err = u.userDepRepo.Add(tx, &relation)
	if err != nil {
		return err
	}
	return nil
}
