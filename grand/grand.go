// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// Package grand provides high performance random bytes/number/string generation functionality.
package grand

import (
	"math/rand"
	"time"
)

var (
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32
	digits     = "0123456789"                                           // 10
	characters = letters + digits + symbols                             // 94
)

func init() {
	rand.Seed(time.Now().UnixNano()) //否则每次随机都一样
}

type Weight struct {
	Id     any
	Weight int
}

func NWeight(w []Weight) any {
	list := []any{}
	totalWeight := 0
	for _, weight := range w {
		totalWeight += weight.Weight
		for i := 1; i <= weight.Weight; i++ {
			list = append(list, weight.Id)
		}
	}
	if totalWeight == 0 {
		return nil
	}
	return list[N(1, totalWeight)-1]
}

// N returns a random int between min and max: [min, max].
// The <min> and <max> also support negative numbers.
func N(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		// Because Intn dose not support negative number,
		// so we should first shift the value to left,
		// then call Intn to produce the random number,
		// and finally shift the result back to right.
		return Intn(max-(min-0)+1) + (min - 0)
	}
	if min < 0 {
		// Because Intn dose not support negative number,
		// so we should first shift the value to right,
		// then call Intn to produce the random number,
		// and finally shift the result back to left.
		return Intn(max+(0-min)+1) - (0 - min)
	}
	return 0
}

//	func RandInt(min, max int) int64 {
//		if max==0{
//			return 0
//		}
//		rand.Seed(time.Now().UnixNano())
//		return int64(min) + rand.Int63n(int64(max+1-min))
//	}
func Intn(n int) int {
	return rand.Intn(n)
}

// B retrieves and returns random bytes of given length <n>.
//func B(n int) []byte {
//	if n <= 0 {
//		return nil
//	}
//	i := 0
//	b := make([]byte, n)
//	for {
//		copy(b[i:], <-bufferChan)
//		i += 4
//		if i >= n {
//			break
//		}
//	}
//	return b
//}

//
//// S returns a random string which contains digits and letters, and its length is <n>.
//// The optional parameter <symbols> specifies whether the result could contain symbols,
//// which is false in default.
//func S(n int, symbols ...bool) string {
//	if n <= 0 {
//		return ""
//	}
//	var (
//		b           = make([]byte, n)
//		numberBytes = B(n)
//	)
//	for i := range b {
//		if len(symbols) > 0 && symbols[0] {
//			b[i] = characters[numberBytes[i]%94]
//		} else {
//			b[i] = characters[numberBytes[i]%62]
//		}
//	}
//	return *(*string)(unsafe.Pointer(&b))
//}
//
//// Str randomly picks and returns <n> count of chars from given string <s>.
//// It also supports unicode string like Chinese/Russian/Japanese, etc.
//func Str(s string, n int) string {
//	if n <= 0 {
//		return ""
//	}
//	var (
//		b     = make([]rune, n)
//		runes = []rune(s)
//	)
//	if len(runes) <= 255 {
//		numberBytes := B(n)
//		for i := range b {
//			b[i] = runes[int(numberBytes[i])%len(runes)]
//		}
//	} else {
//		for i := range b {
//			b[i] = runes[Intn(len(runes))]
//		}
//	}
//	return string(b)
//}
//
//// Digits returns a random string which contains only digits, and its length is <n>.
//func Digits(n int) string {
//	if n <= 0 {
//		return ""
//	}
//	var (
//		b           = make([]byte, n)
//		numberBytes = B(n)
//	)
//	for i := range b {
//		b[i] = digits[numberBytes[i]%10]
//	}
//	return *(*string)(unsafe.Pointer(&b))
//}
//
//// Letters returns a random string which contains only letters, and its length is <n>.
//func Letters(n int) string {
//	if n <= 0 {
//		return ""
//	}
//	var (
//		b           = make([]byte, n)
//		numberBytes = B(n)
//	)
//	for i := range b {
//		b[i] = letters[numberBytes[i]%52]
//	}
//	return *(*string)(unsafe.Pointer(&b))
//}
//
//// Symbols returns a random string which contains only symbols, and its length is <n>.
//func Symbols(n int) string {
//	if n <= 0 {
//		return ""
//	}
//	var (
//		b           = make([]byte, n)
//		numberBytes = B(n)
//	)
//	for i := range b {
//		b[i] = symbols[numberBytes[i]%32]
//	}
//	return *(*string)(unsafe.Pointer(&b))
//}
//

func Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

//
//// Meet randomly calculate whether the given probability <num>/<total> is met.
//func Meet(num, total int) bool {
//	return Intn(total) < num
//}
//
//// MeetProb randomly calculate whether the given probability is met.
//func MeetProb(prob float32) bool {
//	return Intn(1e7) < int(prob*1e7)
//}
