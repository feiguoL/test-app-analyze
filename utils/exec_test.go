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
	excel.ExcelImpl.NewExcelSheet(total)
	err = excel.ExcelImpl.TotalSheet(total, "com.postman.postman_10.24_amd64.deb ",
		"10.24", "amd64", "92080298", "1", "4", "20", "是", baseLineVersion)

	excel.ExcelImpl.NewExcelSheet(baseLineVersion)
	excel.ExcelImpl.DeleteSheet("Sheet1")
	err = excel.ExcelImpl.FirstLineInited(baseLineVersion)
	if err != nil {
		t.Error(err.Error())
	}
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite.so", "本地so库", "libc.so.6", "外部so库", "No", "Yes",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite1.so", "本地so库", "libc.so.6", "外部so库", "Yes", "No",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite2.so", "本地so库", "libc.so.6", "外部so库", "No", "No",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite3.so", "本地so库", "libc.so.6", "外部so库", "Yes", "No",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite4.so", "本地so库", "libc.so.6", "外部so库", "No", "No",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite5.so", "本地so库", "libc.so.6", "外部so库", "No", "No",
		"libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd941a7d000)"})
	err = excel.ExcelImpl.AddRowInfo(baseLineVersion, []string{"dwrite6.so", "本地so库", "libc.so.6", "外部so库", "Yes", "No",
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
