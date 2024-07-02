package service

import (
	"deepin-app-analyze/common"
	"deepin-app-analyze/utils"
	"deepin-app-analyze/utils/excel"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	Yes        = "Yes"
	No         = "No"
	Null       = "/"
	LocalSo    = "本地so库"
	LocalFile  = "本地文件"
	OutSo      = "外部库"
	OutLocalSo = "本地库"
)

var (
	err error // 默认错误信息
	// 全局环境变量
	GlobalBaseLine    []string            // 指定基线版本
	GlobalTempDir     string              // 临时目录
	GlobalTempDirFile string              // 临时目录文件
	GlobalArch        string              // 架构名称
	GlobalKernel      string              // 内核名称
	GlobalSystem      string              // 系统名称
	GlobalAppLib      int                 // app本地库文件数量
	GlobalOutLib      int                 // app外部库文件数量
	GlobalAppExe      int                 // app执行文件数量
	GlobalDebName     string              // app包名称
	GlobalDebArch     string              // app架构
	GlobalDebVersion  string              // app版本
	GlobalDebSize     int64               // app包大小
	GlobalRick        string              // app兼容性风险
	GlobalSoNum       map[string]uint8    // 库文件数量
	GlobalRickMap     map[string]string   // 风险性与基线版本关联
	GlobalLibInfos    map[string][]string // 文件与外部库的关联
)

func init() {
	GlobalSoNum = make(map[string]uint8)
	GlobalRickMap = make(map[string]string)
	GlobalLibInfos = make(map[string][]string)
}

// 系统环境检查
func environment_check() (err error) {
	fmt.Println("environment checking...")
	GlobalArch = utils.SystemArch()
	GlobalKernel, err = utils.SystemkernelVersion()
	if err != nil {
		fmt.Println("get system kernel version failed: ", err)
	}
	GlobalSystem, err = utils.SystemVersion()
	if err != nil {
		fmt.Println("get system version failed: ", err)
	}
	fmt.Printf(`----------------------------------------
System Name:   | %s   
System Kernel: | %s   
System Arch:   | %s   
----------------------------------------
`, GlobalSystem, GlobalKernel, GlobalArch)
	if !common.ListContainersStr(GlobalArch, common.Archs) {
		return fmt.Errorf(common.ArchError)
	}
	fmt.Println("environment check finished!")
	return
}

// 准备工作目录
func ready_work_directory(fp string) (err error) {
	fmt.Println("ready work directory...")
	GlobalTempDir = fmt.Sprintf("%s%s/", common.WorkDir, utils.GetCurrentTimeStr())
	// fmt.Printf("temp dir: %s\n", GlobalTempDir)
	if !utils.FileExist(GlobalTempDir) {
		if !utils.FileMkdir(GlobalTempDir) {
			return fmt.Errorf("%s dir: %s", common.MkdirError, GlobalTempDir)
		}
	}
	GlobalTempDirFile = fmt.Sprintf("%s%s", GlobalTempDir, path.Base(fp))
	if !utils.FileCopy(fp, GlobalTempDirFile) {
		return fmt.Errorf("%s source: %s target: %s", common.CopyError, fp, GlobalTempDirFile)
	}
	pkg, err := utils.PackageInfo(GlobalTempDirFile)
	if err != nil {
		return fmt.Errorf("%s file: %v", common.FileInfoError, err)
	}
	GlobalDebName = pkg.Name
	GlobalDebArch = pkg.Arch
	GlobalDebVersion = pkg.Version
	GlobalDebSize = pkg.Size
	fmt.Printf(`----------------------------------------
Package Name: 		| %s
Package Arch: 		| %s
Package Version: 	| %s
Package Size: 		| %d
----------------------------------------
`, GlobalDebName, GlobalDebArch, GlobalDebVersion, int(GlobalDebSize))
	fmt.Println("ready work directory finished!")
	return
}

