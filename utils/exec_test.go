package utils_test

import (
	"deepin-app-analyze/utils/excel"
	"testing"
)

var err error

const (
	total           = "总览"
	baseLineVersion = "1050-1060"
)

// 测试excel基本功能
func Test_excelinit(t *testing.T) {
	err = excel.ExcelImpl.NewExcelSheet(total)
	if err != nil {
		t.Error(err.Error())
	}

	err = excel.ExcelImpl.TotalSheet(total, baseLineVersion, "com.postman.postman_10.24_amd64.deb ",
		"10.24", "amd64", "92080298", "1", "4", "20")

	err = excel.ExcelImpl.NewExcelSheet(baseLineVersion)
	if err != nil {
		t.Error(err.Error())
	}
	err = excel.ExcelImpl.DeleteSheet("Sheet1")
	if err != nil {
		t.Error(err.Error())
	}
	err = excel.ExcelImpl.FirstLineInited(baseLineVersion)
	if err != nil {
		t.Error(err.Error())
	}
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite.so", "本地so库", "libc.so.6", "外部so库", "是", "否", "否",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite1.so", "本地so库", "libc.so.6", "外部so库", "否", "是", "否",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite2.so", "本地so库", "libc.so.6", "外部so库", "否", "否", "是",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.Save("/home/uos/Downloads/test.xlsx")
	if err != nil {
		t.Error(err.Error())
	}
	if err != nil {
		t.Error(err.Error())
	}
	err = excel.ExcelImpl.Close()
	if err != nil {
		t.Error(err.Error())
	}
}
