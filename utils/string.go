package utils

import (
	"bytes"
	"math/rand"
)

var BaseChars []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

//生成随机字符
func RandString(length int, chars []string) string {
	//rand.Seed(time.Now().UnixNano())
	charLen := len(chars)
	var buf bytes.Buffer
	for start := 0; start < length; start++ {
		t := rand.Intn(charLen)
		buf.WriteString(chars[t])
	}
	return buf.String()
}
