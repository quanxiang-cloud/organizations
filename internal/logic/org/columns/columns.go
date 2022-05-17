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
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

// Columns interface
type Columns interface {
	GetAll(ctx context.Context, r *GetAllColumnsRequest) (*GetAllColumnsResponse, error)
	GetByRoleID(ctx context.Context, r *GetColumnsByRoleRequest) (*GetColumnsByRoleResponse, error)
	Update(ctx context.Context, r *UpdateColumnRequest) (*UpdateColumnResponse, error)
	Set(ctx context.Context, r *SetUseColumnsRequest) (*SetUseColumnsResponse, error)
	Open(ctx context.Context, r *OpenColumnRequest) (*OpenColumnResponse, error)
}

const (
	useStatus   = 1
	unUseStatus = -1
)

// Columns column
type columns struct {
	DB               *gorm.DB
	useColumnsRepo   org.UseColumnsRepo
	tableColumnsRepo org.UserTableColumnsRepo
	redisClient      redis.UniversalClient
	userRepo         org.UserRepo
	conf             configs.Config
}

// NewColumns new
func NewColumns(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) Columns {
	return &columns{
		DB:               db,
		redisClient:      redisClient,
		useColumnsRepo:   mysql2.NewUseColumnsRepo(),
		tableColumnsRepo: mysql2.NewUserTableColumnsRepo(),
		userRepo:         mysql2.NewUserRepo(),
		conf:             conf,
	}
}

// UpdateColumnRequest update column name
type UpdateColumnRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Format    string `json:"format"`
	UpdatedBy string
}

// UpdateColumnResponse update column name
type UpdateColumnResponse struct {
}

// Update update columns alias name
func (c *columns) Update(ctx context.Context, r *UpdateColumnRequest) (*UpdateColumnResponse, error) {
	getByName := c.tableColumnsRepo.GetByName(ctx, c.DB, r.Name)
	if getByName != nil && getByName.ID != r.ID {
		return nil, error2.New(code.ErrColumnExist)
	}

	res := c.tableColumnsRepo.SelectByID(ctx, c.DB, r.ID)
	if res == nil {
		return nil, error2.New(code.DataNotExist)
	}
	tableColumns := org.UserTableColumns{}
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

// GetAllColumnsRequest get column request
type GetAllColumnsRequest struct {
	Status int    `json:"status" form:"status"`
	Name   string `json:"name" form:"name"`
}

// GetAllColumnsResponse get column response
type GetAllColumnsResponse struct {
	All []ColumnResponse `json:"all"`
}

// ColumnResponse response struct
type ColumnResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ColumnName string `json:"columnName"`
	Types      string `json:"types"`
	Len        int    `json:"len"`
	PointLen   int    `json:"pointLen"`
	//1:use,-1:no use
	Status int `json:"status,omitempty"`
	//1:sys,2:alias
	Attr   int    `json:"attr"`
	Format string `json:"format,omitempty"`
	Flag   int    `json:"flag,omitempty"`
}

// GetAll get all column
func (c *columns) GetAll(ctx context.Context, r *GetAllColumnsRequest) (*GetAllColumnsResponse, error) {
	tableColumns, _ := c.tableColumnsRepo.GetAll(ctx, c.DB, r.Status, r.Name)

	all := &GetAllColumnsResponse{}
	for k := range tableColumns {
		if tableColumns[k].ColumnsName == consts.TENANTID {
			continue
		}
		response := ColumnResponse{}
		response.ID = tableColumns[k].ID
		response.Name = tableColumns[k].Name
		response.ColumnName = tableColumns[k].ColumnsName
		response.Len = tableColumns[k].Len
		response.PointLen = tableColumns[k].PointLen
		response.Types = tableColumns[k].Types
		response.Status = unUseStatus
		response.Attr = tableColumns[k].Attr
		response.Format = tableColumns[k].Format
		all.All = append(all.All, response)
	}
	return all, nil
}

