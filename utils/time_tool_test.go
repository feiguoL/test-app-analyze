package utils_test

import (
	"deepin-app-analyze/utils"
	"testing"
)

func Test_GetFastUuid(t *testing.T) {
	ret := utils.GetFastUuid()
	t.Logf(ret)
}

func Test_GetCurrentTimeStr(t *testing.T) {
	ret := utils.GetCurrentTimeStr()
	t.Logf(ret)
}
