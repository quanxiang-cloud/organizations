package goalie

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
	host = "http://goalie/api/v1/goalie/role"

	delOwner = "/del/owner"
)

// Goalie interface api
type Goalie interface {
	DelOwner(ctx context.Context, r *OthDelRequest) (*OthDelResponse, error)
}
type goalie struct {
	client http.Client
}

// NewGoalie new
func NewGoalie(conf client.Config) Goalie {
	return &goalie{
		client: client.New(conf),
	}
}

// OthDelRequest 其它服务删除权限
type OthDelRequest struct {
	IDs   []string `json:"ids"`
	DelBy string   `json:"delBy"`
}

// OthDelResponse 其它服务删除权限
type OthDelResponse struct {
}

// DelOwner del
func (u *goalie) DelOwner(ctx context.Context, r *OthDelRequest) (*OthDelResponse, error) {
	response := &OthDelResponse{}
	err := client.POST(ctx, &u.client, host+delOwner, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}
