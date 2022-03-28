package columns

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
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/core"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	oct "github.com/quanxiang-cloud/organizations/internal/models/octopus"
	mysql3 "github.com/quanxiang-cloud/organizations/internal/models/octopus/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

// Columns columns interface
type Columns interface {
	GetAll(ctx context.Context, r *GetAllColumnsRequest) (*GetAllColumnsResponse, error)
	Update(ctx context.Context, r *UpdateColumnRequest) (*UpdateColumnResponse, error)
	Set(ctx context.Context, r *SetUseColumnsRequest) (*SetUseColumnsResponse, error)
	Add(ctx context.Context, r *AddColumnRequest) (*AddColumnResponse, error)
	Drop(ctx context.Context, r *DropColumnRequest) (*DropColumnResponse, error)
	Open(ctx context.Context, r *OpenColumnRequest) (*OpenColumnResponse, error)
}

const (
	useStatus   = 1
	unUseStatus = -1
)

// Columns column
type columns struct {
	DB               *gorm.DB
	manageColumnRepo oct.ManageColumn
	useColumnsRepo   oct.UseColumnsRepo
	tableColumnsRepo oct.UserTableColumnsRepo
	redisClient      redis.UniversalClient
	conf             configs.Config
	client           http.Client
}

// NewColumns new
func NewColumns(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) Columns {
	return &columns{
		DB:               db,
		redisClient:      redisClient,
		manageColumnRepo: mysql3.NewManageColumnRepo(),
		useColumnsRepo:   mysql3.NewUseColumnsRepo(),
		tableColumnsRepo: mysql3.NewUserTableColumnsRepo(),
		conf:             conf,
		client:           client.New(conf.InternalNet),
	}
}

// UpdateColumnRequest update column name
type UpdateColumnRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TenantID  string `json:"tenantID"`
	Format    string `json:"format"`
	UpdatedBy string
	R         *http.Request
	W         http.ResponseWriter
}

// UpdateColumnResponse update column name
type UpdateColumnResponse struct {
}

// Update update columns alias name
func (c *columns) Update(ctx context.Context, r *UpdateColumnRequest) (*UpdateColumnResponse, error) {
	old := c.tableColumnsRepo.SelectByID(ctx, c.DB, r.ID)
	if old == nil {
		response, err := core.DealRequest(c.client, c.conf.OrgHost, r.R, r)
		if err != nil {
			return nil, err
		}
		core.DealResponse(r.W, response)
		return nil, nil
	}
	res := c.tableColumnsRepo.SelectByIDAndName(ctx, c.DB, r.ID, r.Name)
	if res != nil {
		return nil, error2.New(code.ColumnExist)
	}
	tableColumns := oct.UserTableColumns{}
	tableColumns.ID = r.ID
	tableColumns.Name = r.Name
	tableColumns.UpdatedAt = time.NowUnix()
	tableColumns.UpdatedBy = r.UpdatedBy
	tableColumns.Format = r.Format
	tx := c.DB.Begin()
	err := c.tableColumnsRepo.Update(ctx, tx, &tableColumns)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil

}

// AddColumnRequest add column
type AddColumnRequest struct {
	Name       string `json:"name"`
	ColumnName string `json:"columnName"`
	Types      string `json:"types"`
	Len        int    `json:"len"`
	PointLen   int    `json:"pointLen"`
	Attr       int    `json:"attr"`
	Format     string `json:"format"`
	CreatedBy  string
}

// AddColumnResponse add column
type AddColumnResponse struct {
}

