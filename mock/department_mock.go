// Code generated by MockGen. DO NOT EDIT.
// Source: department.go

// Package org is a generated GoMock package.
package mock

import (
	context "context"
	"errors"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockDepartmentRepo is a mock of DepartmentRepo interface.
type MockDepartmentRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDepartmentRepoMockRecorder
}

// MockDepartmentRepoMockRecorder is the mock recorder for MockDepartmentRepo.
type MockDepartmentRepoMockRecorder struct {
	mock *MockDepartmentRepo
}

var departments = []org.Department{{ID: "1", Name: "test", PID: ""}, {ID: "2", Name: "test1", PID: "1"}}
var departments2 = []org.Department{{ID: "3", Name: "test", PID: "2"}, {ID: "4", Name: "test1", PID: "2"}}

// NewMockDepartmentRepo creates a new mock instance.
func NewMockDepartmentRepo(ctrl *gomock.Controller) *MockDepartmentRepo {
	mock := &MockDepartmentRepo{ctrl: ctrl}
	mock.recorder = &MockDepartmentRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDepartmentRepo) EXPECT() *MockDepartmentRepoMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockDepartmentRepo) Count(ctx context.Context, db *gorm.DB, status int) int64 {

	ret := m.ctrl.Call(m, "Count", ctx, db, status)
	ret0, _ := ret[0].(int64)
	return ret0
}

// Count indicates an expected call of Count.
func (mr *MockDepartmentRepoMockRecorder) Count(ctx, db, status interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDepartmentRepo)(nil).Count), ctx, db, status)
}

// Delete mocks base method.
func (m *MockDepartmentRepo) Delete(ctx context.Context, tx *gorm.DB, id ...string) error {

	varargs := []interface{}{ctx, tx}
	for _, a := range id {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDepartmentRepoMockRecorder) Delete(ctx, tx interface{}, id ...interface{}) *gomock.Call {

	varargs := append([]interface{}{ctx, tx}, id...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDepartmentRepo)(nil).Delete), varargs...)
}

// Get mocks base method.
func (m *MockDepartmentRepo) Get(ctx context.Context, db *gorm.DB, id string) *org.Department {

	ret := m.ctrl.Call(m, "Get", ctx, db, id)
	ret0, _ := ret[0].(*org.Department)

	for k := range departments {
		if id == departments[k].ID {
			department := departments[k]
			department.ID = id
			ret0 = &department
			break
		}
	}
	for k := range departments2 {
		if id == departments2[k].ID {
			department := departments2[k]
			department.ID = id
			ret0 = &department
			break
		}
	}

	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockDepartmentRepoMockRecorder) Get(ctx, db, id interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDepartmentRepo)(nil).Get), ctx, db, id)
}

// GetMaxGrade mocks base method.
func (m *MockDepartmentRepo) GetMaxGrade(ctx context.Context, db *gorm.DB) int64 {

	ret := m.ctrl.Call(m, "GetMaxGrade", ctx, db)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetMaxGrade indicates an expected call of GetMaxGrade.
func (mr *MockDepartmentRepoMockRecorder) GetMaxGrade(ctx, db interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaxGrade", reflect.TypeOf((*MockDepartmentRepo)(nil).GetMaxGrade), ctx, db)
}

// Insert mocks base method.
func (m *MockDepartmentRepo) Insert(ctx context.Context, tx *gorm.DB, req *org.Department) error {

	ret := m.ctrl.Call(m, "Insert", ctx, tx, req)
	ret0, _ := ret[0].(error)
	for k := range departments {
		if departments[k].ID == req.ID {
			return errors.New("has data")
		}
	}
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockDepartmentRepoMockRecorder) Insert(ctx, tx, req interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockDepartmentRepo)(nil).Insert), ctx, tx, req)
}

// InsertBranch mocks base method.
func (m *MockDepartmentRepo) InsertBranch(tx *gorm.DB, req ...org.Department) error {

	varargs := []interface{}{tx}
	for _, a := range req {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertBranch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBranch indicates an expected call of InsertBranch.
func (mr *MockDepartmentRepoMockRecorder) InsertBranch(tx interface{}, req ...interface{}) *gomock.Call {

	varargs := append([]interface{}{tx}, req...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBranch", reflect.TypeOf((*MockDepartmentRepo)(nil).InsertBranch), varargs...)
}

// List mocks base method.
func (m *MockDepartmentRepo) List(ctx context.Context, db *gorm.DB, id ...string) []org.Department {

	varargs := []interface{}{ctx, db}
	for _, a := range id {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]org.Department)
	ret0 = append(ret0, departments...)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockDepartmentRepoMockRecorder) List(ctx, db interface{}, id ...interface{}) *gomock.Call {

	varargs := append([]interface{}{ctx, db}, id...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDepartmentRepo)(nil).List), varargs...)
}