// SetUseColumnsRequest req set use columns
type SetUseColumnsRequest struct {
	RoleID string   `json:"roleID"`
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

// SetUseColumnsResponse req set use columns
type SetUseColumnsResponse struct {
}

// Set set use column
func (c *columns) Set(ctx context.Context, r *SetUseColumnsRequest) (*SetUseColumnsResponse, error) {
	all, _ := c.tableColumnsRepo.GetAll(ctx, c.DB, 0, "")
	columnMap := make(map[string]*org.UserTableColumns)
	for k := range all {
		columnMap[all[k].ID] = &all[k]
	}
	useAll := c.useColumnsRepo.SelectAll(ctx, c.DB, r.RoleID)
	useMap := make(map[string]*org.UseColumns)
	for k := range useAll {
		useMap[useAll[k].ColumnID] = &useAll[k]
	}
	tx := c.DB.Begin()
	adds := make([]org.UseColumns, 0)
	if len(r.Add) > 0 {
		add := make([]string, 0)
		for k := range r.Add {
			if v1, ok := columnMap[r.Add[k]]; ok && v1 != nil {
				add = append(add, r.Add[k])
			}
		}
		for k := range add {
			if v1, ok := useMap[add[k]]; !ok || v1 == nil {
				useColumn := org.UseColumns{}
				useColumn.ID = id.HexUUID(true)
				useColumn.ColumnID = add[k]
				useColumn.RoleID = r.RoleID
				adds = append(adds, useColumn)
			}
		}
	}
	if len(r.Delete) > 0 {
		del := make([]string, 0)
		for k := range r.Delete {
			if v1, ok := useMap[r.Delete[k]]; ok && v1 != nil {
				del = append(del, useMap[r.Delete[k]].ID)
			}
		}
		if len(del) > 0 {
			err := c.useColumnsRepo.DeleteByID(ctx, tx, del...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

	}
	if len(adds) > 0 {
		err := c.useColumnsRepo.Create(ctx, tx, adds)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
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
	_, total := c.tableColumnsRepo.GetAll(ctx, c.DB, 0, "")
	if total > 0 {
		return nil, error2.New(code.ErrFieldColumnUsed)
	}
	userColumns := c.userRepo.GetColumns(ctx, c.DB, new(org.User), c.conf.Mysql.DB)
	tx := c.DB.Begin()

	useColumns := make([]org.UseColumns, 0)
	for k := range userColumns {
		if userColumns[k].GetColumnName() != "" && userColumns[k].GetDataType() != "" {
			tableColumns := org.UserTableColumns{}
			tableColumns.ID = id.HexUUID(true)
			tableColumns.Name = userColumns[k].GetName()
			tableColumns.ColumnsName = userColumns[k].GetColumnName()
			tableColumns.Types = consts.FrontColumns[userColumns[k].GetDataType()]
			tableColumns.Len = userColumns[k].GetCharacterMaximumLength()
			tableColumns.PointLen = userColumns[k].GetNumericScale()
			tableColumns.Attr = consts.SystemAttr
			tableColumns.Status = consts.NormalStatus

			err := c.tableColumnsRepo.Insert(ctx, tx, &tableColumns)

			useColumn := org.UseColumns{}
			useColumn.ID = id.HexUUID(true)
			useColumn.ColumnID = tableColumns.ID
			useColumn.RoleID = "1"

			useColumns = append(useColumns, useColumn)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			continue
		}
		tx.Rollback()
		return nil, errors.New("columns field value err")
	}
	err := c.useColumnsRepo.Create(ctx, tx, useColumns)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return nil, nil
}

// GetColumnsByRoleRequest get column request by role
type GetColumnsByRoleRequest struct {
	RoleID string `json:"roleID" form:"roleID"`
}

// GetColumnsByRoleResponse get column response by role
type GetColumnsByRoleResponse struct {
	All []ColumnResponse `json:"all"`
}

// about web check box select
const (
	Selected = iota + 1
	UnSelected
)

// GetByRoleID get all column by role
func (c *columns) GetByRoleID(ctx context.Context, r *GetColumnsByRoleRequest) (*GetColumnsByRoleResponse, error) {
	tableColumns, _ := c.tableColumnsRepo.GetAll(ctx, c.DB, consts.NormalStatus, "")
	useColumns := c.useColumnsRepo.SelectAll(ctx, c.DB, r.RoleID)
	columIDMap := make(map[string]string)
	for k := range useColumns {
		columIDMap[useColumns[k].ColumnID] = useColumns[k].ColumnID
	}
	all := &GetColumnsByRoleResponse{}
	for k := range tableColumns {
		if tableColumns[k].ColumnsName == consts.TENANTID {
			continue
		}
		response := ColumnResponse{}
		response.ID = tableColumns[k].ID
		response.Name = tableColumns[k].Name
		response.ColumnName = tableColumns[k].ColumnsName
		response.Len = tableColumns[k].Len
		response.PointLen = tableColumns[k].PointLen
		response.Types = tableColumns[k].Types
		response.Status = unUseStatus
		response.Attr = tableColumns[k].Attr
		response.Format = tableColumns[k].Format
		response.Flag = UnSelected
		all.All = append(all.All, response)
	}
	if len(all.All) > 0 {
		for k := range all.All {
			if v1, ok := columIDMap[all.All[k].ID]; ok && v1 != "" {
				all.All[k].Flag = Selected
			}
		}
	}
	return all, nil
}
