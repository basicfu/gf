package ghex

import (
	"math"
	"strings"
)

func B16To62(b string) string {
	s, _ := Conversion(b, 16, 62)
	return s
}

func B62To16(b string) string {
	s, _ := Conversion(b, 62, 16)
	return s
}

var decToB64Map = map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M", 13: "N", 14: "O", 15: "P", 16: "Q", 17: "R", 18: "S", 19: "T", 20: "U", 21: "V", 22: "W", 23: "X", 24: "Y", 25: "Z", 26: "a", 27: "b", 28: "c", 29: "d", 30: "e", 31: "f", 32: "g", 33: "h", 34: "i", 35: "j", 36: "k", 37: "l", 38: "m", 39: "n", 40: "o", 41: "p", 42: "q", 43: "r", 44: "s", 45: "t", 46: "u", 47: "v", 48: "w", 49: "x", 50: "y", 51: "z", 52: "0", 53: "1", 54: "2", 55: "3", 56: "4", 57: "5", 58: "6", 59: "7",
	60: "8", 61: "9", 62: "+", 63: "/"}

// 64进制转10进制
func B64ToDec(b64 string) int {
	b64ToDecMap := mapKeyValueChange(decToB64Map)
	dec := 0
	b64Length := len(b64)
	for i := 0; i < b64Length; i++ {
		dec = dec + b64ToDecMap[string(b64[i])]*int(math.Pow(64, float64(b64Length-i-1)))
	}
	return dec
}

// 十进制转64进制
func DecToB64(dec int) string {
	b64 := ""
	for dec != 0 {
		b64 = decToB64Map[dec%64] + b64
		dec = dec / 64
	}
	return b64
}

// map结构键值互换
func mapKeyValueChange(data map[int]string) map[string]int {
	newData := make(map[string]int)
	for k, v := range data {
		newData[v] = k
	}
	return newData
}

func B32To64(str string) string {
	return DecToB64(BHex2Num(str, 36))
}
func B64To32(str string) string {
	return NumToBHex(B64ToDec(str), 36)
}

func DecToB36(n int) string {
	return NumToBHex(n, 36)
}
func B36ToDec(str string) int {
	return BHex2Num(str, 36)
}

// 36进制
var num2char = "0123456789abcdefghijklmnopqrstuvwxyz"

func NumToBHex(num, n int) string {
	num_str := ""
	for num != 0 {
		yu := num % n
		num_str = string(num2char[yu]) + num_str
		num = num / n
	}
	return strings.ToUpper(num_str)
}
func BHex2Num(str string, n int) int {
	str = strings.ToLower(str)
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int(v)
}
