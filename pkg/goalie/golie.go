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
	"bytes"
	"encoding/json"
	"errors"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"context"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	host = "http://goalie"
	//host = "http://127.0.0.1:8081"

	getLoginUserRolesURI = "/api/v1/goalie/role/now/list"
)

// Goalie interface api
type Goalie interface {
	GetLoginUserRoles(ctx context.Context, r *http.Request) (*GetLoginUserRolesResponse, error)
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

// GetLoginUserRolesRequest get login user roles
type GetLoginUserRolesRequest struct {
	UserID       string   `json:"userID"`
	DepartmentID []string `json:"departmentID"`
}

// GetLoginUserRolesResponse get login user roles response
type GetLoginUserRolesResponse struct {
	Roles []*Role `json:"roles"`
	Total int64   `json:"total"`
}

// Role role
type Role struct {
	ID     string `json:"id"`
	RoleID string `json:"roleID"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
}

// GetLoginUserRoles get user roles
func (u *goalie) GetLoginUserRoles(ctx context.Context, r *http.Request) (*GetLoginUserRolesResponse, error) {
	resp := &GetLoginUserRolesResponse{}
	request := r.Clone(r.Context())
	parse, _ := url.ParseRequestURI(host)
	request.URL = parse
	request.Host = host
	request.URL.Path = getLoginUserRolesURI
	request.RequestURI = ""
	request.Method = http.MethodGet
	request.URL.RawQuery = r.URL.RawQuery

	body, err := json.Marshal(&GetLoginUserRolesRequest{})
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		itoa := strconv.Itoa(len(body))
		request.Header.Set("Content-Length", itoa)
		request.ContentLength = int64(len(body))
		request.Body = io.NopCloser(bytes.NewReader(body))
	}
	response, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}
	deserializationResp, err := DeserializationResp(ctx, response, resp)
	if err != nil {
		return nil, err
	}
	if deserializationResp.Code != 0 {
		return nil, errors.New("request goalie get login user roles err, " + deserializationResp.Msg)
	}
	return resp, nil

}

type respData interface{}

// R response data
type R struct {
	err  error
	Code int64    `json:"code"`
	Msg  string   `json:"msg,omitempty"`
	Data respData `json:"data"`
}

// DeserializationResp marshal response body to struct
func DeserializationResp(ctx context.Context, response *http.Response, entity interface{}) (*R, error) {
	if response.StatusCode != http.StatusOK {
		return nil, error2.New(error2.Internal)
	}
	r := new(R)
	r.Data = entity
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	return r, nil
}
