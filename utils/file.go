package utils

import (
	"os"
)

// 文件夹是否存在
func FolderExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