// Add add self columns
func (c *columns) Add(ctx context.Context, r *AddColumnRequest) (*AddColumnResponse, error) {
	tableColumns := oct.UserTableColumns{}
	tableColumns.ID = id.HexUUID(true)
	tableColumns.Name = r.Name
	tableColumns.ColumnsName = r.ColumnName
	tableColumns.Types = r.Types
	tableColumns.Len = r.Len
	tableColumns.PointLen = r.PointLen
	tableColumns.Attr = consts.AliasAttr
	unix := time.NowUnix()
	tableColumns.CreatedAt = unix
	tableColumns.UpdatedAt = unix
	tableColumns.CreatedBy = r.CreatedBy
	tableColumns.UpdatedBy = r.CreatedBy
	tableColumns.Status = consts.NormalStatus
	tableColumns.Format = r.Format
	tx := c.DB.Begin()
	err := c.tableColumnsRepo.Insert(ctx, tx, &tableColumns)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	err = c.manageColumnRepo.AddColumns(tx, tenantID, r.ColumnName, r.Types, r.Len, r.PointLen)
	if err != nil {
		tx.Rollback()
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, error2.New(code.ColumnExist)
		}

		return nil, err
	}
	tx.Commit()
	return nil, nil

}

// DropColumnRequest del column
type DropColumnRequest struct {
	ID string `json:"id"`
	R  *http.Request
	W  http.ResponseWriter
}

// DropColumnResponse del column
type DropColumnResponse struct {
	ID string `json:"id"`
}

// Drop del column
func (c *columns) Drop(ctx context.Context, r *DropColumnRequest) (*DropColumnResponse, error) {
	res := c.tableColumnsRepo.SelectByID(ctx, c.DB, r.ID)
	if res == nil {
		response, err := core.DealRequest(c.client, c.conf.OrgHost, r.R, r)
		if err != nil {
			return nil, err
		}
		core.DealResponse(r.W, response)
		return nil, nil
	}
	if res.Attr == consts.SystemAttr {
		return nil, error2.New(code.SystemParameter)
	}
	tx := c.DB.Begin()
	err := c.useColumnsRepo.DeleteByID(ctx, tx, res.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tableColumns := &oct.UserTableColumns{}
	tableColumns.ID = r.ID
	tableColumns.DeletedAt = time.NowUnix()
	tableColumns.Status = consts.DelStatus
	err = c.tableColumnsRepo.Update(ctx, tx, tableColumns)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil

}

// GetAllColumnsRequest user table all column request
type GetAllColumnsRequest struct {
	Status int `json:"status" form:"status"`
	R      *http.Request
	W      http.ResponseWriter
}

// GetAllColumnsResponse user table all column response
type GetAllColumnsResponse struct {
	All []ColumnResponse `json:"all"`
}

// ColumnResponse 用户表字段
type ColumnResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ColumnName string `json:"columnName"`
	Types      string `json:"types"`
	Len        int    `json:"len"`
	PointLen   int    `json:"pointLen"`
	//1:use,-1:no use
	Status int `json:"status"`
	//1:default,-1:can not modify,2:can be modify
	Attr   int    `json:"attr"`
	Format string `json:"format"`
}

// GetAll get all column
func (c *columns) GetAll(ctx context.Context, data *GetAllColumnsRequest) (*GetAllColumnsResponse, error) {
	response, err := core.DealRequest(c.client, c.conf.OrgHost, data.R, data)
	if err != nil {
		return nil, err
	}
	columnsResponse := &GetAllColumnsResponse{}
	_, err = core.DeserializationResp(ctx, response, columnsResponse)
	if err != nil {
		return nil, err
	}
	all := &GetAllColumnsResponse{}
	if len(columnsResponse.All) > 0 {
		all.All = append(all.All, columnsResponse.All...)
	}
	tableColumns, _ := c.tableColumnsRepo.GetAll(ctx, c.DB, data.Status)
	useColumns := c.useColumnsRepo.SelectAll(ctx, c.DB, data.Status)

	for k := range tableColumns {
		if tableColumns[k].ColumnsName == consts.TENANTID {
			continue
		}
		res := ColumnResponse{}
		res.ID = tableColumns[k].ID
		res.Name = tableColumns[k].Name
		res.ColumnName = tableColumns[k].ColumnsName
		res.Len = tableColumns[k].Len
		res.PointLen = tableColumns[k].PointLen
		res.Types = tableColumns[k].Types
		res.Status = unUseStatus
		res.Attr = tableColumns[k].Attr
		res.Format = tableColumns[k].Format
		all.All = append(all.All, res)
	}
	if len(all.All) > 0 {
		for k := range all.All {
			for _, v1 := range useColumns {
				if all.All[k].ID == v1.ColumnID {
					all.All[k].Status = useStatus
				}
			}
		}
	}

	return all, nil
}

