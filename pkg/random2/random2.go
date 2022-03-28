package random2

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
	"math/rand"
	"time"
)

const (
	upperLetters = 1
	lowerLetters = 2
	numbers      = 3
	symbols      = 4
)

//Numbers number 0-9
var Numbers = []byte{ /*48, 49,*/ 50, 51, 52, 53, 54, 55, 56, 57}

//UpperLetters A-Z
var UpperLetters = []byte{
	65, 66, 67, 68, 69, 70, 71, 72 /*73,*/, 74, 75, 76, 77, 78 /*79,*/, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
}

//LowerLetters a-z
var LowerLetters = []byte{
	97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107 /*108,*/, 109, 110 /*111,*/, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
}

//Symbols symbols
var Symbols = []byte{
	33, 34, 35, 36, 37, 38 /*, 39,40, 41*/, 42, 43, 44, 45, 46, /*47,*/
	/*58, 59, 60, 61, 62,*/ 63, 64,
	/*91, 92, 93,*/ 94, 95, 96,
	/*123, 124, 125,*/ 126,
}

//Illegal not allow symbols
var Illegal = []byte{
	39, 40, 41, 47, 58, 59, 60, 61, 62, 91, 92, 93, 23, 124, 125,
}

//RandomString random length（n） string，num is type（some number，UpperLetters，LowerLetters，symbols）
func RandomString(n int, num int64) string {
	m := getPwd(n, num)

	res := make([]byte, 0)
	rand.Seed(time.Now().UnixNano())
	for k := range m {
		switch k {
		case numbers:
			for i := 0; i < m[k]; i++ {
				rd := rand.Intn(len(Numbers))
				res = append(res, Numbers[rd])
			}
		case upperLetters:
			for i := 0; i < m[k]; i++ {
				rd := rand.Intn(len(UpperLetters))
				res = append(res, UpperLetters[rd])
			}
		case lowerLetters:
			for i := 0; i < m[k]; i++ {
				rd := rand.Intn(len(LowerLetters))
				res = append(res, LowerLetters[rd])
			}
		case symbols:
			for i := 0; i < m[k]; i++ {
				rd := rand.Intn(len(Symbols))
				res = append(res, Symbols[rd])
			}
		}
	}

	for i := len(res) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		res[i], res[num] = res[num], res[i]
	}

	return string(res)
}

func getPwd(n int, num int64) map[int]int {
	return getN(n, decodeRule(num))
}

func decodeRule(num int64) map[int]int {
	m := make(map[int]int)
	if num == 0 {
		//todo 采用默认全部规则
		m[numbers] = -1
		m[upperLetters] = -1
		m[lowerLetters] = -1
		m[symbols] = -1
		return m
	}
	sprintf := fmt.Sprintf("%b", num)

	for i := len(sprintf) - 1; i >= 0; i-- {
		if sprintf[i] == 49 {
			m[len(sprintf)-i] = -1
		}
	}
	return m
}

func getN(n int, m map[int]int) map[int]int {
	for k := range m {
		m[k] = -1
	}

	l := len(m)
	if n < l {
		n = l
	}
	if l == 1 {
		for k := range m {
			m[k] = n
		}
		return m
	}
	if n == l {
		for k := range m {
			m[k] = 1
		}
		return m
	}
	rand.Seed(time.Now().Unix())
	var j = 0
	var c = 0
	for k := range m {
		if j == 0 {
			b := rand.Intn(n / 2)
			if b == 0 {
				b = 1
			}
			if b > (n / 2) {
				getN(n, m)
			}
			if b > 0 && b < n {
				m[k] = b
				j = j + b
				c = c + 1
				continue
			}
			getN(n, m)
		}
		if l-c == 1 {
			if n-j > 0 {
				m[k] = n - j
			}
		} else {
			b := rand.Intn((n - j) / 2)
			if b == 0 {
				b = 1
			}
			if b > 0 && b < n-j {
				m[k] = b
				j = j + b
				c = c + 1
				if n-j < l-c {
					getN(n, m)
				}
				continue
			}
			getN(n, m)
		}

	}
	var i = 0
	for _, v := range m {
		i = i + v
	}
	if i != n {
		return getN(n, m)
	}
	return m
}

//CheckPassword check password rule
func CheckPassword(s string, minLength, num int64) bool {
	rule := decodeRule(num)
	runes := []rune(s)
	if int64(len(runes)) < minLength {
		return false
	}
	for _, v := range runes {
		fmt.Println(string(v))
		if byte(v) > 126 {
			return false
		}
		for _, v1 := range Illegal {
			if byte(v) == v1 {
				return false
			}
		}
	}
A:
	for k := range rule {
		switch k {
		case numbers:
			for _, v := range runes {
				for _, v1 := range Numbers {
					if byte(v) == v1 {
						rule[k] = 1
						continue A
					}
				}
			}
		case upperLetters:
			for _, v := range runes {
				for _, v1 := range UpperLetters {
					if byte(v) == v1 {
						rule[k] = 1
						continue A
					}
				}
			}
		case lowerLetters:
			for _, v := range runes {
				for _, v1 := range LowerLetters {
					if byte(v) == v1 {
						rule[k] = 1
						continue A
					}
				}
			}
		case symbols:
			for _, v := range runes {
				for _, v1 := range Symbols {
					if byte(v) == v1 {
						rule[k] = 1
						continue A
					}
				}
			}
		}
	}
	for _, v := range rule {
		if v != 1 {
			return false
		}
	}
	return true
}
