package utils

import (
	"bufio"
	"fmt"
	"os/exec"
)

var Debug = false

// @auth ut000198  (2024/06/17)
// @description 运行命令的函数方法
// @param name 执行的命令名称 args 执行命令的参数
// @return 返回执行的输出与错误
func RunCmd(name string, args []string, workdir string) ([]string, error) {
	var result []string
	cmd := exec.Command(name, args...)
	if workdir != "" {
		cmd.Dir = workdir
	}
	// 调试打印执行命令方便排查
	if Debug {
		fmt.Println(cmd.String())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// 读取系统输出的结果
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		ret := scanner.Text()
		if ret != "" {
			result = append(result, ret)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return result, err
}

// @auth ut000198  (2024/06/18)
// @description 查找文件目录路径
// @param fp 执行的目录路径 codition 查找的条件
// @return 返回找到结果的列表
func FindFilePath(fp, codition string) ([]string, error) {
	var args []string
	args = append(args, fp)
	args = append(args, "-name")
	args = append(args, codition)
	out, err := RunCmd("find", args, "")
	if err != nil {
		return nil, err
	}
	return out, err
}

// @auth ut000198  (2024/06/18)
// @description 解压deb文件
// @param fp 解压文件路径 dir 解压的目标目录
// @return 返回执行错误
func DecompressDeb(fp, dir string) error {
	outs, err := RunCmd("dpkg", []string{"-x", fp, dir}, "")
	if err != nil {
		fmt.Println(outs)
		return err
	}
	return err
}

// @auth ut000198  (2024/06/18)
// @description 获取外部库的数量信息
// @param fp 执行的目录路径
// @return 返回找到结果的列表
func LDDFile(fp string) ([]string, error) {
	out, err := RunCmd("ldd", []string{fp}, "")
	if err != nil {
		return nil, err
	}
	return out, err
}

// @auth ut000198  (2024/06/19)
// @description 获取app执行文件列表方法
// @param fp 执行的目录路径
// @return 返回找到结果的列表
func FindDebExecutNum(fp string) (int, error) {
	var args []string
	args = append(args, "-c")
	cmd := fmt.Sprintf(`find %s -type f -executable -exec file {} \; | grep "ELF 64-bit LSB executable"`, fp)
	args = append(args, cmd)
	outs, err := RunCmd("bash", args, fp)
	if err != nil {
		return 0, err
	}
	if Debug {
		for _, line := range outs {
			fmt.Println(line)
		}
	}
	return len(outs), err
}