// ldd获取文件信息
func ldd_file_handle(fp string) []string {
	outs, err := utils.LDDFile(fp)
	if err != nil {
		fmt.Printf("ldd file error: %s %v", fp, err.Error())
		return outs
	}
	for _, val := range outs {
		// fmt.Println(val)
		so := strings.Split(val, " ")
		GlobalSoNum[so[0]] += 1
	}
	GlobalOutLib = len(GlobalSoNum)
	return outs
}

func Get_baseline_files() (string, error) {
	var result string
	baseLineFiles, err := utils.FileWalkDir(common.BaseLineDir)
	if err != nil {
		return "", err
	}
	if len(baseLineFiles) == 0 {
		return "", fmt.Errorf("no baseline files found")
	}
	for _, filePath := range baseLineFiles {
		result += fmt.Sprintf("    %s\n", strings.ReplaceAll(path.Base(filePath), ".json", ""))
	}
	return result, nil
}

// excel的初始化总览数据
func excel_handler(GlobalAppExe int) error {
	var baseline string
	excel.ExcelImpl.NewExcelSheet(common.TotalSheet)
	excel.ExcelImpl.DeleteSheet("Sheet1")

	baseLineFiles, err := utils.FileWalkDir(common.BaseLineDir)
	if err != nil {
		return err
	}
	for _, val := range baseLineFiles {
		baseline += strings.ReplaceAll(path.Base(val), ".json", "") + " "
	}
	filesize := strconv.Itoa(int(GlobalDebSize))
	execnum := strconv.Itoa(GlobalAppExe)
	localnum := strconv.Itoa(GlobalAppLib)
	outnum := strconv.Itoa(GlobalOutLib)
	GlobalRick = No
	err = excel.ExcelImpl.TotalSheet(common.TotalSheet, GlobalDebName,
		GlobalDebVersion, GlobalDebArch, filesize, execnum, localnum, outnum, GlobalRick, baseline)

	return err
}

// 执行abi diff的详细操作
func diffabi(wg *sync.WaitGroup, mutex *sync.Mutex, baselinefile, locso string, outso []string) {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()
	var config common.Config
	var change, nochange, notfound string
	sheetName := strings.ReplaceAll(path.Base(baselinefile), ".json", "")
	_ = excel.ExcelImpl.FirstLineInited(sheetName)
	out, err := utils.FileRead(baselinefile)
	if err != nil {
		fmt.Println("Failed to read config file")
	}
	err = json.Unmarshal(out, &config)
	if err != nil {
		fmt.Println("Failed to unmarshal config file")
	}
	// fmt.Println("zzzzz", GlobalRickMap)
	// 对应CPU架构的处理
	switch GlobalDebArch {
	case "i386":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.I386.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.I386.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "amd64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.AMD64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.AMD64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "arm64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.ARM64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.ARM64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "sw_64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.SW_64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.SW_64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "loong64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.LOONG64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.LOONG64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "loongarch64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.LOONGARCH64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.LOONGARCH64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "mips64el":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.MIPS64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.MIPS64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "riscv64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.RISCV64.Change) {
				GlobalRick = Yes
				GlobalRickMap[sheetName] = Yes
				change = Yes
			} else {
				change = No
			}
			if common.ListContainersStr(so, config.RISCV64.NoChange) {
				nochange = Yes
			} else {
				nochange = No
			}
			if change == No && nochange == No {
				notfound = No
			} else {
				notfound = Yes
			}
			var args []string
			args = append(args, locso)
			if strings.Contains(locso, ".so") {
				args = append(args, LocalSo)
			} else {
				args = append(args, LocalFile)
			}
			args = append(args, strings.TrimSpace(so))
			out = strings.TrimSpace(out)
			if strings.Contains(out, "/lib/") {
				args = append(args, OutSo)
			} else {
				args = append(args, OutLocalSo)
			}
			if strings.Contains(out, "=> not found") {
				out = strings.ReplaceAll(out, " => not found", "")
				change = Null
				notfound = Null
			}
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	}
}

