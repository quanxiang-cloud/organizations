package systems

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
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

const (
	systemRedis         = "orgs:systems:secret"
	expireTime          = 2 //hour
	passwordMinLength   = 8
	defaultPasswordRule = 15 //0x1111
)

//GetSecurityInfo get info
func GetSecurityInfo(ctx context.Context, conf configs.Config, redisClient redis.UniversalClient) *SecurityInfo {
	newSystems := NewSystems(conf.InternalNet)
	info, err := newSystems.GetSecurityInfo(ctx)
	securityInfo := SecurityInfo{}
	if err != nil {
		val := redisClient.Get(ctx, systemRedis).Val()
		if val != "" {
			json.Unmarshal([]byte(val), &securityInfo)
			return &securityInfo
		}
		securityInfo.PwdType = defaultPasswordRule
		securityInfo.PwdCount = int64(conf.MaxLoginErrNum)
		securityInfo.PwdCountWait = int64(conf.LockAccountTime)
		securityInfo.PwdMinLen = passwordMinLength
		return &securityInfo
	}
	if info == nil || info.ID == "" || info.PwdType == 0 || info.PwdMinLen == 0 {
		securityInfo.PwdType = defaultPasswordRule
		securityInfo.PwdCount = int64(conf.MaxLoginErrNum)
		securityInfo.PwdCountWait = int64(conf.LockAccountTime)
		securityInfo.PwdMinLen = passwordMinLength
		return &securityInfo
	}

	marshal, _ := json.Marshal(info)
	redisClient.SetEX(ctx, systemRedis, string(marshal), expireTime*time.Hour)
	return info

}
