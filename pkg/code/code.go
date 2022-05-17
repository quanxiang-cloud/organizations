package code

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
import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// InvalidURI invalid url
	InvalidURI = 50014000000
	// InvalidParams invalid params
	InvalidParams = 50014000001
	// InvalidTimestamp invalid timestamp
	InvalidTimestamp = 50014000002
	// NameUsed name used
	NameUsed = 50014000003
	// InvalidDELDEP invalid delete
	InvalidDELDEP = 50014000004
	// InvalidFile invalid file
	InvalidFile = 50034000006
	// AccountPasswordCountErr account password err count
	AccountPasswordCountErr = 50034000007
	// CanNotDel can not del
	CanNotDel = 50014000008
	// InvalidAccount invalid account
	InvalidAccount = 50034000009
	// InvalidVerificationCode invalid code
	InvalidVerificationCode = 50044000010
	// ValidVerificationCode  valid verification code
	ValidVerificationCode = 50044000011
	// AccountExist account exist
	AccountExist = 50034000012
	// ChangeDepErr change department err
	ChangeDepErr = 50014000013
	// InvalidUpdate invalid update
	InvalidUpdate = 50014000014
	// InvalidDEPID invalid dep id
	InvalidDEPID = 50014000015
	// LockedAccount locked account
	LockedAccount = 50034000016
	// CanNotMoveParentToChild can not move parent to child node
	CanNotMoveParentToChild = 50014000017
	// ResetAccountPasswordErr reset password err
	ResetAccountPasswordErr = 50034000018
	// InvalidPWD invalid password
	InvalidPWD = 50034000019
	// InvalidPhone invalid phone
	InvalidPhone = 50034000020
	// TopDepExist supper department exist
	TopDepExist = 50034000021
	// SystemParameter default system param
	SystemParameter = 50034000022
	// ColumnExist column exist
	ColumnExist = 50034000023
	// FieldNameIsNull  name  is null
	FieldNameIsNull = 50034000024
	// InvalidEmail invaild email
	InvalidEmail = 50034000025
	// EmailRequired email is must
	EmailRequired = 50034000026
	// SelfEmailRequired self email is must
	SelfEmailRequired = 50034000027
	// BatchDeleteNotSupported not support batch del
	BatchDeleteNotSupported = 50034000028
	// MismatchPasswordRule not match password rule
	MismatchPasswordRule = 50034000029
	// DataNotExist data not exist
	DataNotExist = 50034000030
	// CanNotModifyYourself can not modify self
	CanNotModifyYourself = 50034000031
	// NotExistAccountErr account not exist
	NotExistAccountErr = 50034000032
	// ExpireVerificationCode code was expired
	ExpireVerificationCode = 50034000033
	// ErrTooLong too long
	ErrTooLong = 50034000034
	// ErrInvalidRuleAccount not match account rule
	ErrInvalidRuleAccount = 50034000035
	// ErrFirstResetInvalid first reset invalid
	ErrFirstResetInvalid = 50034000036
	// ErrHasBeActive some data is active
	ErrHasBeActive = 50034000037
	// ErrFieldColumnUsed alias column is open
	ErrFieldColumnUsed = 50034000038
	// ErrCircleData make a circle data
	ErrCircleData = 50034000039
	// ErrColumnExist make a circle data
	ErrColumnExist = 50034000040
	// ErrNoPower no power
	ErrNoPower = 50034000041
)

// CodeTable 码表
var CodeTable = error2.Table{
	InvalidURI:              "无效的URI.",
	InvalidParams:           "无效或错误参数.",
	InvalidTimestamp:        "无效的时间格式.",
	NameUsed:                "名称已被使用！请检查后重试！",
	InvalidDELDEP:           "当前部门下还存在关联用户或子部门，不能进操作！",
	InvalidFile:             "无效的文件.",
	AccountPasswordCountErr: "账号或密码(验证码)错误，请检查后重试！你还有%d次尝试机会！",
	InvalidAccount:          "账号不存在或已被禁用，请检查后重试！",
	CanNotDel:               "删除失败！当前节点下还存在子部门或者当前节点为根节点！",
	InvalidVerificationCode: "无效的验证码！",
	ValidVerificationCode:   "验证码未失效!",
	AccountExist:            "账户已存在！请检查后重试！",
	ChangeDepErr:            "批量调整部门异常，请检查后再试！",
	InvalidUpdate:           "修改的对象不存在，请检查后再试！",
	InvalidDEPID:            "上级部门不能选择自己，请检查后重试！",
	LockedAccount:           "密码(验证码)错误达到安全限制，账户已锁定，请于24小时之后重试或者联系管理员！",
	CanNotMoveParentToChild: "上级部门不能选择当前部门的子部门，请检查后重试！",
	ResetAccountPasswordErr: "重置密码错误，请检查后重试！",
	InvalidPWD:              "密码格式有误！请输入包含英文、特殊字符、数字且不少于8位的密码！",
	InvalidPhone:            "电话号码格式错误！请请输入正确电话号码！",
	TopDepExist:             "顶级部门已存在！",
	SystemParameter:         "系统默认参数不可删除或修改！",
	ColumnExist:             "字段或别名已经存在！",
	FieldNameIsNull:         "字段名字为空，请先完善！",
	InvalidEmail:            "无效的邮箱！",
	EmailRequired:           "邮箱不能为空！",
	SelfEmailRequired:       "私人邮箱不能为空！",
	BatchDeleteNotSupported: "不支持批量删除！",
	MismatchPasswordRule:    "不匹配密码规则",
	DataNotExist:            "数据不存在，无法操作",
	CanNotModifyYourself:    "不能对自己操作",
	NotExistAccountErr:      "账号或密码(验证码)错误，请检查后重试！",
	ExpireVerificationCode:  "验证码已失效，请重新获取！",
	ErrTooLong:              "长度超过限制！",
	ErrInvalidRuleAccount:   "不符合规则的账户！",
	ErrFirstResetInvalid:    "首次重置密码已完成！",
	ErrHasBeActive:          "数据中包含已被激活数据，请选择正确数据再操作！",
	ErrFieldColumnUsed:      "扩展字段功能已被开启，请不要重复操作！",
	ErrCircleData:           "数据关系成环，请检查后提交！",
	ErrColumnExist:          "字段命名重复，请修改后再尝试！",
	ErrNoPower:              "没有权限",
}
