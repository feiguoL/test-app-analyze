package utils_test

import (
	"deepin-app-analyze/common"
	"deepin-app-analyze/utils"
	"encoding/json"
	"path/filepath"
	"testing"
)

const (
	workdir    = "/tmp/test/"                                                // 测试目录
	appfile    = "/home/uos/Downloads/com.postman.postman_10.24_amd64.deb"   // 测试app文件
	configfile = "/home/uos/go/src/deepin-app-analyze/config/1050-1060.json" // 测试配置文件基线配置
	configdir  = "/home/uos/go/src/deepin-app-analyze/config/"               // 测试配置文件基线目录
)

// 测试文件是否存在的方法
func Test_FileExist(t *testing.T) {
	bol := utils.FileExist(workdir)
	if bol {
		t.Logf("file is exist")
	} else {
		t.Errorf("file is not exist")
	}
}

// 测试创建目录的方法
func Test_FileMkdir(t *testing.T) {
	bol := utils.FileMkdir(workdir)
	if bol {
		t.Logf("directory mkdir success")
	} else {
		t.Errorf("directory mkdir failed")
	}
}

// 测试文件复制的方法
func Test_FileCopy(t *testing.T) {
	bol := utils.FileCopy(appfile, workdir+filepath.Base(appfile))
	if bol {
		t.Logf("copy file to work dir success")
	} else {
		t.Errorf("copy file to work dir failed")
	}
}

// 测试文件读取
func Test_ReadFile(t *testing.T) {
	var config common.Config
	out, err := utils.FileRead(configfile)
	if err != nil {
		t.Errorf("Failed to read")
	}
	err = json.Unmarshal(out, &config)
	if err != nil {
		t.Errorf("Failed to unmarshal")
	}
	t.Log(config.AMD64.Change)
}

// 测试目录遍历方法
func Test_FileWalkDir(t *testing.T) {
	outs, err := utils.FileWalkDir(configdir)
	if err != nil {
		t.Errorf("Failed to read")
	}
	for _, val := range outs {
		t.Log(val)
	}
}

// 测试获取包信息
func Test_GetPackgeInfo(t *testing.T) {
	out, err := utils.PackageInfo(appfile)
	if err != nil {
		t.Errorf("Failed to get package info: %v", err)
	}
	t.Log(out.Name, out.Version, out.Arch, out.Size)
}
