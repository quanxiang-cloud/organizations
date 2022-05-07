// Code generated by MockGen. DO NOT EDIT.
// Source: landlord.go

// Package landlord is a generated GoMock package.
package mock

import (
	context "context"
	"github.com/quanxiang-cloud/organizations/pkg/landlord"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLandlord is a mock of Landlord interface.
type MockLandlord struct {
	ctrl     *gomock.Controller
	recorder *MockLandlordMockRecorder
}

// MockLandlordMockRecorder is the mock recorder for MockLandlord.
type MockLandlordMockRecorder struct {
	mock *MockLandlord
}

// NewMockLandlord creates a new mock instance.
func NewMockLandlord(ctrl *gomock.Controller) *MockLandlord {
	mock := &MockLandlord{ctrl: ctrl}
	mock.recorder = &MockLandlordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLandlord) EXPECT() *MockLandlordMockRecorder {
	return m.recorder
}

// CancelRelation mocks base method.
func (m *MockLandlord) CancelRelation(ctx context.Context, header http.Header, r *landlord.CancelRelationRequest) (*landlord.CancelRelationResponse, error) {

	ret := m.ctrl.Call(m, "CancelRelation", ctx, header, r)
	ret0, _ := ret[0].(*landlord.CancelRelationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CancelRelation indicates an expected call of CancelRelation.
func (mr *MockLandlordMockRecorder) CancelRelation(ctx, header, r interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelRelation", reflect.TypeOf((*MockLandlord)(nil).CancelRelation), ctx, header, r)
}

// Register mocks base method.
func (m *MockLandlord) Register(ctx context.Context, header http.Header, r *landlord.RegisterRequest) (*landlord.RegisterResponse, error) {

	_ = m.ctrl.Call(m, "Register", ctx, header, r)
	//ret0, _ := ret[0].(*RegisterResponse)
	//ret1, _ := ret[1].(error)
	response := &landlord.RegisterResponse{
		ID: "test",
	}
	return response, nil
}

// Register indicates an expected call of Register.
func (mr *MockLandlordMockRecorder) Register(ctx, header, r interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockLandlord)(nil).Register), ctx, header, r)
}

// SetRelation mocks base method.
func (m *MockLandlord) SetRelation(ctx context.Context, header http.Header, r *landlord.SetRelationRequest) (*landlord.SetRelationResponse, error) {

	ret := m.ctrl.Call(m, "SetRelation", ctx, header, r)
	ret0, _ := ret[0].(*landlord.SetRelationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetRelation indicates an expected call of SetRelation.
func (mr *MockLandlordMockRecorder) SetRelation(ctx, header, r interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRelation", reflect.TypeOf((*MockLandlord)(nil).SetRelation), ctx, header, r)
}

// MockrespData is a mock of respData interface.
type MockrespData struct {
	ctrl     *gomock.Controller
	recorder *MockrespDataMockRecorder
}

// MockrespDataMockRecorder is the mock recorder for MockrespData.
type MockrespDataMockRecorder struct {
	mock *MockrespData
}

// NewMockrespData creates a new mock instance.
func NewMockrespData(ctrl *gomock.Controller) *MockrespData {
	mock := &MockrespData{ctrl: ctrl}
	mock.recorder = &MockrespDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrespData) EXPECT() *MockrespDataMockRecorder {
	return m.recorder
}
