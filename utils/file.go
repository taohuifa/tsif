package utils

import (
	"os"
)

// 文件夹是否存在
func FolderExist(filedir string) bool {
	_, err := os.Stat(filedir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 文件夹不存在则创建
func FolderNoCreate(filedir string) {
	if FolderExist(filedir) {
		return // 已经存在
	}
	os.Mkdir(filedir, os.ModePerm)
}

// 文件夹创建例子
// "path/filepath"
// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// if err != nil {
// 	log.Fatal(err)
// }
