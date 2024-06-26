package common

import "strings"

// 判断字符串是否在列表中
func ListContainersStr(s string, list []string) bool {
	for _, val := range list {
		if val == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}
