package header2

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

	"github.com/gin-gonic/gin"
)

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
	_tenantID     = "Tenant-Id"
)

// Profile proile
type Profile struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	DepartmentID string `json:"department_id"`
	TenantID     string `json:"tenant_id"`
}

// GetProfile get info from header
func GetProfile(c *gin.Context) Profile {
	return getProfileFromGIN(c)
}

func getProfileFromGIN(c *gin.Context) Profile {
	userID := c.GetHeader(_userID)
	userName := c.GetHeader(_userName)
	departmentID := c.GetHeader(_departmentID)
	tenantID := c.GetHeader(_tenantID)

	return Profile{
		UserID:       userID,
		UserName:     userName,
		DepartmentID: strings.Split(departmentID, ",")[0],
		TenantID:     tenantID,
	}
}

// GetDepartments get departments
func GetDepartments(c *gin.Context) [][]string {
	departmentID := c.GetHeader(_departmentID)
	res := make([][]string, 0)
	if departmentID == "" {
		return nil
	}
	split := strings.Split(departmentID, "|")
	for k := range split {
		res = append(res, strings.Split(split[k], ","))
	}
	return res
}

// CloneProfile clone header
func CloneProfile(dst *http.Header, src http.Header) {
	dst.Set(_userID, deepCopy(src.Values(_userID)))
	dst.Set(_userName, deepCopy(src.Values(_userName)))
	dst.Set(_departmentID, deepCopy(src.Values(_departmentID)))
	dst.Set(_tenantID, deepCopy(src.Values(_tenantID)))
}

func deepCopy(src []string) string {
	for _, elem := range src {
		if elem != "" {
			return elem
		}
	}
	return ""
}

const (
	roleName = "Role"
	roleID   = "Role-Id"
)

// SetRole header set role info
func SetRole(c *gin.Context, role ...string) {
	roles := strings.Join(role, ",")
	c.Request.Header.Set(roleName, roles)
	c.Writer.Header().Set(roleName, roles)
}

// SetRoleID set role id
func SetRoleID(c *gin.Context, roleIDs ...string) {
	roles := strings.Join(roleIDs, ",")
	c.Request.Header.Set(roleID, roles)
	c.Writer.Header().Set(roleID, roles)
}

// Role role
type Role struct {
	Role   []string
	RoleID []string
}

// IsSuper role is super
func (r *Role) IsSuper() bool {
	for _, role := range r.Role {
		if role == "super" {
			return true
		}
	}
	return false
}

// GetRole get role
func GetRole(c *gin.Context) *Role {
	roleStr := c.Request.Header.Get(roleName)
	roles := strings.Split(roleStr, ",")
	return &Role{Role: roles}
}

// GetRoleID get role id
func GetRoleID(c *gin.Context) *Role {
	roleStr := c.Request.Header.Get(roleID)
	roles := strings.Split(roleStr, ",")
	return &Role{RoleID: roles}
}

// SetContext set context
func SetContext(ctx context.Context, name, value interface{}) context.Context {
	return context.WithValue(ctx, name, value)
}
