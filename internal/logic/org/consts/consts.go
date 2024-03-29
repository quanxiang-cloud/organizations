package consts

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
//Common const
const (
	DelStatus = -1

	NormalStatus = 1

	UnWork = -3

	ActiveStatus = 2

	FirsGrade = 1

	UnNormalStatus = -2

	RedisTokenUserInfo = "organizations:users:"

	RedisTokenUserInfoEx = 60 //minute

	ResetPasswordStatus = 0

	SystemAttr = 1

	AliasAttr = 2

	AllAttr = 0

	FieldAdminStatus = 0

	FieldViewerStatus = 1

	OwnerDepName = "所在部门名称"

	NamePhoneEmailNotNull = "姓名、私人邮箱、公司邮箱不能为空！"

	NameLengthIsLong = "姓名长度超过限制"

	NotPhone = "手机格式不正确"

	NotDepartment = "部门错误或者不存在"

	EmailLengthIsLong = "邮箱超过长度"

	NotEmail = "邮箱格式不正确"

	EmailPhoneRepeat = "邮箱或者手机重复"

	EmailPhoneExist = "邮箱或者手机号已存在"

	EmailExist = "帐户邮箱已被占用"

	PhoneExist = "手机帐户已被占用"

	RelationDepartmentFail = "关联部门失败"
)

// SYSTEM column
const (
	ID = "id"

	NAME = "name"

	EMAIL = "email"

	PHONE = "phone"

	SELFEMAIL = "selfEmail"

	ADDRESS = "address"

	LEADERID = "leaderID"

	IDCARD = "idCard"

	AVATAR = "avatar"

	PASSWORDSTATUS = "passwordStatus"

	UPDATEDAT = "updatedAt"

	CREATEDAT = "createdAt"

	UPDATEDBY = "updatedBy"

	CREATEDBY = "createdBy"

	USESTATUS = "useStatus"

	COMPANYID = "companyID"

	REMARK = "remark"

	DEPNAME = "depName"

	DEPID = "depID"

	TENANTID = "tenantID"

	JOBNUMBER = "jobNumber"

	GENDER = "gender"

	ENTRYTIME = "entryTime"

	SOURCE = "source"
)

// Field type
const (
	STRING = "string"

	TEXT = "text"

	INT = "int"

	INT64 = "int64"

	FLOAT = "float"

	TIME = "time"
)

// DBColumns key is for front, value is db column type
var DBColumns = map[string]string{
	"string":   "varchar",
	"text":     "text",
	"longtext": "longtext",
	"int":      "int",
	"float":    "decimal",
	"time":     "bigint",
}

// FrontColumns key is for db column type, value is for  front
var FrontColumns = map[string]string{
	"varchar":  "string",
	"text":     "text",
	"longtext": "longtext",
	"int":      "int",
	"bigint":   "time",
	"decimal":  "float",
}
