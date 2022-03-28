package landlord

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

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	host = "http://landlord/api/v1/landlord"

	registerURI       = "/m/tenant/add"
	cancelRelationURI = "/m/tenant/cancel/relation"
	setRelationURI    = "/m/tenant/set/relation"
)

// Landlord interface
type Landlord interface {
	Register(ctx context.Context, header http.Header, r *RegisterRequest) (*RegisterResponse, error)
	CancelRelation(ctx context.Context, header http.Header, r *CancelRelationRequest) (*CancelRelationResponse, error)
	SetRelation(ctx context.Context, header http.Header, r *SetRelationRequest) (*SetRelationResponse, error)
}
type landlord struct {
	client http.Client
}

// NewLandlord 初始化对象
func NewLandlord(conf client.Config) Landlord {
	return &landlord{
		client: client.New(conf),
	}
}

// RegisterRequest other server add landlord request
type RegisterRequest struct {
	Name    string `json:"name" `
	OwnerID string `json:"ownerID" `
}

// RegisterResponse other server add landlord or dep to org response
type RegisterResponse struct {
	ID string `json:"id"`
}

// Register 实际请求
func (u *landlord) Register(ctx context.Context, header http.Header, r *RegisterRequest) (*RegisterResponse, error) {
	response := &RegisterResponse{}
	err := POST(ctx, &u.client, header, host+registerURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// CancelRelationRequest cancel relation
type CancelRelationRequest struct {
	UserID   []string `json:"userID" binding:"required"`
	TenantID string   `json:"tenantID"`
}

// CancelRelationResponse cancel relation
type CancelRelationResponse struct {
}

// CancelRelation 实际请求
func (u *landlord) CancelRelation(ctx context.Context, header http.Header, r *CancelRelationRequest) (*CancelRelationResponse, error) {
	response := &CancelRelationResponse{}
	err := POST(ctx, &u.client, header, host+cancelRelationURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// SetRelationRequest set relation
type SetRelationRequest struct {
	UserID   []string `json:"userID" binding:"required"`
	TenantID string   `json:"tenantID"`
}

// SetRelationResponse set relation
type SetRelationResponse struct {
}

// SetRelation 实际请求
func (u *landlord) SetRelation(ctx context.Context, header http.Header, r *SetRelationRequest) (*SetRelationResponse, error) {
	response := &SetRelationResponse{}
	err := POST(ctx, &u.client, header, host+setRelationURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GET http get
func GET(ctx context.Context, client *http.Client, headers http.Header, uri string, params map[string]string, entity interface{}) error {

	if len(params) > 0 {
		uri = uri + "?"
		for k, v := range params {
			uri = uri + k + "=" + v + "&"
		}
		uri = uri[:len(uri)-1]
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	if headers != nil {
		req.Header = headers
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = DeserializationResp(ctx, response, entity)
	if err != nil {
		return err
	}
	return err
}

// POST http post
func POST(ctx context.Context, client *http.Client, header http.Header, uri string, params interface{}, entity interface{}) error {
	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return err
	}
	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = DeserializationResp(ctx, response, entity)
	if err != nil {
		return err
	}
	return err
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
func DeserializationResp(ctx context.Context, response *http.Response, entity interface{}) error {
	if response.StatusCode != http.StatusOK {
		return error2.New(error2.Internal)
	}

	r := new(R)
	r.Data = entity

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	if r.Code != error2.Success {
		return error2.NewErrorWithString(r.Code, r.Msg)
	}

	return nil
}
