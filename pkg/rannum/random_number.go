package rannum

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
	"math/rand"
	"strconv"
	"sync"
)

/*RandomNumber random number
length: length
*/
func (s *Set) RandomNumber(length int) int64 {
	digitNumber := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	for {
		ranNumber := ""
		for j := 1; j < length; j++ {
			ranNumber += digitNumber[rand.Intn(len(digitNumber))]
		}
		//0开头的丢弃
		if ranNumber[:1] == "0" {
			continue
		}
		atoi, err := strconv.Atoi(ranNumber)
		if err != nil {
			continue
		}
		res := int64(atoi)
		if !s.Has(res) {
			s.Add(res)
			return res
		}
	}

}

//Set set
type Set struct {
	m map[int64]bool
	sync.RWMutex
}

//New new
func New() *Set {
	return &Set{
		m: map[int64]bool{},
	}
}

//Add add
func (s *Set) Add(item int64) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

//Remove remove
func (s *Set) Remove(item int64) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

//Has has
func (s *Set) Has(item int64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

//Len len
func (s *Set) Len() int {
	return len(s.List())
}

//Clear clear
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int64]bool{}
}

//IsEmpty is empty
func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

//List list
func (s *Set) List() []int64 {
	s.RLock()
	defer s.RUnlock()
	list := make([]int64, 0)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
