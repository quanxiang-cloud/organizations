package account

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
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/logic/org/user"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/code"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/encode2"
	"github.com/quanxiang-cloud/organizations/pkg/ladp"
	"github.com/quanxiang-cloud/organizations/pkg/message"
	"github.com/quanxiang-cloud/organizations/pkg/random2"
	"github.com/quanxiang-cloud/organizations/pkg/systems"
	"github.com/quanxiang-cloud/organizations/pkg/verification"
)

// Account account interface
type Account interface {
	CheckPassword(c context.Context, account *LoginAccountRequest) (*LoginAccountResponse, error)
	UpdatePassword(c context.Context, account *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	FirstUpdatePassword(c context.Context, account *FirstSetPasswordRequest) (*FirstSetPasswordResponse, error)
	ForgetUpdatePassword(c context.Context, account *ForgetResetRequest) (*ForgetResetResponse, error)
	AdminUpdatePassword(c context.Context, accounts *AdminUpdatePasswordRequest) (*AdminUpdatePasswordResponse, error)
	GetCode(ctx context.Context, r *CodeRequest) (*CodeResponse, error)
}

const (
	loginTypePwd        = "pwd"
	loginTypeLdap       = "ldapClient"
	loginTypeCode       = "code"
	accountLength       = 50
	passwordLength      = 8
	codeLength          = 6
	redisAccountPWDErr  = "organizations:accountPWDErr:"
	resetPasswordStatus = -1
	codeKey             = "code"
)

// account
type account struct {
	DB          *gorm.DB
	accountRepo org.AccountRepo
	user        org.UserRepo
	message     message.Message
	redisClient redis.UniversalClient
	ldapClient  ldap.Ldap
	depRepo     org.DepartmentRepo
	conf        configs.Config
	userDepRepo org.UserDepartmentRelationRepo
}

// NewAccount new
func NewAccount(conf configs.Config, db *gorm.DB, redisClient redis.UniversalClient) Account {
	return &account{
		accountRepo: mysql2.NewAccountRepo(),
		DB:          db,
		message:     message.NewMessage(conf.InternalNet),
		redisClient: redisClient,
		user:        mysql2.NewUserRepo(),
		ldapClient:  ldap.NewLdap(conf.InternalNet),
		depRepo:     mysql2.NewDepartmentRepo(),
		conf:        conf,
		userDepRepo: mysql2.NewUserDepartmentRelationRepo(),
	}
}

// UpdatePasswordRequest reset password request
type UpdatePasswordRequest struct {
	UserID      string `json:"userID"`
	OldPassword string `json:"oldPassword" binding:"required,password"`
	NewPassword string `json:"newPassword" binding:"required,password"`
	//1:normal,2:invalid,-1:del
	UseStatus int `json:"useStatus"`
	Password  string
}

// UpdatePasswordResponse reset password response
type UpdatePasswordResponse struct {
	UserID string `json:"userID"`
}

// UpdatePassword update password
func (u *account) UpdatePassword(c context.Context, r *UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	accounts := u.accountRepo.SelectByUserID(u.DB, r.UserID)
	if accounts == nil {
		return nil, error2.New(code.ResetAccountPasswordErr)
	}

	if accounts[0].Password == encode2.MD5Encode(r.OldPassword) {
		//todo get info from system server
		info := systems.GetSecurityInfo(c, u.conf, u.redisClient)
		f := random2.CheckPassword(r.NewPassword, info.PwdMinLen, info.PwdType)
		if !f {
			return nil, error2.New(code.MismatchPasswordRule)
		}
		tx := u.DB.Begin()
		u2 := org.Account{
			UserID:   accounts[0].UserID,
			Password: encode2.MD5Encode(r.NewPassword),
		}
		err := u.accountRepo.UpdatePasswordByUserID(tx, &u2)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		u3 := UpdatePasswordResponse{
			UserID: r.UserID,
		}
		return &u3, nil
	}
	return nil, error2.New(code.ResetAccountPasswordErr)

}

// AdminUpdatePasswordRequest admin reset password request
type AdminUpdatePasswordRequest struct {
	UserIDs     []string           `json:"userIDs"`
	CreatedBy   string             `json:"createdBy"`
	SendMessage []user.SendMessage `json:"sendMessage"`
	Header      http.Header
}

// AdminUpdatePasswordResponse admin reset password response
type AdminUpdatePasswordResponse struct {
	Users []ResetPasswordResponse `json:"Users"`
}

// ResetPasswordResponse reset password response
type ResetPasswordResponse struct {
	UserID   string `json:"userID"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

// AdminUpdatePassword admin reset password
func (u *account) AdminUpdatePassword(c context.Context, r *AdminUpdatePasswordRequest) (*AdminUpdatePasswordResponse, error) {
	tx := u.DB.Begin()
	m := make(map[string]string)

	u3 := AdminUpdatePasswordResponse{}
	send := make(map[string]user.SendMessage)
	for k := range r.SendMessage {
		send[r.SendMessage[k].UserID] = r.SendMessage[k]
	}
	for k := range r.UserIDs {
		newPWD := user.CreatePassword(c, u.conf, u.redisClient)
		u2 := org.Account{
			UserID:   r.UserIDs[k],
			Password: encode2.MD5Encode(newPWD),
		}
		err := u.accountRepo.UpdatePasswordByUserID(tx, &u2)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		m[r.UserIDs[k]] = newPWD

		res := u.user.Get(c, u.DB, r.UserIDs[k])
		if r.SendMessage != nil && r.SendMessage[k].SendChannel != user.NO {
			user.SendAccountAndPWDOrCode(c, u.message, send[res.ID].SendTo, "", u.conf.MessageTemplate.ResetPWD, newPWD, r.SendMessage[k].SendChannel)
		}
		response := ResetPasswordResponse{}
		response.UserID = res.ID
		if u.conf.POC {
			response.Password = newPWD
			response.Email = res.Email
			response.Name = res.Name
		}
		u3.Users = append(u3.Users, response)
	}

	tx.Commit()

	return &u3, nil
}

// LoginAccountRequest login request
type LoginAccountRequest struct {
	// email or phone
	UserName string `json:"userName" binding:"required,max=60,emailOrPhone"`
	Password string `json:"password" binding:"required,min=6"`
	// login type: pwd,ldapClient,code
	Types  string `json:"types" binding:"required"` //登录模式
	Header http.Header
}

// LoginAccountResponse login response
type LoginAccountResponse struct {
	UserID string `json:"userID"`
	//1:normal，-2:invalid，-1:del，2:active
	UseStatus int    `json:"useStatus"`
	Name      string `json:"-"`
	UserName  string `json:"-"`
}

// CheckPassword check password
func (u *account) CheckPassword(c context.Context, r *LoginAccountRequest) (*LoginAccountResponse, error) {
	res := &LoginAccountResponse{}
	acc := u.accountRepo.SelectByAccount(u.DB, r.UserName)
	if acc == nil {
		return nil, error2.New(code.NotExistAccountErr)
	}
	oldUser := u.user.Get(c, u.DB, acc.UserID)
	res.UserName = oldUser.Name
	res.UserID = oldUser.ID
	res.UseStatus = oldUser.UseStatus
	if oldUser.UseStatus != consts.NormalStatus {
		return res, error2.New(code.InvalidAccount)
	}
	res.UserID = oldUser.ID
	res.Name = oldUser.Email
	res.UseStatus = oldUser.UseStatus

	//TODO get info from system server
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)

	val, err1 := u.redisClient.Get(c, redisAccountPWDErr+acc.UserID).Result()
	if err1 != nil {
		if err1 != redis.Nil {
			logger.Logger.Error(err1)
			return res, err1
		}
	}
	var errNum = 0
	if val != "" {
		atoi, err := strconv.Atoi(val)
		if err != nil {
			return res, err
		}
		errNum = atoi
	}

	if errNum >= int(info.PwdCount) {
		return nil, error2.New(code.LockedAccount)
	}

	if u.conf.Ldap.Open {
		split := strings.Split(r.UserName, "@")
		if split[1] == u.conf.Ldap.Regex {
			r.Types = loginTypeLdap
		}
	}
	var flag = false
	var err error = nil
	switch r.Types {
	case loginTypePwd:
		flag, err = u.pwd(c, r, acc.Password)
	case loginTypeLdap:
		c = context.WithValue(c, user.TenantID, oldUser.TenantID)
		flag, err = u.ldap(c, r.Header, r)
	case loginTypeCode:
		flag, err = u.code(c, r)
	}
	if err != nil {
		return nil, err
	}
	if !flag {
		err = error2.New(code.AccountPasswordCountErr, int(info.PwdCount)-(errNum+1))
		u.redisClient.SetEX(c, redisAccountPWDErr+acc.UserID, errNum+1, time.Duration(info.PwdCountWait)*time.Minute)
		return nil, err
	}
	u.redisClient.Del(c, redisAccountPWDErr+acc.UserID)
	return res, nil

}

// pwd password
func (u *account) pwd(ctx context.Context, r *LoginAccountRequest, comparePassword string) (bool, error) {
	if encode2.MD5Encode(r.Password) == comparePassword {
		return true, nil
	}
	return false, nil
}

// ldap ldap
func (u *account) ldap(ctx context.Context, header http.Header, r *LoginAccountRequest) (bool, error) {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	authReq := &ldap.AuthReq{
		TenantID: tenantID,
		UserName: r.UserName,
		Password: r.Password,
	}
	auth, err := u.ldapClient.Auth(ctx, header, authReq)
	if err != nil {
		return false, err
	}
	return auth.Flag, nil
}

// code email or phone code
func (u *account) code(ctx context.Context, account *LoginAccountRequest) (bool, error) {
	cacheCode := u.redisClient.Get(ctx, u.conf.VerificationCode.LoginCode+":"+account.UserName).Val()
	if cacheCode == "" {
		return false, error2.New(code.ExpireVerificationCode)
	}
	if account.Password != cacheCode || account.Password == "" {
		return false, error2.New(code.InvalidVerificationCode)
	}
	u.redisClient.Del(ctx, u.conf.VerificationCode.LoginCode+":"+account.UserName)
	return true, nil
}

// CodeRequest get code request
type CodeRequest struct {
	UserName string `json:"userName" form:"userName" binding:"required,max=60,emailOrPhone"`
	Model    string
}

// CodeResponse code response
type CodeResponse struct {
	Code string `json:"code"`
}

// GetCode get code
func (u *account) GetCode(ctx context.Context, r *CodeRequest) (*CodeResponse, error) {
	if r.Model != u.conf.VerificationCode.RegisterCode {
		acc := u.accountRepo.SelectByAccount(u.DB, r.UserName)
		if acc == nil {
			return nil, error2.New(code.InvalidAccount)
		}
		user := u.user.Get(ctx, u.DB, acc.UserID)
		ctx = context.WithValue(ctx, user.TenantID, user.TenantID)
	}

	rd := strings.ToLower(random2.RandomString(codeLength, 6))
	var templateAlias = ""
	switch r.Model {
	case u.conf.VerificationCode.LoginCode:
		b := u.redisClient.Exists(ctx, u.conf.VerificationCode.LoginCode+":"+r.UserName).Val()
		if b == 1 {
			return nil, error2.New(code.ValidVerificationCode)
		}
		u.redisClient.SetEX(ctx, u.conf.VerificationCode.LoginCode+":"+r.UserName, rd, u.conf.VerificationCode.ExpireTime*time.Second)
		templateAlias = u.conf.MessageTemplate.LoginCode
	case u.conf.VerificationCode.ResetCode:
		b := u.redisClient.Exists(ctx, u.conf.VerificationCode.ResetCode+":"+r.UserName).Val()
		if b == 1 {
			return nil, error2.New(code.ValidVerificationCode)
		}
		u.redisClient.SetEX(ctx, u.conf.VerificationCode.ResetCode+":"+r.UserName, rd, u.conf.VerificationCode.ExpireTime*time.Second)
		templateAlias = u.conf.MessageTemplate.ResetCode

	case u.conf.VerificationCode.ForgetCode:
		b := u.redisClient.Exists(ctx, u.conf.VerificationCode.ForgetCode+":"+r.UserName).Val()
		if b == 1 {
			return nil, error2.New(code.ValidVerificationCode)
		}
		u.redisClient.SetEX(ctx, u.conf.VerificationCode.ForgetCode+":"+r.UserName, rd, u.conf.VerificationCode.ExpireTime*time.Second)
		templateAlias = u.conf.MessageTemplate.ForgetCode
	case u.conf.VerificationCode.RegisterCode:
		b := u.redisClient.Exists(ctx, u.conf.VerificationCode.RegisterCode+":"+r.UserName).Val()
		if b == 1 {
			return nil, error2.New(code.ValidVerificationCode)
		}
		u.redisClient.SetEX(ctx, u.conf.VerificationCode.RegisterCode+":"+r.UserName, rd, u.conf.VerificationCode.ExpireTime*time.Second)
		templateAlias = u.conf.MessageTemplate.RegisterCode
	}

	if len(r.UserName) > accountLength {
		return nil, error2.NewErrorWithString(code.ErrTooLong, "接收信息账户超过限制长度")
	}

	if verification.CheckEmail(r.UserName) {
		user.SendAccountAndPWDOrCode(ctx, u.message, r.UserName, "", templateAlias, rd, 1)
	} else if verification.CheckPhone(r.UserName) {
		user.SendAccountAndPWDOrCode(ctx, u.message, r.UserName, "", templateAlias, rd, 2)
	} else {
		return nil, error2.New(code.ErrInvalidRuleAccount)
	}
	res := &CodeResponse{}
	res.Code = rd
	return res, nil

}

// ForgetResetRequest reset password when forgot password
type ForgetResetRequest struct {
	//email、phone
	UserName    string `json:"userName" binding:"required,max=60,emailOrPhone"`
	NewPassword string `json:"newPassword" binding:"required,password"`
	Code        string `json:"code" binding:"required"`
	Header      http.Header
}

// ForgetResetResponse reset password when forgot password response
type ForgetResetResponse struct {
	UserID string `json:"userID"`
}

// ForgetUpdatePassword reset password when forgot password
func (u *account) ForgetUpdatePassword(c context.Context, r *ForgetResetRequest) (*ForgetResetResponse, error) {
	acc := u.accountRepo.SelectByAccount(u.DB, r.UserName)
	if acc == nil {
		return nil, error2.New(code.NotExistAccountErr)
	}
	oldUser := u.user.Get(c, u.DB, acc.UserID)
	if oldUser.UseStatus != consts.NormalStatus {
		return nil, error2.New(code.InvalidAccount)
	}

	val := u.redisClient.Get(c, u.conf.VerificationCode.ForgetCode+":"+r.UserName).Val()
	if val == "" {
		return nil, error2.New(code.ExpireVerificationCode)
	}
	if val != r.Code {
		return nil, error2.New(code.InvalidVerificationCode)
	}
	//todo get info from system server
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)

	f := random2.CheckPassword(r.NewPassword, info.PwdMinLen, info.PwdType)
	if !f {
		return nil, error2.New(code.MismatchPasswordRule)
	}
	tx := u.DB.Begin()
	u2 := org.Account{
		UserID:   oldUser.ID,
		Password: encode2.MD5Encode(r.NewPassword),
	}
	err := u.accountRepo.UpdatePasswordByUserID(tx, &u2)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	u3 := ForgetResetResponse{
		UserID: oldUser.ID,
	}
	u.redisClient.Del(c, u.conf.VerificationCode.ForgetCode+":"+r.UserName)
	return &u3, nil

}

// ResetLdapPassword reset ldapClient password
func ResetLdapPassword(ctx context.Context, header http.Header, ldapClient ldap.Ldap, id, email, password string, depNumberID int64) error {
	updateReq := &ldap.UserUpdatePasswordReq{}
	updateReq.ID = id
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	updateReq.TenantID = tenantID
	updateReq.Mail = email
	updateReq.UserPassword = password
	updateReq.GidNumber = strconv.FormatInt(depNumberID, 10)
	_, err := ldapClient.UpdatePassword(ctx, header, updateReq)
	return err
}

// FirstSetPasswordRequest first login reset password request
type FirstSetPasswordRequest struct {
	UserID      string `json:"userID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword" binding:"required,password"`
}

// FirstSetPasswordResponse  first login reset password response
type FirstSetPasswordResponse struct {
}

// FirstUpdatePassword first login reset password request
func (u *account) FirstUpdatePassword(c context.Context, r *FirstSetPasswordRequest) (*FirstSetPasswordResponse, error) {
	accounts := u.accountRepo.SelectByUserID(u.DB, r.UserID)
	if accounts == nil {
		return nil, error2.New(code.ResetAccountPasswordErr)
	}
	oldUser := u.user.Get(c, u.DB, r.UserID)
	//todo get info from system server
	info := systems.GetSecurityInfo(c, u.conf, u.redisClient)

	if oldUser.PasswordStatus&consts.NormalStatus != 0 {
		return nil, error2.New(code.ErrFirstResetInvalid)
	}
	f := random2.CheckPassword(r.NewPassword, info.PwdMinLen, info.PwdType)
	if !f {
		return nil, error2.New(code.MismatchPasswordRule)
	}
	tx := u.DB.Begin()
	u2 := &org.Account{
		UserID:   r.UserID,
		Password: encode2.MD5Encode(r.NewPassword),
	}
	err := u.accountRepo.UpdatePasswordByUserID(tx, u2)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	oldUser.PasswordStatus = oldUser.PasswordStatus + 1
	err = u.user.UpdateByID(c, tx, oldUser)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	u3 := FirstSetPasswordResponse{}

	return &u3, nil

}
