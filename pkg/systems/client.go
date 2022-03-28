package systems

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
	"encoding/json"
	"io"
	"net/http"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
)

const (
	host = "http://systems/api/v1/systems"

	enterpriseURI = "/enterprise/t/info"
	securityURI   = "/security/t/info"
)

// Systems interface
type Systems interface {
	// GetEnterpriseInfo get enterprise info
	GetEnterpriseInfo(ctx context.Context) (*EnterpriseInfo, error)
	// GetSecurityInfo get security info
	GetSecurityInfo(ctx context.Context) (*SecurityInfo, error)
}

type systems struct {
	client http.Client
}

// NewSystems new
func NewSystems(conf client.Config) Systems {
	return &systems{
		client: client.New(conf),
	}
}

// EnterpriseInfo Tenant info
type EnterpriseInfo struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
}

// SecurityInfo security info
type SecurityInfo struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	EnterpriseID  string `json:"enterpriseId"`
	IPCount       int64  `json:"ipCount"`
	IPCountWait   int64  `json:"ipCountWait"`
	PwdCount      int64  `json:"pwdCount"`
	PwdCountWait  int64  `json:"pwdCountWait"`
	PwdMinLen     int64  `json:"pwdMinLen"`
	PwdExpireDays int64  `json:"pwdExpireDays"`
	PwdNoticeDays int64  `json:"pwdNoticeDays"`
	PwdType       int64  `json:"pwdType"`
	LoginType     int64  `json:"loginType"`
	PwdChange     bool   `json:"pwdChange"`
	M2FA          bool   `json:"M2FA"`
}

// GetEnterpriseInfo  get enterprise info
func (s *systems) GetEnterpriseInfo(ctx context.Context) (*EnterpriseInfo, error) {
	res := new(EnterpriseInfo)
	err := GET(ctx, &s.client, host+enterpriseURI, nil, res)
	return res, err
}

// GetSecurityInfo get security info
func (s *systems) GetSecurityInfo(ctx context.Context) (*SecurityInfo, error) {
	res := new(SecurityInfo)
	err := GET(ctx, &s.client, host+securityURI, nil, res)
	return res, err
}

// GET http get
func GET(ctx context.Context, client *http.Client, uri string, params map[string]string, entity interface{}) error {

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
	req.Header.Set(ginheader.GetTenantID(ctx).Wreck())
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

// R response
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
