package excel

import "deepin-app-analyze/utils/excel/impl"

// excel 接口关联
var ExcelImpl impl.Service

// excel 方法初始化
func init() {
	ExcelImpl = impl.NewExcelFileImpl()
}
