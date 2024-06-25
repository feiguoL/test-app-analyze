package impl

// 设置当前表格的sheet页
func (e *excelImpl) SetActiveSheet(sheetIndex int) {
	e.excelFile.SetActiveSheet(sheetIndex)
}

func (e *excelImpl) DeleteSheet(sheetName string) {
	e.excelFile.DeleteSheet(sheetName)
}
