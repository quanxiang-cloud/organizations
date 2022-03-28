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
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/core"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	oct "github.com/quanxiang-cloud/organizations/internal/models/octopus"
	mysql3 "github.com/quanxiang-cloud/organizations/internal/models/octopus/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	ldap "github.com/quanxiang-cloud/organizations/pkg/ladp"
	"github.com/quanxiang-cloud/organizations/pkg/message"
)

// User interface
type User interface {
	Add(ctx context.Context, req *AddUserRequest, r *http.Request) (*AddUserResponse, error)
	Update(ctx context.Context, req *UpdateUserRequest, r *http.Request) (*UpdateUserResponse, error)
	AdminSelectByID(c context.Context, req *SearchOneUserRequest, r *http.Request) (*SearchOneUserResponse, error)
	UserSelectByID(c context.Context, req *ViewerSearchOneUserRequest, r *http.Request) (*ViewerSearchOneUserResponse, error)
	Template(c context.Context, req *GetTemplateFileRequest, r *http.Request) (*GetTemplateFileResponse, error)
}

const (
	ownerDepName = "所在部门名称"
)

// User user
type user struct {
	DB *gorm.DB
	//message     message.Message
	redisClient redis.UniversalClient
	columnRepo  oct.UserTableColumnsRepo
	message     message.Message
	ldap        ldap.Ldap
	conf        configs.Config
	extend      oct.ExtendRepo
	client      http.Client
}

// NewUser new
func NewUser(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) User {
	return &user{
		DB:          db,
		redisClient: redisClient,
		columnRepo:  mysql3.NewUserTableColumnsRepo(),
		message:     message.NewMessage(conf.InternalNet),
		ldap:        ldap.NewLdap(conf.InternalNet),
		conf:        conf,
		extend:      mysql3.NewExtendRepo(),
		client:      client.New(conf.InternalNet),
	}
}

// AddUserRequest add user request
type AddUserRequest struct {
	UserInfo map[string]interface{} `json:"userInfo"`
	DepIDs   []string               `json:"depIDs"`
	//1:normal,2:will be active
	Status int `json:"status"`
}

// AddUserResponse 管理员可见字段
type AddUserResponse struct {
	Response *http.Response
}

// Add  add user
func (u *user) Add(ctx context.Context, req *AddUserRequest, r *http.Request) (*AddUserResponse, error) {
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, req.UserInfo)
	if err != nil {
		return nil, err
	}
	in := new(core.INResponse)
	resp, err := core.DeserializationResp(ctx, response, in)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if resp != nil && resp.Code == 0 {
		tx := u.DB.Begin()
		_, aliasFilter := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldAdminStatus, consts.AliasAttr)
		if aliasFilter != nil {
			core.Filter(&req.UserInfo, aliasFilter, core.IN)
			req.UserInfo[consts.ID] = in.ID
			err = u.extend.Insert(u.DB, tx, tenantID, req.UserInfo)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
		}
	}
	return &AddUserResponse{Response: response}, nil

}

// UpdateUserRequest update user
type UpdateUserRequest struct {
	UserInfo map[string]interface{} `json:"userInfo"`
	DepIDs   []string               `json:"depIDs"`
}

// UpdateUserResponse update user response
type UpdateUserResponse struct {
	Response *http.Response
}

// Update update
func (u *user) Update(ctx context.Context, req *UpdateUserRequest, r *http.Request) (*UpdateUserResponse, error) {
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, req.UserInfo)
	if err != nil {
		return nil, err
	}
	in := new(core.INResponse)
	resp, err := core.DeserializationResp(ctx, response, in)
	if resp != nil && resp.Code == 0 {
		_, tenantID := ginheader.GetTenantID(ctx).Wreck()
		tx := u.DB.Begin()
		_, aliasFilter := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldAdminStatus, consts.AliasAttr)
		if aliasFilter != nil {
			core.Filter(&req.UserInfo, aliasFilter, core.IN)
			extend := &oct.Extend{}
			extend.ID = in.ID
			err = u.extend.UpdateByID(u.DB, tx, tenantID, extend, req.UserInfo)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
		}
	}
	return &UpdateUserResponse{
		Response: response,
	}, nil
}

