package utils_test

import (
	"deepin-app-analyze/utils"
	"fmt"
	"path/filepath"
	"testing"
)

// 测试运行命令的方法函数
func Test_RunCmd(t *testing.T) {
	// TODO: 执行解压软件包命令是否执行成功
	out, err := utils.RunCmd("dpkg", []string{"-x", fmt.Sprintf("%s%s", workdir, filepath.Base(appfile)), workdir}, "")
	if err != nil {
		t.Errorf("执行该命令错误：%s", err.Error())
	}

	t.Log(out)

}

// 查找文件列表
func Test_FindFilePath(t *testing.T) {
	outs, err := utils.FindFilePath(workdir, "*.so")
	if err != nil {
		t.Errorf("查找文件执行失败: %s", err.Error())
	}
	for _, val := range outs {
		fmt.Println(val)
	}
	t.Logf("查找成功")
}

// 测试获取外部库文件信息
func Test_LDDFile(t *testing.T) {
	fp := fmt.Sprintf("%s%s", workdir, "/opt/apps/com.postman.postman/files/Postman/app/libEGL.so")
	outs, err := utils.LDDFile(fp)
	if err != nil {
		t.Errorf("获取外部库执行失败: %s", err.Error())
	}
	for _, val := range outs {
		fmt.Println(val)
	}
	t.Logf("获取成功")

}

// 测试解压deb文件方法
func Test_DecompressDeb(t *testing.T) {
	err := utils.DecompressDeb(appfile, workdir)
	if err != nil {
		t.Errorf("解压deb文件失败: %s", err.Error())
	}
	t.Logf("解压deb文件成功")
}

// 测试查找deb执行文件数量
func Test_FindDebExecutNum(t *testing.T) {
	num, err := utils.FindDebExecutNum(workdir)
	if err != nil {
		t.Errorf("Failed to find deb execut number: %v", err)
	}
	t.Log(num)
}