// PageList mocks base method.
func (m *MockDepartmentRepo) PageList(ctx context.Context, db *gorm.DB, status, page, limit int) ([]org.Department, int64) {

	ret := m.ctrl.Call(m, "PageList", ctx, db, status, page, limit)

	ret1, _ := ret[1].(int64)
	return departments, ret1
}

// PageList indicates an expected call of PageList.
func (mr *MockDepartmentRepoMockRecorder) PageList(ctx, db, status, page, limit interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PageList", reflect.TypeOf((*MockDepartmentRepo)(nil).PageList), ctx, db, status, page, limit)
}

// SelectByPID mocks base method.
func (m *MockDepartmentRepo) SelectByPID(ctx context.Context, db *gorm.DB, pid string, status, page, limit int) ([]org.Department, int64) {

	_ = m.ctrl.Call(m, "SelectByPID", ctx, db, pid, status, page, limit)
	//ret0, _ := ret[0].([]org.Department)
	//ret1, _ := ret[0].(int64)
	var flag = false
	for k := range departments {
		if departments[k].ID == pid {
			flag = true
			break
		}
	}
	if flag {
		return departments[1:], int64(len(departments[1:]))
	}
	return nil, 0
}

// SelectByPID indicates an expected call of SelectByPID.
func (mr *MockDepartmentRepoMockRecorder) SelectByPID(ctx, db, pid, status, page, limit interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByPID", reflect.TypeOf((*MockDepartmentRepo)(nil).SelectByPID), ctx, db, pid, status, page, limit)
}

// SelectByPIDAndName mocks base method.
func (m *MockDepartmentRepo) SelectByPIDAndName(ctx context.Context, db *gorm.DB, pid, name string) *org.Department {

	ret := m.ctrl.Call(m, "SelectByPIDAndName", ctx, db, pid, name)
	ret0, _ := ret[0].(*org.Department)
	ret0 = nil
	return ret0
}

// SelectByPIDAndName indicates an expected call of SelectByPIDAndName.
func (mr *MockDepartmentRepoMockRecorder) SelectByPIDAndName(ctx, db, pid, name interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByPIDAndName", reflect.TypeOf((*MockDepartmentRepo)(nil).SelectByPIDAndName), ctx, db, pid, name)
}

// SelectByPIDs mocks base method.
func (m *MockDepartmentRepo) SelectByPIDs(ctx context.Context, db *gorm.DB, status int, pid ...string) []org.Department {

	varargs := []interface{}{ctx, db, status}
	for _, a := range pid {
		varargs = append(varargs, a)
	}
	_ = m.ctrl.Call(m, "SelectByPIDs", varargs...)
	//ret0, _ := ret[0].([]org.Department)
	var flag = false
A:
	for k := range departments2 {
		for _, v1 := range pid {
			if departments2[k].ID == v1 {
				flag = true
				break A
			}
		}
	}
	if flag {
		return nil
	}
	return departments2
}

// SelectByPIDs indicates an expected call of SelectByPIDs.
func (mr *MockDepartmentRepoMockRecorder) SelectByPIDs(ctx, db, status interface{}, pid ...interface{}) *gomock.Call {

	varargs := append([]interface{}{ctx, db, status}, pid...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByPIDs", reflect.TypeOf((*MockDepartmentRepo)(nil).SelectByPIDs), varargs...)
}

// SelectSupper mocks base method.
func (m *MockDepartmentRepo) SelectSupper(ctx context.Context, db *gorm.DB) *org.Department {

	ret := m.ctrl.Call(m, "SelectSupper", ctx, db)
	ret0, _ := ret[0].(*org.Department)
	ret0 = nil
	return ret0
}

// SelectSupper indicates an expected call of SelectSupper.
func (mr *MockDepartmentRepoMockRecorder) SelectSupper(ctx, db interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectSupper", reflect.TypeOf((*MockDepartmentRepo)(nil).SelectSupper), ctx, db)
}

// Update mocks base method.
func (m *MockDepartmentRepo) Update(ctx context.Context, tx *gorm.DB, req *org.Department) error {

	ret := m.ctrl.Call(m, "Update", ctx, tx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockDepartmentRepoMockRecorder) Update(ctx, tx, req interface{}) *gomock.Call {

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDepartmentRepo)(nil).Update), ctx, tx, req)
}
