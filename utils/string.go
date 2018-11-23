package utils

import (
	"bytes"
	"math/rand"
	"strings"
	// "unicode"
	// "fmt"
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

// 裁剪成map
func SplitMap(str string, sepA string, sepB string) *map[string]string {
	mmap := make(map[string]string)

	// 遍历第一层解析
	sstrs := strings.Split(str, sepA)
	slen := len(sstrs)
	for i := 0; i < slen; i++ {
		// fmt.Println(i, sstrs[i], strings.FieldsFunc(sstrs[i], unicode.IsSpace))
		// 解析第二层
		sargs := strings.Split(sstrs[i], sepB)
		alen := len(sargs)
		if alen <= 0 {
			mmap[sstrs[i]] = "" // 空值
			continue
		} else if alen == 1 {
			mmap[sargs[0]] = "" // 空值
			continue
		}
		mmap[sargs[0]] = sargs[1]
	}
	return &mmap
}
