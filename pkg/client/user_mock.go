package client

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

	"github.com/stretchr/testify/mock"
)

type userMock struct {
	mock.Mock
}

// NewUserMock mock
func NewUserMock() User {
	// create an instance of our test object
	user := new(userMock)

	// setup expectations
	user.On("OthAddUsers", mock.Anything).Return(map[string]*AddListResponse{
		"-1": {
			Result: map[int]*Result{
				0: &Result{
					ID:     "1",
					Remark: "ok",
					Attr:   11,
				},
			},
		},
	}, nil)
	// setup expectations
	user.On("OthAddDeps", mock.Anything).Return(map[string]*AddListResponse{
		"-1": {Result: map[int]*Result{
			0: &Result{
				ID:     "1",
				Remark: "ok",
				Attr:   11,
			},
		},
		},
	}, nil)

	// setup expectations
	user.On("GetUserInfo", mock.Anything).Return(map[string]*OneUserResponse{
		"-1": {
			ID:        "xxx",
			Name:      "xxx",
			UseStatus: 1,
		},
	}, nil)
	// setup expectations
	user.On("GetUserByIDs", mock.Anything).Return(map[string]*GetUserByIDsResponse{
		"-1": {
			Users: []OneUserResponse{{
				ID:   "xxx",
				Name: "xxx",
			},
			},
		},
	}, nil)

	// setup expectations
	user.On("GetUsersByDepID", mock.Anything).Return(map[string]*GetUsersByDepIDResponse{
		"-1": {
			Users: []OneUserResponse{{
				ID:   "xxx",
				Name: "xxx",
			},
			},
		},
	}, nil)
	// setup expectations
	user.On("GetDepByIDs", mock.Anything).Return(map[string]*GetDepByIDsResponse{
		"-1": {
			Deps: []DepOneResponse{
				{
					ID:   "xxx",
					Name: "xxx",
				},
			},
		},
	}, nil)

	// setup expectations
	user.On("GetDepMaxGrade", mock.Anything).Return(map[string]*GetDepMaxGradeResponse{
		"-1": {
			Grade: 4,
		},
	}, nil)
	return user
}

func (m *userMock) OthAddUsers(ctx context.Context, reqs *AddUsersRequest) (*AddListResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*AddListResponse)
	resp := &AddListResponse{}
	if len(reqs.Users) > 0 {
		data, _ := res["-1"]
		return data, nil
	}

	return resp, args.Error(1)
}

func (m *userMock) OthAddDeps(ctx context.Context, reqs *AddDepartmentRequest) (*AddListResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*AddListResponse)
	resp := &AddListResponse{}
	if len(reqs.Deps) > 0 {
		data, ok := res["-1"]
		if !ok {
			// 以空数据返回
		}
		return data, nil
	}

	return resp, args.Error(1)
}

func (m *userMock) GetUserInfo(ctx context.Context, reqs *OneUserRequest) (*OneUserResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*OneUserResponse)

	if reqs.ID != "" {
		data, ok := res["-1"]
		if !ok {
			// 以空数据返回
		}
		return data, nil
	}

	return nil, args.Error(1)
}

func (m *userMock) GetUserByIDs(ctx context.Context, reqs *GetUserByIDsRequest) (*GetUserByIDsResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*GetUserByIDsResponse)
	resp := new(GetUserByIDsResponse)
	if len(reqs.IDs) > 0 {
		data, ok := res["-1"]
		if !ok {
			// 以空数据返回
		}
		return data, nil
	}

	return resp, args.Error(1)
}

func (m *userMock) GetUsersByDepID(ctx context.Context, r *GetUsersByDepIDRequest) (*GetUsersByDepIDResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (m *userMock) GetDepByIDs(ctx context.Context, reqs *GetDepByIDsRequest) (*GetDepByIDsResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*GetDepByIDsResponse)
	resp := new(GetDepByIDsResponse)
	if reqs != nil {
		data, ok := res["-1"]
		if !ok {
			// 以空数据返回
		}
		return data, nil
	}

	return resp, args.Error(1)
}

func (m *userMock) GetDepMaxGrade(ctx context.Context, reqs *GetDepMaxGradeRequest) (*GetDepMaxGradeResponse, error) {
	args := m.Called()
	res := args.Get(0).(map[string]*GetDepMaxGradeResponse)
	resp := new(GetDepMaxGradeResponse)
	if reqs != nil {
		data, ok := res["-1"]
		if !ok {
			// 以空数据返回
		}
		return data, nil
	}

	return resp, args.Error(1)
}
