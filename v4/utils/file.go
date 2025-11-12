package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileExists 检查文件是否存在
//
// 参数：
//   path - 文件路径
//
// 返回：
//   bool - true表示文件存在,false表示不存在
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

// ListFiles 递归列出目录下的所有文件
//
// 功能说明：
//   遍历目录及其子目录,收集所有符合条件的文件路径
//
// 参数：
//   path - 目录路径
//   pattern - 正则表达式模式(空字符串表示不过滤)
//   maxFiles - 最大文件数量(0表示不限制)
//
// 返回：
//   []string - 文件路径列表
//   error - 错误信息
func ListFiles(path string, pattern string, maxFiles int) ([]string, error) {
	var fileList []string

	// 读取目录
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return fileList, fmt.Errorf("读取目录失败: %w", err)
	}

	separator := string(filepath.Separator)

	for _, fi := range dir {
		// 如果是目录,递归遍历
		if fi.IsDir() {
			subFiles, err := ListFiles(filepath.Join(path, fi.Name()), pattern, maxFiles-len(fileList))
			if err == nil {
				for _, subFile := range subFiles {
					fileList = append(fileList, subFile)
					if maxFiles > 0 && len(fileList) >= maxFiles {
						return fileList, nil
					}
				}
			}
		} else {
			// 检查是否符合模式
			shouldAdd := true

			// 如果有正则模式,进行匹配
			if pattern != "" {
				matched, _ := regexp.MatchString(pattern, fi.Name())
				shouldAdd = matched
			} else {
				// 排除临时文件
				if strings.HasSuffix(fi.Name(), ".tmp") ||
					strings.HasSuffix(fi.Name(), ".temp") ||
					strings.HasSuffix(fi.Name(), ".dealing") {
					shouldAdd = false
				}
			}

			if shouldAdd {
				fileList = append(fileList, path+separator+fi.Name())
			}

			// 检查是否达到最大数量
			if maxFiles > 0 && len(fileList) >= maxFiles {
				return fileList, nil
			}
		}
	}

	return fileList, nil
}

// EnsureDir 确保目录存在
//
// 功能说明：
//   如果目录不存在则创建,如果已存在则不做任何操作
//
// 参数：
//   dir - 目录路径
//
// 返回：
//   error - 错误信息
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// RemoveDir 删除目录及其内容
//
// 参数：
//   dir - 目录路径
//
// 返回：
//   error - 错误信息
func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}
