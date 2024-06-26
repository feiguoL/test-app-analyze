package impl

import (
	"github.com/xuri/excelize/v2"
)

var (
	SheetIndex map[string]int // 存储过程中产生sheet页值index
	SheetLines map[string]int // 存储过程中产生对应sheet页的行数
)

type excelImpl struct {
	excelFile *excelize.File // 接口方法主体
}

// 标识红色边框
func (e *excelImpl) GetRedColor() int {
	styleId, _ := e.excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "right",
				Color: "#FF0000",
				Style: 2,
			},
			{
				Type:  "left",
				Color: "#FF0000",
				Style: 2,
			},
			{
				Type:  "top",
				Color: "#FF0000",
				Style: 2,
			},
			{
				Type:  "bottom",
				Color: "#FF0000",
				Style: 2,
			},
		},
	})
	return styleId
}

// 初始化excel文件方法服务
func NewExcelFileImpl() Service {
	SheetIndex = make(map[string]int)
	SheetLines = make(map[string]int)
	return &excelImpl{
		excelFile: excelize.NewFile(),
	}
}
