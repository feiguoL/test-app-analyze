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
	Yes = "yes"
	No  = "no"
)

var (
	err         error  // 默认错误信息
	arch        string // 架构名称
	kernel      string // 内核名称
	system      string // 系统名称
	tempdir     string // 临时目录
	tempdirfile string // 临时目录文件
	applib      int    // app本地库文件数量
	outlib      int    // app外部库文件数量
	appexec     int    // app执行文件数量
	debname     string // app包名称
	debarch     string // app架构
	debversion  string // app版本
	debsize     int64  // app包大小
	soNum       map[string]uint8
	wg          sync.WaitGroup
	mutex       sync.Mutex
)

// 系统环境检查
func environment_check() (err error) {
	fmt.Println("environment checking...")
	arch = utils.SystemArch()
	kernel, err = utils.SystemkernelVersion()
	if err != nil {
		fmt.Println("get system kernel version failed: ", err)
	}
	system, err = utils.SystemVersion()
	if err != nil {
		fmt.Println("get system version failed: ", err)
	}
	fmt.Printf(`----------------------------------------
System Name:   | %s   
System Kernel: | %s   
System Arch:   | %s   
----------------------------------------
`, system, kernel, arch)
	if !common.ListContainersStr(arch, common.Archs) {
		return fmt.Errorf(common.ArchError)
	}
	fmt.Println("environment check finished!")
	return
}

// 准备工作目录
func ready_work_directory(fp string) (err error) {
	fmt.Println("ready work directory...")
	tempdir = fmt.Sprintf("%s%s/", common.WorkDir, utils.GetCurrentTimeStr())
	fmt.Printf("temp dir: %s\n", tempdir)
	if !utils.FileExist(tempdir) {
		if !utils.FileMkdir(tempdir) {
			return fmt.Errorf("%s dir: %s", common.MkdirError, tempdir)
		}
	}
	tempdirfile = fmt.Sprintf("%s%s", tempdir, path.Base(fp))
	if !utils.FileCopy(fp, tempdirfile) {
		return fmt.Errorf("%s source: %s target: %s", common.CopyError, fp, tempdirfile)
	}
	pkg, err := utils.PackageInfo(tempdirfile)
	if err != nil {
		return fmt.Errorf("%s file: %v", common.FileInfoError, err)
	}
	debname = pkg.Name
	debarch = pkg.Arch
	debversion = pkg.Version
	debsize = pkg.Size
	fmt.Printf(`----------------------------------------
Package Name: 		| %s
Package Arch: 		| %s
Package Version: 	| %s
Package Size: 		| %d
----------------------------------------
`, debname, debarch, debversion, int(debsize))
	fmt.Println("ready work directory finished!")
	return
}

// ldd获取文件信息
func ldd_file_handle(fp string) []string {
	outs, err := utils.LDDFile(fp)
	if err != nil {
		fmt.Println("ldd file error: ", err.Error())
		return outs
	}
	for _, val := range outs {
		fmt.Println(val)
		so := strings.Split(val, " ")
		soNum[so[0]] += 1
	}
	outlib = len(soNum)
	return outs
}

// excel的初始化总览数据
func excel_handler(appexec int) error {
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
	filesize := strconv.Itoa(int(debsize))
	execnum := strconv.Itoa(appexec)
	localnum := strconv.Itoa(applib)
	outnum := strconv.Itoa(outlib)
	err = excel.ExcelImpl.TotalSheet(common.TotalSheet, baseline, debname,
		debversion, debarch, filesize, execnum, localnum, outnum)

	return err
}

// 执行abi diff的详细操作
func diffabi(wg *sync.WaitGroup, baselinefile, locso string, outso []string) {
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
	switch debarch {
	case "i386":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.I386.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "amd64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.AMD64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, strings.TrimSpace(so))
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, strings.TrimSpace(out))
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "arm64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.ARM64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "sw_64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.SW_64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "loong64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.LOONG64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "loongarch64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.LOONGARCH64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "mips64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.MIPS64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	case "riscv64":
		for _, out := range outso {
			so := strings.Split(out, " ")[0]
			if common.ListContainersStr(so, config.RISCV64.Change) {
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
				notfound = Yes
			} else {
				notfound = No
			}
			var args []string
			args = append(args, locso)
			args = append(args, "本地so库")
			args = append(args, so)
			args = append(args, "外部so库")
			args = append(args, nochange)
			args = append(args, change)
			args = append(args, notfound)
			args = append(args, out)
			excel.ExcelImpl.AddRowInfo(sheetName, args)
		}
	}
}

// 执行并发多个基线文件处理的关联
func libdo(wg *sync.WaitGroup, locso string, outso []string) {
	defer wg.Done()
	baselinefiles, err := utils.FileWalkDir(common.BaseLineDir)
	if err != nil {
		fmt.Printf("Not Found Base Line File, Shell Script Stop Error: %s", err)
		return
	}
	subwg := sync.WaitGroup{}
	subwg.Add(len(baselinefiles))
	for _, baselinefile := range baselinefiles {
		baseline_sheet := strings.ReplaceAll(path.Base(baselinefile), ".json", "")
		excel.ExcelImpl.NewExcelSheet(baseline_sheet)
		go diffabi(&subwg, baselinefile, locso, outso)
	}
	subwg.Wait()

}

// 处理文件目录内容
func handle_file() error {
	fmt.Println("decompress deb file and app analysis handler...")
	libinfo := make(map[string][]string)
	err = utils.DecompressDeb(tempdirfile, tempdir)
	if err != nil {
		return err
	}
	fileList, err := utils.FindFilePath(tempdir, "*.so")
	if err != nil {
		return err
	}
	applib = len(fileList)
	soNum = make(map[string]uint8)
	for _, val := range fileList {
		so := path.Base(val)
		fmt.Println(so)
		outs := ldd_file_handle(val)
		libinfo[so] = outs
	}
	appexec, err := utils.FindDebExecutNum(tempdir)
	if err != nil {
		return err
	}
	fmt.Printf(`--------------------------------
App Lib File Number:   | %d
App Exec File Number:  | %d
App So File Number:    | %d
--------------------------------
`, applib, appexec, outlib)

	fmt.Println("app diff abi checking...")
	err = excel_handler(appexec)
	// 批量处理abi库比对
	wg := sync.WaitGroup{}
	wg.Add(applib)
	for k, v := range libinfo {
		go libdo(&wg, k, v)
	}
	wg.Wait()
	utils.FileClean(tempdir)

	return err
}

// @auth ut000198  (2024/06/18)
// @description 运行app兼容性的主服务逻辑
// @param fp软件包的对应目录
// @return 返回检查执行过程中的错误
func Appcheck(fp string) (err error) {
	fmt.Println(`ABI Check Compatibility Analysis Tool execution, Please wait for moment...`)
	err = environment_check()
	if err != nil {
		return
	}
	err = ready_work_directory(fp)
	if err != nil {
		return
	}
	handle_file()
	excel.ExcelImpl.Save("/home/uos/Downloads/test.xlsx")
	fmt.Println("App Analysis ABI Check Complete. The Report File Path: /home/uos/Downloads/test.xlsx")
	excel.ExcelImpl.Close()
	return
}
