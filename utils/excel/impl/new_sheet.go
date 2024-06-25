package impl

// 新建excel sheet标签页
func (e *excelImpl) NewExcelSheet(sheet string) {
	index := e.excelFile.NewSheet(sheet)
	SheetIndex[sheet] = index
	return
}
