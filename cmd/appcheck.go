/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"deepin-app-analyze/service"
	"deepin-app-analyze/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// @auth ut000198  (2024/06/17)
// @description app兼容性检查命令方法
var appcheckCmd = &cobra.Command{
	Use:   "appcheck",
	Short: "ABI Check Compatibility Analysis Tool",
	Long:  `ABI Check Compatibility Analysis Tool, Output Application Compatibility Analysis Report`,
	Run: func(cmd *cobra.Command, args []string) {
		fp, err := cmd.Flags().GetString("filepath")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		exp, err := cmd.Flags().GetString("example")
		// 调试模式，默认隐藏
		if exp == "debug" {
			utils.Debug = true
		}
		if exp == "test" {
			fmt.Println(`example: deepin-app-analyze -f ${filepath} | deepin-app-analyze --filepath ${filepath}
			deepin-app-analyze -f ./dde_2023.02.16_all.deb
			deepin-app-analyze --filepath /home/uos/Downloads/dde_2023.02.16_all.deb `)
			return
		}
		service.Appcheck(fp)
	},
}

// @auth ut000198  (2024/06/17)
// @description 初始化命令与参数
func init() {
	// 子命令使用的调用参数
	appcheckCmd.Flags().StringP("filepath", "f", "", "Path to deb package file")
	appcheckCmd.Flags().StringP("version", "v", "", "Default use the all version check or just the onces version")
	appcheckCmd.Flags().StringP("example", "e", "", "Use the example command default input 'test'")

	rootCmd.AddCommand(appcheckCmd)
}
