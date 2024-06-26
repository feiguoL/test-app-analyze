package impl

// excel 统一的服务接口
type Service interface {
	NewExcelSheet(sheet string)
	TotalSheet(sheetName, pkgName, pkgVersion, pkgArch,
		pkgSize, execNum, applib, outlib, rick, baseline string) (err error)
	FirstLineInited(sheet string) (err error)
	AddRowInfo(sheetName string, args []string) (err error)
	UpdateRick(idx int, totalName, sheetVersion, rick string) (err error)
	SetActiveSheet(sheetIndex int)
	DeleteSheet(sheetName string)
	Save(filepath string) (err error)
	Close() error
}
