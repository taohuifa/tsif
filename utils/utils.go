package utils

import (
	"io/ioutil"
	"os"
	"strconv"
)

// 写入pid文件
func WritePIDFile(filePath string) error {
	return ioutil.WriteFile(filePath, []byte(strconv.Itoa(os.Getpid())), 0644)
}

// 数组合并
func ArrayAdd(array []interface{}, adds []interface{}) []interface{} {
	asize := len(adds)
	for i := 0; i < asize; i++ {
		array = append(array, adds[i])
	}
	return array
}
