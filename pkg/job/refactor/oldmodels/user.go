package oldmodels

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
import "gorm.io/gorm"

// User user
type User struct {
	ID             string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	UserName       string `gorm:"column:user_name;type:varchar(64);" json:"userName"`
	Phone          string `gorm:"column:phone;type:varchar(64);" json:"phone"`
	Email          string `gorm:"column:email;type:varchar(64);" json:"email"`
	IDCard         string `gorm:"column:id_card;type:varchar(64);" json:"idCard"`
	Address        string `gorm:"column:address;type:varchar(200);" json:"address"`
	BankCardNumber string `gorm:"column:bank_card_number;type:varchar(64);" json:"bankCardNumber"`
	BankAddress    string `gorm:"column:bank_address;type:varchar(64);" json:"bankAddress"`
	LeaderID       string `gorm:"column:leader_id;type:varchar(64); " json:"leaderId"`
	UseStatus      int    `gorm:"column:use_status; " json:"useStatus"`
	CompanyID      string `gorm:"column:company_id;type:varchar(64); " json:"companyId"`
	Position       string `gorm:"column:position;type:varchar(64); " json:"position"`
	CreateTime     int64  `gorm:"column:create_time;type:bigint; " json:"createTime"`
	UpdateTime     int64  `gorm:"column:update_time;type:bigint; " json:"updateTime"`
	CreatBy        string `gorm:"column:creat_by;type:varchar(64); " json:"creatBy"`
	Avatar         string `gorm:"column:avatar;type:text; " json:"avatar"`
	PasswordStatus int    `gorm:"column:password_status; " json:"passwordStatus"`
}

//TableName table name
func (User) TableName() string {
	return "t_user"
}

// UserRepo interface
type UserRepo interface {
	All(db *gorm.DB) (res []User)
}

type userRepo struct {
}

// All all
func (d *userRepo) All(db *gorm.DB) (one []User) {
	list := make([]User, 0)
	affected := db.Find(&list).RowsAffected
	if affected > 0 {
		return list
	}
	return nil
}

//NewUserRepo new
func NewUserRepo() UserRepo {
	return new(userRepo)
}
