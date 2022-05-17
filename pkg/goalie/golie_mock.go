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
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type goalieMock struct {
	mock.Mock
}

// NewGoalieMock mock
func NewGoalieMock() Goalie {
	// create an instance of our test object
	user := new(goalieMock)

	// setup expectations
	user.On("GetLoginUserRoles", mock.Anything).Return(map[string]*GetLoginUserRolesResponse{
		"-1": {
			Total: 1,
			Roles: []*Role{
				{
					ID:     "1",
					RoleID: "1",
					Name:   "roleName",
					Tag:    "tag",
				},
			},
		},
	}, nil)
	return user
}

func (m *goalieMock) GetLoginUserRoles(ctx context.Context, r *http.Request) (*GetLoginUserRolesResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*GetLoginUserRolesResponse)
	resp := &GetLoginUserRolesResponse{}
	if r != nil {
		data, _ := res["-1"]
		return data, nil
	}

	return resp, args.Error(1)
}
