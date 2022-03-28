package main

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
	"flag"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"github.com/quanxiang-cloud/organizations/pkg/job/refactor/logic"
)

var (
	configPATH = flag.String("config", "configs/config.yml", "-config=配置文件地址")
)

func main() {
	flag.Parse()
	conf, err := configs.NewConfig(*configPATH)
	if err != nil {
		panic(err)
	}
	adaptedLogger := logger.Logger
	if err != nil {
		panic(err)
	}

	v1, err := logic.NewCleanV1(conf, adaptedLogger)
	if err != nil {
		panic(err)
	}
	err = v1.CleanV1()
	if err != nil {
		panic(err)
	}
}
