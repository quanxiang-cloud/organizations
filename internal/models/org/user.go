package org

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

	"gorm.io/gorm"
)

// User user system info
type User struct {
	ID        string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id,omitempty" comment:"ID"`
	Name      string `gorm:"column:name;type:varchar(64);" json:"name,omitempty" comment:"姓名"`
	Phone     string `gorm:"column:phone;type:varchar(64);" json:"phone,omitempty" comment:"手机号"`
	Email     string `gorm:"column:email;type:varchar(64);" json:"email,omitempty" comment:"邮箱"`
	SelfEmail string `gorm:"column:self_email;type:varchar(64);" json:"self_email,omitempty" comment:"私人邮箱"`
	IDCard    string `gorm:"column:id_card;type:varchar(64);" json:"idCard,omitempty" comment:"身份证"`
	Address   string `gorm:"column:address;type:varchar(200);" json:"address,omitempty" comment:"住址"`
	//1:normal，-2:invalid，-1:del，2:active,-3:no word
	UseStatus int    `gorm:"column:use_status;type:int(4); " json:"useStatus,omitempty" comment:"状态"`
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID,omitempty" comment:"租户ID"`
	Position  string `gorm:"column:position;type:varchar(64); " json:"position,omitempty" comment:"职位"`
	Avatar    string `gorm:"column:avatar;type:text; " json:"avatar,omitempty" comment:"头像"`
	JobNumber string `gorm:"column:job_number;type:text; " json:"jobNumber,omitempty" comment:"工号"`
	//0:null,1:man,2:woman
	Gender int    `gorm:"column:gender;type:int(4); " json:"gender,omitempty" comment:"性别"`
	Source string `gorm:"column:source;type:varchar(64); " json:"source,omitempty" comment:"信息来源"`
	//1:normal,0:must reset password
	PasswordStatus int    `gorm:"column:password_status;type:int(4); " json:"passwordStatus,omitempty" comment:"密码状态"`
	CreatedAt      int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" comment:"创建时间"`
	UpdatedAt      int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" comment:"更新时间"`
	DeletedAt      int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" comment:"删除时间"`
	CreatedBy      string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty" comment:"创建者"`
	UpdatedBy      string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty" comment:"修改者"`
	DeletedBy      string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty" comment:"删除者"`
}

// TableName table name
func (*User) TableName() string {
	return "org_user"
}

// UserRepo interface
type UserRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, r *User) (err error)
	InsertBranch(tx *gorm.DB, req ...*User) error
	UpdateByID(ctx context.Context, tx *gorm.DB, r *User) (err error)
	Get(ctx context.Context, db *gorm.DB, id string) (res *User)
	List(ctx context.Context, db *gorm.DB, id ...string) (list []*User)
	PageList(ctx context.Context, db *gorm.DB, status, page, limit int, userIDs []string) (list []*User, total int64)
	SelectByEmailOrPhone(ctx context.Context, db *gorm.DB, info string) (res *User)
	GetColumns(ctx context.Context, db *gorm.DB, user *User, schema string) []Columns
	Count(ctx context.Context, db *gorm.DB, status, activeStatus int) (totalUser, activeUserNum int64)
}

// Columns db column interface
type Columns interface {
	New() Columns
	GetName() string
	GetColumnName() string
	GetDataType() string
	GetCharacterMaximumLength() int
	GetNumericScale() int
}

// MysqlUserColumn 用户表字段
type MysqlUserColumn struct {
	Name                   string `gorm:"column:NAME" json:"NAME"`
	ColumnName             string `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
	DataType               string `gorm:"column:DATA_TYPE" json:"DATA_TYPE"`
	CharacterMaximumLength int    `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"CHARACTER_MAXIMUM_LENGTH"`
	NumericScale           int    `gorm:"column:NUMERIC_SCALE" json:"NUMERIC_SCALE"`
}

// NewMysqlUserColumn new
func NewMysqlUserColumn() *MysqlUserColumn {
	return &MysqlUserColumn{}
}

// New new
func (m *MysqlUserColumn) New() Columns {
	return NewMysqlUserColumn()
}

// GetName get columns alias name
func (m *MysqlUserColumn) GetName() string {
	return m.Name
}

// GetColumnName get columns name
func (m *MysqlUserColumn) GetColumnName() string {
	return m.ColumnName
}

// GetDataType get type
func (m *MysqlUserColumn) GetDataType() string {
	return m.DataType
}

// GetCharacterMaximumLength get length
func (m *MysqlUserColumn) GetCharacterMaximumLength() int {
	return m.CharacterMaximumLength
}

// GetNumericScale get point
func (m *MysqlUserColumn) GetNumericScale() int {
	return m.NumericScale
}
