package impl

import (
	"fmt"
)

// 首页概览数据
func (e *excelImpl) TotalSheet(sheetName, baseLineVersion, pkgName, pkgVersion, pkgArch,
	pkgSize, execNum, applib, outlib string) (err error) {

	e.excelFile.SetColWidth(sheetName, "A", "B", 50)
	e.excelFile.MergeCell(sheetName, "A1", "B1")
	e.excelFile.SetCellStr(sheetName, "A1", "APP分析报告")

	e.excelFile.SetCellStr(sheetName, "A2", "基线系统库版本")
	e.excelFile.SetCellStr(sheetName, "B2", baseLineVersion)

	e.excelFile.SetCellStr(sheetName, "A3", "Deb包名称")
	e.excelFile.SetCellStr(sheetName, "B3", pkgName)

	e.excelFile.SetCellStr(sheetName, "A4", "架构")
	e.excelFile.SetCellStr(sheetName, "B4", pkgArch)

	e.excelFile.SetCellStr(sheetName, "A5", "包版本")
	e.excelFile.SetCellStr(sheetName, "B5", pkgVersion)

	e.excelFile.SetCellStr(sheetName, "A6", "包大小(单位字节/byte)")
	e.excelFile.SetCellStr(sheetName, "B6", pkgSize)

	e.excelFile.SetCellStr(sheetName, "A7", "可执行文件数量")
	e.excelFile.SetCellStr(sheetName, "B7", execNum)

	e.excelFile.SetCellStr(sheetName, "A8", "本地库数量")
	e.excelFile.SetCellStr(sheetName, "B8", applib)

	e.excelFile.SetCellStr(sheetName, "A9", "外部库数量")
	e.excelFile.SetCellStr(sheetName, "B9", outlib)

	return
}

// 首行默认写入
func (e *excelImpl) FirstLineInited(sheetName string) (err error) {

	oneLine := map[string]string{
		"A1": "文件",
		"B1": "类型",
		"C1": "依赖外部库",
		"D1": "外部库的类型",
		"E1": "无abi变更",
		"F1": "已知abi变更",
		"G1": "未知abi变更",
		"H1": "详细信息",
	}
	e.excelFile.SetColWidth(sheetName, "A", "G", 18)
	e.excelFile.SetColWidth(sheetName, "H", "H", 100)
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
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("F%d", line), args[5])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("G%d", line), args[6])
	err = e.excelFile.SetCellStr(sheetName, fmt.Sprintf("H%d", line), args[7])
	SheetLines[sheetName] += 1
	return
}
