package core

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
	"reflect"
	"strconv"

	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
)

// data in or out
const (
	IN  = "IN"
	OUT = "OUT"
)

// Filter filter
func Filter(data interface{}, filter map[string]string, inOrOut string) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		dataFilter(data, filter, inOrOut)
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(data)
		for i := 0; i < of.Len(); i++ {
			if value := of.Index(i); value.IsValid() {
				Filter(of.Index(i).Interface(), filter, inOrOut)
			}
			continue
		}
	case reflect.Ptr:
		if reflect.ValueOf(data).IsValid() {
			Filter(reflect.ValueOf(data).Elem().Interface(), filter, inOrOut)
		}
	default:
	}
}

func dataFilter(oldSchema interface{}, filter map[string]string, inOrOut string) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			if _, ok := filter[iter.Key().String()]; ok {
				if !iter.Value().IsNil() {
					switch reflect.TypeOf(iter.Value().Interface()).Kind() {
					case reflect.Slice:
						switch filter[iter.Key().String()] {
						case consts.STRING, consts.TEXT:
							v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(fmt.Sprintf("%s", iter.Value().Interface())))
						case consts.TIME, consts.INT64:
							n, _ := strconv.ParseInt(fmt.Sprintf("%s", iter.Value().Interface()), 10, 64)
							v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(n))
						case consts.INT:
							n, _ := strconv.Atoi(fmt.Sprintf("%s", iter.Value().Interface()))
							v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(n))
						case consts.FLOAT:
							float, _ := strconv.ParseFloat(fmt.Sprintf("%s", iter.Value().Interface()), 64)
							v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(float))
						}
					}

					continue
				} else {
					// TODO delete
					v.SetMapIndex(iter.Key(), reflect.Value{})
				}
			} else {
				if inOrOut == IN {
					switch iter.Key().String() {
					case consts.ID:
						v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(fmt.Sprintf("%s", iter.Value().Interface())))
					default:
						// TODO delete
						v.SetMapIndex(iter.Key(), reflect.Value{})
					}
				} else {
					switch iter.Key().String() {
					case consts.ID, consts.NAME, consts.EMAIL, consts.SELFEMAIL, consts.PHONE, consts.AVATAR, consts.PASSWORDSTATUS:
						v.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(fmt.Sprintf("%s", iter.Value().Interface())))
					default:
						// TODO delete
						v.SetMapIndex(iter.Key(), reflect.Value{})
					}
				}
			}
		}
	default:
	}
}
