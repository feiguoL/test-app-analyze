package utils_test

import (
	"deepin-app-analyze/utils"
	"testing"
)

// 测试获取系统CPU架构
func Test_GetSystemArch(t *testing.T) {
	arch := utils.SystemArch()
	t.Log(arch)
}

// 测试获取系统内核信息
func Test_GetSystemkernelVersion(t *testing.T) {
	version, err := utils.SystemkernelVersion()
	if err != nil {
		t.Errorf("获取系统内核信息错误：%s", err.Error())
	}
	t.Log(version)
}

// 测试获取系统版本
func Test_GetSystemVersion(t *testing.T) {
	version, err := utils.SystemVersion()
	if err != nil {
		t.Errorf("获取系统版本信息错误：%s", err.Error())
	}
	t.Log(version)
}
