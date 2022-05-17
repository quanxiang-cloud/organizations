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
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/user"
	"github.com/quanxiang-cloud/organizations/pkg/goalie"
	"net/http"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/logic/octopus/core"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	oct "github.com/quanxiang-cloud/organizations/internal/models/octopus"
	mysql3 "github.com/quanxiang-cloud/organizations/internal/models/octopus/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

// OthServer other server interface
type OthServer interface {
	AddUsers(c context.Context, rq *AddRequest) (res *AddListResponse, err error)
}

// othersServer other struct
type othersServer struct {
	DB             *gorm.DB
	columnRepo     oct.UserTableColumnsRepo
	conf           configs.Config
	client         http.Client
	extend         oct.ExtendRepo
	columnRoleRepo oct.UseColumnsRepo
	goalieClient   goalie.Goalie
}

// NewOtherServer new
func NewOtherServer(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) OthServer {
	return &othersServer{

		DB:             db,
		columnRepo:     mysql3.NewUserTableColumnsRepo(),
		conf:           conf,
		client:         client.New(conf.InternalNet),
		extend:         mysql3.NewExtendRepo(),
		columnRoleRepo: mysql3.NewUseColumnsRepo(),
		goalieClient:   goalie.NewGoalie(conf.InternalNet),
	}
}

// AddRequest other server add user and department to org request
type AddRequest struct {
	Users      []map[string]interface{} `json:"users"`
	IsUpdate   int                      `json:"isUpdate"`   //是否更新已有数据，1更新，-1不更新
	SyncID     string                   `json:"syncID"`     //同步中心id
	SyncSource string                   `json:"syncSource"` //同步来源
	R          *http.Request
	W          http.ResponseWriter
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

// AddUsers  other server add user to org
func (u *othersServer) AddUsers(c context.Context, r *AddRequest) (res *AddListResponse, err error) {
	result, err := u.addUser(c, r.Users, r.IsUpdate, r.R, r.W)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *othersServer) addUser(c context.Context, reqData []map[string]interface{}, isUpdate int, r *http.Request, w http.ResponseWriter) (*AddListResponse, error) {
	response, err := core.DealRequest(u.client, u.conf.OrgHost, r, reqData)
	if err != nil {
		return nil, err
	}
	result := new(AddListResponse)
	resp, err := core.DeserializationResp(c, response, result)
	if err != nil {
		core.DealResponse(w, response)
		return nil, err
	}
	if resp != nil && resp.Code == 0 {
		_, tenantID := ginheader.GetTenantID(c).Wreck()
		columnIDs, _, err := user.GetRoles(c, u.DB, r, u.columnRoleRepo, u.goalieClient)
		if err != nil {
			return nil, err
		}
		_, aliasFilter := u.columnRepo.GetFilter(c, u.DB, consts.AliasAttr, false, columnIDs...)
		if aliasFilter != nil {
			for k := range reqData {
				if result.Result[k].ID != "" {
					data := reqData[k]
					core.Filter(&data, aliasFilter, core.IN)
					data[consts.ID] = result.Result[k].ID
					if result.Result[k].Attr == 12 {
						extend := &oct.Extend{}
						extend.ID = result.Result[k].ID
						err = u.extend.UpdateByID(u.DB, u.DB, tenantID, extend, data)
					} else {
						err = u.extend.Insert(u.DB, u.DB, tenantID, data)
					}

					if err != nil {
						result.Result[k].Attr = 0
						return result, err
					}
					return result, nil
				}

			}

		}
	}
	core.DealResponse(w, response)
	return nil, nil
}
