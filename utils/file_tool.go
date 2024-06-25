package utils

import (
	"deepin-app-analyze/common"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// @auth ut000198  (2024/06/17)
// @description 判断文件目录是否存在
// @param fp文件路径
// @return 返回判断结果 true or false
func FileExist(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// @auth ut000198  (2024/06/17)
// @description 复制文件从源目录到目标目录
// @param src源文件路径 dst目标目录路径
// @return 返回判断结果 true or false
func FileCopy(src, dst string) bool {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return false
	} else {
		err = ioutil.WriteFile(dst, input, 0644)
		if err != nil {
			return false
		} else {
			return true
		}
	}
}

// @auth ut000198  (2024/06/17)
// @description 创建文件目录
// @param fp文件目录路径
// @return 返回判断结果 true or false
func FileMkdir(fp string) bool {
	err := os.MkdirAll(fp, 0777)
	if err != nil {
		return false
	} else {
		return true
	}
}

// @auth ut000198  (2024/06/18)
// @description 文件内容读取
// @param fp文件目录路径
// @return 返回文件内容
func FileRead(fp string) ([]byte, error) {
	info, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer info.Close()
	data, err := ioutil.ReadAll(info)
	if err != nil {
		return nil, err
	}

	return data, err
}

// @auth ut000198  (2024/06/18)
// @description 基线文件目录遍历
// @param dir文件目录路径
// @return 文件列表
func FileWalkDir(dir string) ([]string, error) {
	var data []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.Contains(path, ".json") {
				data = append(data, path)
			}
		}
		return nil
	})
	return data, err
}

func FileClean(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// @auth ut000198  (2024/06/18)
// @description 查询deb文件的包信息
// @param path 文件目录路径
// @return 包的基本信息
func PackageInfo(path string) (*common.PackageInfo, error) {
	var pkg common.PackageInfo
	fileinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	pkg.Size = fileinfo.Size()
	pkg.Name = fileinfo.Name()
	if pkg.Name != "" {
		names := strings.Split(pkg.Name, "_")
		pkg.Version = names[1]
		pkg.Arch = strings.ReplaceAll(names[2], ".deb", "")
	}

	return &pkg, nil
}
