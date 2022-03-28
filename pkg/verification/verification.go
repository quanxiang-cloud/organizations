package verification

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
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidation register validation
func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("phone", phone)
		err = v.RegisterValidation("emailOrPhone", emailOrPhone)
		err = v.RegisterValidation("password", password)
		if err != nil {
			panic(err)
		}
	}
}

const (
	phoneRegexString = `^1[3-9]\d{9}$`
	emailRegexString = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
)

var phone validator.Func = func(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if ok {
		return CheckPhone(s)
	}
	return false
}

var emailOrPhone validator.Func = func(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if ok {
		if CheckEmail(s) || CheckPhone(s) {
			return true
		}
		return false
	}
	return false
}

// CheckEmail 检查是否邮箱
func CheckEmail(userName string) bool {
	emailCompile := regexp.MustCompile(emailRegexString)
	return emailCompile.MatchString(userName)
}

// CheckPhone 检查是否手机
func CheckPhone(userName string) bool {
	phoneCompile := regexp.MustCompile(phoneRegexString)
	return phoneCompile.MatchString(userName)
}

var password validator.Func = func(fl validator.FieldLevel) bool {
	fmt.Println("校验密码")
	s, ok := fl.Field().Interface().(string)
	if ok {
		return CheckPassword(s)
	}
	return false
}

// CheckPassword 校验密码
func CheckPassword(s string) bool {
	var lenFlag = false
	if len(s) >= 8 {
		lenFlag = true
	}

	//字母数字
	var numFlag = false
	//字母
	var azFlag = false
	//特殊符号
	var flag = false
	runes := []rune(s)
	for _, v := range runes {
		if v >= 48 && v <= 57 {
			numFlag = true
			continue
		}
		if (v >= 65 && v <= 90) || (v >= 97 && v <= 122) {
			azFlag = true
			continue
		}
		if (v >= 33 && v <= 38) || (v >= 42 && v <= 46) || (v >= 63 || v <= 64) || (v >= 94 || v <= 96) || (v == 126) {
			flag = true
			continue
		}
	}
	if lenFlag && numFlag && azFlag && flag {
		return true
	}
	return false
}