// 执行并发多个基线文件处理的关联
func libdo(baselinefile string) {
	subwg := sync.WaitGroup{}
	subwg.Add(len(GlobalLibInfos))
	submutex := sync.Mutex{}
	for locso, outso := range GlobalLibInfos {
		go diffabi(&subwg, &submutex, baselinefile, locso, outso)
	}
	subwg.Wait()
}

// 处理文件目录内容
func handle_file() error {
	fmt.Println("decompress deb file and app analysis handler...")

	err = utils.DecompressDeb(GlobalTempDirFile, GlobalTempDir)
	if err != nil {
		fmt.Printf("error decompressing %v", err)
		return err
	}
	fileList, err := utils.FindFilePath(GlobalTempDir, "*.so")
	if err != nil {
		fmt.Printf("error finding file list %v", err)
		return err
	}
	GlobalAppLib = len(fileList)

	for _, val := range fileList {
		so := path.Base(val)
		// fmt.Println(so)
		outs := ldd_file_handle(val)
		GlobalLibInfos[so] = outs
	}
	GlobalAppExe, _ := utils.FindDebExecutFile(GlobalTempDir)
	for _, val := range GlobalAppExe {
		so := path.Base(val)
		outs := ldd_file_handle(val)
		GlobalLibInfos[so] = outs
	}
	fmt.Printf(`----------------------------------------
App Lib File Number:   | %d
App Exec File Number:  | %d
App So File Number:    | %d
----------------------------------------
`, GlobalAppLib, len(GlobalAppExe), GlobalOutLib)

	fmt.Println("app diff abi checking...")

	err = excel_handler(len(GlobalAppExe))
	baselinefiles, err := utils.FileWalkDir(common.BaseLineDir)
	if err != nil {
		fmt.Printf("Not Found Base Line File, Shell Script Stop Error: %s", err)
		return err
	}

	// 批量处理abi库比对
	for _, baselinefile := range baselinefiles {
		baseline_sheet := strings.ReplaceAll(path.Base(baselinefile), ".json", "")
		if len(GlobalBaseLine) > 0 {
			if common.ListContainersStr(baseline_sheet, GlobalBaseLine) {
				GlobalRickMap[baseline_sheet] = No
				excel.ExcelImpl.NewExcelSheet(baseline_sheet)
				fmt.Println("check baseline version: ", baseline_sheet)
				libdo(baselinefile)
			} else {
				err = fmt.Errorf("Not Found baseline version, Please recheck your input version !!!")
				goto OUTRUN
			}
		} else {
			GlobalRickMap[baseline_sheet] = No
			excel.ExcelImpl.NewExcelSheet(baseline_sheet)
			fmt.Println("check baseline version: ", baseline_sheet)
			libdo(baselinefile)
		}
	}
OUTRUN:
	utils.FileClean(GlobalTempDir)

	return err
}

// @auth ut000198  (2024/06/18)
// @description 运行app兼容性的主服务逻辑
// @param fp软件包的对应目录
// @return 返回检查执行过程中的错误
func Appcheck(fp string, baseline_list []string) (err error) {
	GlobalBaseLine = baseline_list
	defer excel.ExcelImpl.Close()
	fmt.Println(`ABI Check Compatibility Analysis Tool execution, Please wait for moment...`)
	err = environment_check()
	if err != nil {
		return
	}
	err = ready_work_directory(fp)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}
	err = handle_file()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}
	path := utils.GetPWD()
	rick := 0
	// 兼容性风险最终结果更新
	for k, v := range GlobalRickMap {
		excel.ExcelImpl.UpdateRick(rick, common.TotalSheet, k, v)
		rick += 1
	}
	name := strings.ReplaceAll(GlobalDebName, ".deb", "")
	excel.ExcelImpl.Save(fmt.Sprintf("%s/%s.xlsx", path, name))
	fmt.Println(fmt.Sprintf("App Analysis ABI Check Complete. The Report File Path: %s/%s.xlsx", path, name))
	return
}
