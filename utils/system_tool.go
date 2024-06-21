package utils

import (
	"fmt"
	"runtime"
	"strings"
)

// @auth ut000198  (2024/06/18)
// @description 获取系统架构信息
// @return 返回系统架构信息
func SystemArch() string {
	return runtime.GOARCH
}

// @auth ut000198  (2024/06/18)
// @description 获取系统内核信息
// @return 返回系统内核
func SystemkernelVersion() (string, error) {
	outs, err := RunCmd("uname", []string{"-r"}, "")
	if err != nil {
		return "", err
	}
	return outs[0], nil
}

// @auth ut000198  (2024/06/18)
// @description 获取系统当前版本
// @return 返回系统版本
func SystemVersion() (string, error) {
	outs, err := RunCmd("grep", []string{"PRETTY_NAME", "/etc/os-release"}, "")
	if err != nil {
		return "", err
	}
	if outs[0] != "" && strings.Contains(outs[0], "=") {
		return strings.Split(outs[0], "=")[1], err
	}
	return outs[0], fmt.Errorf("get system information failed")
}
