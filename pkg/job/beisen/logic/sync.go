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
	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"gorm.io/gorm"
	"net/http"

	client2 "github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/other"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/header2"
)

// Sync beisein sync job
type Sync interface {
	SyncData(ctx context.Context, req *SyncRequest) (*SyncResponse, error)
}

type sync struct {
	Oth    other.OthServer
	Client http.Client
}

// NewSync new
func NewSync(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) Sync {
	user.NewSearch(db)
	return &sync{
		Oth:    other.NewOtherServer(conf, db, redisClient),
		Client: client2.New(conf.InternalNet),
	}
}

// AddDataV1 add data
type AddDataV1 struct {
	Users      []other.AddUser `json:"users"`
	Deps       []other.AddDep  `json:"deps"`
	SyncDEP    int             `json:"syncDep"`
	IsUpdate   int             `json:"isUpdate"`
	SyncID     string          `json:"syncID"`
	SyncSource string          `json:"syncSource"`
	Profile    header2.Profile
}

// SyncRequest sync
type SyncRequest struct {
	SyncDEP    int
	IsUpdate   int
	RequestURL string
	TenantID   string
}

// SyncResponse sync
type SyncResponse struct {
}

// SyncData sync data
func (s *sync) SyncData(ctx context.Context, req *SyncRequest) (*SyncResponse, error) {
	//1、get data from sync server
	data, err := s.getData(ctx, req)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	au := &other.AddUsersRequest{}
	au.Users = append(au.Users, data.Users...)
	au.SyncDEP = req.SyncDEP
	au.IsUpdate = req.IsUpdate
	if req.SyncDEP == 1 {
		ad := &other.AddDepartmentRequest{}
		ad.Deps = append(ad.Deps, data.Deps...)
		ad.SyncDEP = req.SyncDEP
		ad.IsUpdate = req.IsUpdate
		_, err := s.Oth.AddDepartments(ctx, ad)
		if err != nil {
			logger.Logger.Error(err)
			return nil, err
		}
	}
	_, err = s.Oth.AddUsers(ctx, au)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	sig1 := make(chan int, 1)
	sig2 := make(chan int, 1)
	to := make(chan int)

	go s.Oth.PushUserToSearch(ctx, sig2, to)
	go s.Oth.PushDepToSearch(ctx, sig1)

	var i = 0
	var total = 0
	for {
		if total != 0 && i >= (total) {
			return nil, nil
		}
		select {
		case total = <-to:
		case n := <-sig1:
			i = i + n
		case n := <-sig2:
			i = i + n

		}
	}
}

// getData
func (s *sync) getData(c context.Context, parma *SyncRequest) (*AddDataV1, error) {
	response := &AddDataV1{}
	err := client2.POST(c, &s.Client, parma.RequestURL, &parma, &response)
	if err != nil {
		return nil, err
	}
	return response, err
}