// ManageDepartmentResponse admin can be view
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
	CreatedBy string `json:"createdBy,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

// DepOneResponse department response
type DepOneResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	LeaderID  string `json:"leaderID"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID,omitempty"`
	CompanyID string `json:"companyID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	//1:company,2:department
	Attr int `json:"attr"`
}

//ImportFileRequest import file
type ImportFileRequest struct {
	//1:normal，-2:invalid，-1:del，2:active
	UseStatus int `json:"useStatus" form:"useStatus" binding:"required,max=64"`
	//1:update,-1:not update old data
	IsUpdate int `json:"isUpdate" form:"isUpdate" `
}

// ImportFileResponse import file response data
type ImportFileResponse struct {
	AddSuccessTotal    int                      `json:"addSuccessTotal"`
	AddData            []map[string]interface{} `json:"addData"`
	UpdateSuccessTotal int                      `json:"updateSuccessTotal"`
	UpdateData         []map[string]interface{} `json:"updateData"`
	FailTotal          int                      `json:"failTotal"`
	FailUsers          []map[string]interface{} `json:"failUsers"`
}

// SearchOneUserRequest search one
type SearchOneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// SearchOneUserResponse admin response
type SearchOneUserResponse struct {
	Response *http.Response
}

// AdminSelectByID admin get one by id
func (u *user) AdminSelectByID(ctx context.Context, req *SearchOneUserRequest, r *http.Request) (*SearchOneUserResponse, error) {
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, req)
	if err != nil {

		return nil, err
	}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	old := u.extend.SelectByID(u.DB, tenantID, req.ID)
	if old != nil {
		m := make(map[string]interface{})
		resp, _ := core.DeserializationResp(ctx, response, &m)
		_, filter := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldAdminStatus, consts.AllAttr)
		if filter != nil {
			core.Filter(&old, filter, core.OUT)
			for k, v1 := range old {
				m[k] = v1
			}
			resp.Data = m
			marshal, _ := json.Marshal(resp)

			response.Body = io.NopCloser(bytes.NewReader(marshal))
			l := len(marshal)
			itoa := strconv.Itoa(l)
			response.Header.Set("Content-Length", itoa)
		}
	}
	return &SearchOneUserResponse{
		Response: response,
	}, nil
}

// ViewerSearchOneUserRequest user get one by id
type ViewerSearchOneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// ViewerSearchOneUserResponse user response
type ViewerSearchOneUserResponse struct {
	Response *http.Response
}

// UserSelectByID user select by id
func (u *user) UserSelectByID(ctx context.Context, req *ViewerSearchOneUserRequest, r *http.Request) (*ViewerSearchOneUserResponse, error) {
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, req)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	old := u.extend.SelectByID(u.DB, tenantID, req.ID)
	if old != nil {
		m := make(map[string]interface{})
		resp, _ := core.DeserializationResp(ctx, response, &m)
		_, filter := u.columnRepo.GetFilter(ctx, u.DB, consts.FieldViewerStatus, consts.AliasAttr)
		if filter != nil {
			core.Filter(&old, filter, core.OUT)
			for k, v1 := range old {
				m[k] = v1
			}
			resp.Data = m
			marshal, _ := json.Marshal(resp)

			response.Body = io.NopCloser(bytes.NewReader(marshal))
			l := len(marshal)
			itoa := strconv.Itoa(l)
			response.Header.Set("Content-Length", itoa)
		}
	}
	return &ViewerSearchOneUserResponse{
		Response: response,
	}, nil

}

// GetTemplateFileRequest temp file
type GetTemplateFileRequest struct {
	Octopus int `json:"octopus"`
}

// GetTemplateFileResponse temp file
type GetTemplateFileResponse struct {
	Data        []byte            `json:"data"`
	FileName    string            `json:"fileName"`
	ExcelColumn map[string]string `json:"excelColumn"`
}

// Template get xlsx template
func (u *user) Template(c context.Context, req *GetTemplateFileRequest, r *http.Request) (*GetTemplateFileResponse, error) {
	if r.URL.RawQuery == "" {
		r.URL.RawQuery = "octopus=1"
	} else {
		r.URL.RawQuery = r.URL.RawQuery + "&octopus=1"
	}
	req.Octopus = 1
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, req)
	if err != nil {
		return nil, err
	}
	fileResponse := &GetTemplateFileResponse{}
	_, err = core.DeserializationResp(c, response, fileResponse)
	if err != nil {
		return nil, err
	}
	if len(fileResponse.ExcelColumn) == 0 {
		return nil, error2.New(code.FieldNameIsNull)
	}
	xlsxFields := u.columnRepo.GetXlsxField(c, u.DB, consts.FieldAdminStatus)
	if len(xlsxFields) > 0 {
		for k, v := range fileResponse.ExcelColumn {
			xlsxFields[k] = v
		}
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
	s = append(s, ownerDepName)
	sort.Strings(s)
	for k := range s {
		cell := row.AddCell()
		cell.SetValue(s[k])
	}
	buffer := new(bytes.Buffer)
	newFile.Write(buffer)
	res := &GetTemplateFileResponse{}
	res.ExcelColumn = xlsxFields
	res.Data = buffer.Bytes()
	res.FileName = u.conf.TemplateName
	return res, nil

}
