package client

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
	"net/http"

	"context"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	host = "http://org/api/v1/org"

	othAddUsersURI  = "/o/user/add"
	othAddDepsURI   = "/o/department/add"
	oneUserURI      = "/o/user/info"
	usersByIDsURI   = "/o/user/ids"
	depByIDsURI     = "/o/dep/ids"
	usersByDepIDURI = "/o/user/dep/id"
	depMaxGradeURI  = "/o/dep/max/grade"
)

// User interface api
type User interface {
	OthAddUsers(ctx context.Context, r *AddUsersRequest) (*AddListResponse, error)
	OthAddDeps(ctx context.Context, r *AddDepartmentRequest) (*AddListResponse, error)
	GetUserInfo(ctx context.Context, r *OneUserRequest) (*OneUserResponse, error)
	GetUserByIDs(ctx context.Context, r *GetUserByIDsRequest) (*GetUserByIDsResponse, error)
	GetDepByIDs(ctx context.Context, r *GetDepByIDsRequest) (*GetDepByIDsResponse, error)
	GetUsersByDepID(ctx context.Context, r *GetUsersByDepIDRequest) (*GetUsersByDepIDResponse, error)
	GetDepMaxGrade(ctx context.Context, r *GetDepMaxGradeRequest) (*GetDepMaxGradeResponse, error)
}
type user struct {
	client http.Client
}

// NewUser new
func NewUser(conf client.Config) User {
	return &user{
		client: client.New(conf),
	}
}

//AddUsersRequest other server add user request
type AddUsersRequest struct {
	Users []AddUser `json:"users"`
	//1:update old data,2:no
	IsUpdate   int    `json:"isUpdate"`
	SyncID     string `json:"syncID"`
	SyncSource string `json:"syncSource"`
}

//AddUser other server add user to org
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
	DepIDs    []string `json:"depIDs,omitempty"`
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

// OthAddUsers add
func (u *user) OthAddUsers(ctx context.Context, r *AddUsersRequest) (*AddListResponse, error) {
	response := &AddListResponse{}
	err := client.POST(ctx, &u.client, host+othAddUsersURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// AddDepartmentRequest other server add  department to org request
type AddDepartmentRequest struct {
	Deps []AddDep `json:"deps"`
	//1: sync department data,-1:no
	SyncDep int `json:"syncDep"`
	//1:update old data,-1,no
	IsUpdate   int    `json:"isUpdate"`
	SyncID     string `json:"syncID"`
	SyncSource string `json:"syncSource"`
}

// AddDep other server add department to org
type AddDep struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	//1:normal,-1:delete,-2:invalid
	UseStatus int `json:"useStatus"`
	//1:company,2;:department
	Attr      int    `json:"attr"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID"`
	CompanyID string `json:"companyID"`
	Grade     int    `json:"grade"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	Remark    string `json:"remark,omitempty"`
}

//OthAddDeps add dep
func (u *user) OthAddDeps(ctx context.Context, r *AddDepartmentRequest) (*AddListResponse, error) {
	response := &AddListResponse{}
	err := client.POST(ctx, &u.client, host+othAddDepsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// OneUserRequest request
type OneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// OneUserResponse response
type OneUserResponse struct {
	ID        string `json:"id,omitempty" `
	Name      string `json:"name,omitempty" `
	Phone     string `json:"phone,omitempty" `
	Email     string `json:"email,omitempty" `
	SelfEmail string `json:"selfEmail,omitempty" `
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `json:"useStatus,omitempty" `
	TenantID  string `json:"tenantID,omitempty" `
	Position  string `json:"position,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	JobNumber string `json:"jobNumber,omitempty" `
	// 0x1111 right first 0:need reset password
	Status int                 `json:"status"`
	Dep    [][]DepOneResponse  `json:"deps,omitempty"`
	Leader [][]OneUserResponse `json:"leaders,omitempty"`
}

// DepOneResponse 用于用户部门层级线索
type DepOneResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	LeaderID  string `json:"leaderID"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

//GetUserInfo get user info
func (u *user) GetUserInfo(ctx context.Context, r *OneUserRequest) (*OneUserResponse, error) {
	response := &OneUserResponse{}
	err := client.POST(ctx, &u.client, host+oneUserURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

//GetUserByIDsRequest get user by ids request
type GetUserByIDsRequest struct {
	IDs []string `json:"ids"`
}

// GetUserByIDsResponse get user by ids response
type GetUserByIDsResponse struct {
	Users []OneUserResponse `json:"users"`
}

//GetUserByIDs get user by ids
func (u *user) GetUserByIDs(ctx context.Context, r *GetUserByIDsRequest) (*GetUserByIDsResponse, error) {
	response := &GetUserByIDsResponse{}
	err := client.POST(ctx, &u.client, host+usersByIDsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetDepByIDsRequest request
type GetDepByIDsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// GetDepByIDsResponse response
type GetDepByIDsResponse struct {
	Deps []DepOneResponse `json:"deps"`
}

//GetDepByIDs get dep by ids
func (u *user) GetDepByIDs(ctx context.Context, r *GetDepByIDsRequest) (*GetDepByIDsResponse, error) {
	response := &GetDepByIDsResponse{}
	err := client.POST(ctx, &u.client, host+depByIDsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetUsersByDepIDRequest get users by department id request
type GetUsersByDepIDRequest struct {
	DepID string `json:"depID"`
	//1:include
	IsIncludeChild int `json:"isIncludeChild"`
}

// GetUsersByDepIDResponse  get users by department id response
type GetUsersByDepIDResponse struct {
	Users []OneUserResponse `json:"users"`
}

// GetUsersByDepID get users by department id
func (u *user) GetUsersByDepID(ctx context.Context, r *GetUsersByDepIDRequest) (*GetUsersByDepIDResponse, error) {
	response := &GetUsersByDepIDResponse{}
	err := client.POST(ctx, &u.client, host+usersByDepIDURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetDepMaxGradeRequest request
type GetDepMaxGradeRequest struct {
}

// GetDepMaxGradeResponse response
type GetDepMaxGradeResponse struct {
	Grade int64 `json:"grade"`
}

//GetDepMaxGrade get dep max grade
func (u *user) GetDepMaxGrade(ctx context.Context, r *GetDepMaxGradeRequest) (*GetDepMaxGradeResponse, error) {
	response := &GetDepMaxGradeResponse{}
	err := client.POST(ctx, &u.client, host+depMaxGradeURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}
