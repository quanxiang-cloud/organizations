package department

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
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	"testing"
)

func TestTreeRouter(t *testing.T) {
	list := make([]org.Department, 0)
	list = append(list, org.Department{
		ID:   "1",
		PID:  "",
		Name: "Qingcloud",
	}, org.Department{
		ID:   "2",
		PID:  "1",
		Name: "dep",
	}, org.Department{
		ID:   "3",
		PID:  "1",
		Name: "研发",
	}, org.Department{
		ID:   "4",
		PID:  "3",
		Name: "超级部门",
	})
	router := NewDepartmentRouter()
	router.AddRoute(list)
	d := "Qingcloud/dep"
	d2 := "Qingcloud/超级部门"
	d3 := "Qingcloud/研发/超级部门"

	node1 := router.GetRoute(d)
	fmt.Println(node1)
	node2 := router.GetRoute(d2)
	fmt.Println(node2)
	node3 := router.GetRoute(d3)
	fmt.Println(node3)
}
