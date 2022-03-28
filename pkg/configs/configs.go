package configs

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
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/cabin/tailormade/db/redis"
)

// DefaultPath default
var DefaultPath = "./configs/config.yml"

// DefaultOctopusPath default
var DefaultOctopusPath = "./configs/octopus.yml"

// Config config
type Config struct {
	MaxLoginErrNum   int              `yaml:"maxLoginErrNum"`
	LockAccountTime  time.Duration    `yaml:"lockAccountTime"`
	InternalNet      client.Config    `yaml:"internalNet"`
	ProcessPort      string           `yaml:"processPort"`
	Port             string           `yaml:"port"`
	Model            string           `yaml:"model"`
	OrgHost          string           `yaml:"orgHost"`
	TemplatePath     string           `yaml:"templatePath"`
	TemplateName     string           `yaml:"templateName"`
	POC              bool             `yaml:"poc"`
	Log              logger.Config    `yaml:"log"`
	Mysql            mysql.Config     `yaml:"mysql"`
	Redis            redis.Config     `yaml:"redis"`
	VerificationCode VerificationCode `yaml:"verificationCode"`
	MessageTemplate  MessageTemplate  `yaml:"messageTemplate"`
	Elastic          elastic.Config   `yaml:"elastic"`
	Ldap             Ldap             `yaml:"ldap"`
}

// Service service config
type Service struct {
	DB string `yaml:"db"`
}

// VerificationCode code
type VerificationCode struct {
	LoginCode    string        `yaml:"loginCode"`
	ResetCode    string        `yaml:"resetCode"`
	ForgetCode   string        `yaml:"forgetCode"`
	RegisterCode string        `yaml:"registerCode"`
	ExpireTime   time.Duration `yaml:"expireTime"`
}

// MessageTemplate message
type MessageTemplate struct {
	LoginCode    string `yaml:"loginCode"`
	ResetCode    string `yaml:"resetCode"`
	ForgetCode   string `yaml:"forgetCode"`
	RegisterCode string `yaml:"registerCode"`
	ResetPWD     string `yaml:"resetPWD"`
	NewPWD       string `yaml:"newPWD"`
}

// Ldap ldap
type Ldap struct {
	Open  bool   `yaml:"open"`
	Regex string `yaml:"regex"`
}

// NewConfig new
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewOctopusConfig new
func NewOctopusConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultOctopusPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
