package ldap

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
	host = "http://ldap/api/v1/ldap"

	addUserURI     = "/user/add"
	updateUserURI  = "/user/update"
	delUserURI     = "/user/del"
	setPasswordURI = "/user/set/password"
	authURI        = "/auth"
)

// UserAddReq user add req
type UserAddReq struct {
	ID        string `json:"id" binding:"-"`
	Mail      string `json:"mail" binding:"required"`
	GroupName string `json:"groupName" binding:"required"`
	//firstName
	FirstName string `json:"firstName" binding:"required"`
	//lastName
	LastName string `json:"lastName" binding:"required"`
	// userName
	UserName string `json:"userName" binding:"required"`
	// type must int number
	UIDNumber                  int    `json:"uidNumber" binding:"required"`
	UserPassword               string `json:"userPassword" binding:"required"`
	BusinessCategory           string `json:"businessCategory"`
	CarLicense                 string `json:"carLicense"`
	DepartmentNumber           string `json:"departmentNumber"`
	Description                string `json:"description"`
	EmployeeNumber             string `json:"employeeNumber"`
	FacsimileTelephoneNumber   string `json:"facsimileTelephoneNumber"`
	Gecos                      string `json:"gecos"`
	HomePhone                  string `json:"homePhone"`
	LabeledURI                 string `json:"labeledURI"`
	Mobile                     string `json:"mobile"`
	L                          string `json:"l"`
	O                          string `json:"o"`
	PhysicalDeliveryOfficeName string `json:"physicalDeliveryOfficeName"`
	PostalAddress              string `json:"postalAddress"`
	PostalCode                 string `json:"postalCode"`
	PostOfficeBox              string `json:"postOfficeBox"`
	RegisteredAddress          string `json:"registeredAddress"`
	RoomNumber                 string `json:"roomNumber"`
	ST                         string `json:"st"`
	Street                     string `json:"street"`
	TelephoneNumber            string `json:"telephoneNumber"`
	Title                      string `json:"title"`
	TenantID                   string `json:-`
	GidNumber                  string `json:"gidNumber"`
}

// UserAddResp user add resp
type UserAddResp struct {
}

// UserUpdateReq user update req
type UserUpdateReq struct {
	ID        string `json:"id" binding:"required"`
	Mail      string `json:"mail" binding:"required"`
	GroupName string `json:"groupName" binding:"-"`
	//firstName
	FirstName string `json:"firstName" binding:"-"`
	//lastName
	LastName string `json:"lastName" binding:"-"`
	// userName
	UserName string `json:"userName" binding:"-"`
	// type must int number
	UIDNumber                  int    `json:"uidNumber" binding:"-"`
	UserPassword               string `json:"userPassword"`
	BusinessCategory           string `json:"carLicense"`
	CarLicense                 string `json:"businessCategory"`
	DepartmentNumber           string `json:"departmentNumber"`
	Description                string `json:"description"`
	EmployeeNumber             string `json:"employeeNumber"`
	FacsimileTelephoneNumber   string `json:"facsimileTelephoneNumber"`
	Gecos                      string `json:"gecos"`
	HomePhone                  string `json:"homePhone"`
	LabeledURI                 string `json:"labeledURI"`
	Mobile                     string `json:"mobile"`
	L                          string `json:"l"`
	O                          string `json:"o"`
	PhysicalDeliveryOfficeName string `json:"physicalDeliveryOfficeName"`
	PostalAddress              string `json:"postalAddress"`
	PostalCode                 string `json:"postalCode"`
	PostOfficeBox              string `json:"postOfficeBox"`
	RegisteredAddress          string `json:"registeredAddress"`
	RoomNumber                 string `json:"roomNumber"`
	ST                         string `json:"st"`
	Street                     string `json:"street"`
	TelephoneNumber            string `json:"telephoneNumber"`
	Title                      string `json:"title"`
	TenantID                   string `json:-`
	GidNumber                  string `json:"gidNumber"`
}

// UserUpdateResp user update resp
type UserUpdateResp struct {
}

// UserDelReq del one user
type UserDelReq struct {
	TenantID  string `json:-`
	ID        string `json:"ID" binding:"required"`
	GroupName string `json:"groupName" binding:"required"`
}

// UserDelResp del one user resp
type UserDelResp struct {
}

// AuthReq auth
type AuthReq struct {
	TenantID string `json:-`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// AuthResp auth resp
type AuthResp struct {
	Flag bool `json:"flag"`
}

// UserUpdatePasswordReq user update password req
type UserUpdatePasswordReq struct {
	ID           string `json:"id" binding:"required"`
	Mail         string `json:"mail" binding:"required"`
	UserPassword string `json:"userPassword"`
	GidNumber    string `json:"gidNumber"` //对应部门数字id
	TenantID     string `json:"tenantID"`
}

// UserUpdatePasswordResp user update password resp
type UserUpdatePasswordResp struct {
}

// Ldap ldap服务
type Ldap interface {
	AddUser(ctx context.Context, header http.Header, r *UserAddReq) (*UserAddResp, error)
	UpdateUser(ctx context.Context, header http.Header, r *UserUpdateReq) (*UserUpdateResp, error)
	DelUser(ctx context.Context, header http.Header, r *UserDelReq) (*UserDelResp, error)
	UpdatePassword(ctx context.Context, header http.Header, r *UserUpdatePasswordReq) (*UserUpdatePasswordResp, error)
	Auth(ctx context.Context, header http.Header, r *AuthReq) (*AuthResp, error)
}
type ldap struct {
	client http.Client
}

// NewLdap 初始化对象
func NewLdap(conf client.Config) Ldap {
	return &ldap{
		client: client.New(conf),
	}
}

// AddUser add
func (l *ldap) AddUser(ctx context.Context, header http.Header, reqs *UserAddReq) (*UserAddResp, error) {
	response := UserAddResp{}
	err := POST(ctx, &l.client, header, host+addUserURI, reqs, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// UpdateUser update
func (l *ldap) UpdateUser(ctx context.Context, header http.Header, reqs *UserUpdateReq) (*UserUpdateResp, error) {
	response := UserUpdateResp{}
	err := POST(ctx, &l.client, header, host+updateUserURI, reqs, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// DelUser del
func (l *ldap) DelUser(ctx context.Context, header http.Header, reqs *UserDelReq) (*UserDelResp, error) {
	response := UserDelResp{}
	err := POST(ctx, &l.client, header, host+delUserURI, reqs, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// UpdatePassword update password
func (l *ldap) UpdatePassword(ctx context.Context, header http.Header, reqs *UserUpdatePasswordReq) (*UserUpdatePasswordResp, error) {
	response := UserUpdatePasswordResp{}
	err := POST(ctx, &l.client, header, host+setPasswordURI, reqs, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// Auth auth
func (l *ldap) Auth(ctx context.Context, header http.Header, reqs *AuthReq) (*AuthResp, error) {
	response := AuthResp{}
	var flag = false
	err := POST(ctx, &l.client, header, host+authURI, reqs, &flag)
	if err != nil {
		return nil, err
	}
	response.Flag = flag
	return &response, err
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
	responseData, err := client.Do(req)
	if err != nil {
		return err
	}
	defer responseData.Body.Close()

	err = DeserializationResp(ctx, responseData, entity)
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
	responseData, err := client.Do(req)
	if err != nil {
		return err
	}
	defer responseData.Body.Close()

	err = DeserializationResp(ctx, responseData, entity)
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
