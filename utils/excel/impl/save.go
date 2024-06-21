package impl

// 保存文件到目录路径
func (e *excelImpl) Save(filepath string) (err error) {
	return e.excelFile.SaveAs(filepath)
}

// 关闭文件上下文
func (e *excelImpl) Close() error {
	return e.excelFile.Close()
}