// SetUseColumnsRequest req set use columns
type SetUseColumnsRequest struct {
	Columns   []SetUseColumn `json:"columns"`
	CreatedBy string
	R         *http.Request
	W         http.ResponseWriter
}

// SetUseColumn set will use column
type SetUseColumn struct {
	ColumnID     string `json:"columnID"`
	ViewerStatus int    `json:"viewerStatus"`
}

// SetUseColumnsResponse req set use columns
type SetUseColumnsResponse struct {
}

// Set set use column
func (c *columns) Set(ctx context.Context, r *SetUseColumnsRequest) (*SetUseColumnsResponse, error) {
	all, _ := c.tableColumnsRepo.GetAll(ctx, c.DB, 0)
	aliasColumnMap := make(map[string]*oct.UserTableColumns)
	for k := range all {
		aliasColumnMap[all[k].ID] = &all[k]
	}
	systemsColumns := make([]SetUseColumn, 0)
	aliasClolumns := make([]SetUseColumn, 0)
	for k := range r.Columns {
		if v, ok := aliasColumnMap[r.Columns[k].ColumnID]; ok && v != nil {
			aliasClolumns = append(aliasClolumns, r.Columns[k])
		} else {
			systemsColumns = append(systemsColumns, r.Columns[k])
		}
	}
	var flag = false
	if len(systemsColumns) > 0 {
		request := SetUseColumnsRequest{}
		request.Columns = systemsColumns
		response, err := core.DealRequest(c.client, c.conf.OrgHost, r.R, request)
		if err != nil {
			return nil, err
		}
		in := new(core.INResponse)
		resp, err := core.DeserializationResp(ctx, response, in)
		if err != nil {
			return nil, err
		}
		if resp != nil && resp.Code == 0 {
			flag = true
		} else {
			core.DealResponse(r.W, response)
			return nil, nil
		}
	} else {
		flag = true
	}
	if flag {
		tx := c.DB.Begin()
		if len(aliasClolumns) > 0 {
			unix := time.NowUnix()
			useColumns := make([]oct.UseColumns, 0)
			for _, v := range aliasClolumns {
				useColumn := oct.UseColumns{}
				useColumn.ID = id.HexUUID(true)
				useColumn.ColumnID = v.ColumnID
				useColumn.ViewerStatus = v.ViewerStatus
				useColumn.UpdatedBy = r.CreatedBy
				useColumn.UpdatedAt = unix
				useColumn.CreatedBy = r.CreatedBy
				useColumn.CreatedAt = unix
				useColumns = append(useColumns, useColumn)
			}

			err := c.useColumnsRepo.Update(ctx, tx, useColumns)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
			return nil, nil
		}
		err := c.useColumnsRepo.Update(ctx, tx, nil)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return nil, nil
	}
	return nil, nil

}

// OpenColumnRequest tenant alias field
type OpenColumnRequest struct {
}

// OpenColumnResponse resp
type OpenColumnResponse struct {
}

// Open open colum field
func (c *columns) Open(ctx context.Context, r *OpenColumnRequest) (*OpenColumnResponse, error) {
	_, total := c.tableColumnsRepo.GetAll(ctx, c.DB, 0)
	if total > 0 {
		return nil, error2.New(code.ErrFieldColumnUsed)
	}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if c.manageColumnRepo.CheckTableExist(c.DB, tenantID) {
		return nil, error2.New(code.ErrFieldColumnUsed)
	}
	tx := c.DB.Begin()
	err := c.manageColumnRepo.CreateTable(tx, tenantID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return nil, nil
}
