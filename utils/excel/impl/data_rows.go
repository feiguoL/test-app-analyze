package impl

import (
	"fmt"
)

// 首页概览数据
func (e *excelImpl) TotalSheet(sheetName, pkgName, pkgVersion, pkgArch,
	pkgSize, execNum, applib, outlib, rick, baseLineVersion string) (err error) {

	e.excelFile.SetColWidth(sheetName, "A", "Z", 50)
	e.excelFile.MergeCell(sheetName, "A1", "B1")
	e.excelFile.SetCellStr(sheetName, "A1", "APP兼容性分析报告")

	e.excelFile.SetCellStr(sheetName, "A2", "软件包名称")
	e.excelFile.SetCellStr(sheetName, "B2", pkgName)

	e.excelFile.SetCellStr(sheetName, "A3", "架构")
	e.excelFile.SetCellStr(sheetName, "B3", pkgArch)

	e.excelFile.SetCellStr(sheetName, "A4", "包版本")
	e.excelFile.SetCellStr(sheetName, "B4", pkgVersion)

	e.excelFile.SetCellStr(sheetName, "A5", "包大小(单位字节/byte)")
	e.excelFile.SetCellStr(sheetName, "B5", pkgSize)

	e.excelFile.SetCellStr(sheetName, "A6", "可执行文件数量")
	e.excelFile.SetCellStr(sheetName, "B6", execNum)

	e.excelFile.SetCellStr(sheetName, "A7", "本地库数量")
	e.excelFile.SetCellStr(sheetName, "B7", applib)

	e.excelFile.SetCellStr(sheetName, "A8", "外部库数量")
	e.excelFile.SetCellStr(sheetName, "B8", outlib)

	e.excelFile.SetCellStr(sheetName, "A9", "软件包是否存在兼容性风险")
	e.excelFile.SetCellStr(sheetName, "B9", rick)

	e.excelFile.SetCellStr(sheetName, "A10", "基线系统库版本")
	e.excelFile.SetCellStr(sheetName, "B10", baseLineVersion)
	e.excelFile.SetDefaultFont("Noto Sans")
	return
}

// 首行默认写入
func (e *excelImpl) FirstLineInited(sheetName string) (err error) {

	oneLine := map[string]string{
		"A1": "文件",
		"B1": "类型",
		"C1": "依赖外部库",
		"D1": "外部库的类型",
		"E1": "是否存在abi变化",
		"F1": "是否不在基线范围内",
		"G1": "详细信息",
	}
	e.excelFile.SetColWidth(sheetName, "A", "F", 18)
	e.excelFile.SetColWidth(sheetName, "G", "G", 100)
	SheetLines[sheetName] += 1
	for idx, line := range oneLine {
		err = e.excelFile.SetCellStr(sheetName, idx, line)
	}
	return err
}

// 在对应页签新增一条数据
func (e *excelImpl) AddRowInfo(sheetName string, args []string) (err error) {
	line := SheetLines[sheetName] + 1
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("A%d", line), args[0])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("B%d", line), args[1])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("C%d", line), args[2])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("D%d", line), args[3])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("E%d", line), args[4])
	if args[4] == "Yes" {
		style := e.GetRedColor()
		e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("E%d", line), fmt.Sprintf("E%d", line), style)
	}
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("F%d", line), args[5])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("G%d", line), args[6])
	SheetLines[sheetName] += 1
	return
}

func IntToString(i int) string {
	return string(byte(i))
}

// 更新是否有兼容性风险
func (e *excelImpl) UpdateRick(idx int, TotalName, sheetVersion, rick string) (err error) {
	char := IntToString(66 + idx)
	if rick == "Yes" {
		style := e.GetRedColor()
		e.excelFile.SetCellStyle(TotalName, fmt.Sprintf("%s9", char), fmt.Sprintf("%s10", char), style)
		err = e.excelFile.SetCellStr(TotalName, fmt.Sprintf("%s9", char), "是")
	} else {
		err = e.excelFile.SetCellStr(TotalName, fmt.Sprintf("%s9", char), "否")
	}
	err = e.excelFile.SetCellStr(TotalName, fmt.Sprintf("%s10", char), sheetVersion)
	return
}
