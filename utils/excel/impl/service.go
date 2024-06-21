package impl

// excel 统一的服务接口
type Service interface {
	NewExcelSheet(sheet string) (err error)
	TotalSheet(sheetName, baseLineVersion, pkgName, pkgVersion, pkgArch,
		pkgSize, execNum, applib, outlib string) (err error)
	FirstLineInited(sheet string) (err error)
	AddRowInfo(sheetName string, args []string) (err error)
	SetActiveSheet(sheetIndex int)
	DeleteSheet(sheetName string) (err error)
	Save(filepath string) (err error)
	Close() error
}
