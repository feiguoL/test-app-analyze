package impl

// 新建excel sheet标签页
func (e *excelImpl) NewExcelSheet(sheet string) (err error) {
	index, err := e.excelFile.NewSheet(sheet)
	if err != nil {
		return err
	}

	SheetIndex[sheet] = index
	return
}
