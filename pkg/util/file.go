package util

import (
	"path/filepath"
	"strings"
)

func GetFileNameAndExt(filename string) (string, string) {
	// 获取文件名（带后缀）
	fileNameWithExt := filepath.Base(filename)

	// 获取文件后缀（包含 .）
	ext := filepath.Ext(fileNameWithExt)

	// 获取文件名（去掉后缀）
	fileName := strings.TrimSuffix(fileNameWithExt, ext)

	return fileName, ext
}
