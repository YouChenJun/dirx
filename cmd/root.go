package cmd

import (
	"fmt"
	"os"

	"github.com/YouChenJun/dirx/libs"
	"github.com/YouChenJun/dirx/utils"
	"github.com/spf13/cobra"
)

var options = libs.Options{}

var RootCmd = &cobra.Command{
	Use:   libs.BINARY,
	Short: fmt.Sprintf("%s - %s", libs.BINARY, libs.DESC),
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&options.Target, "target", "u", "", "需要扫描的url目标")
	RootCmd.PersistentFlags().StringVarP(&options.TargetFile, "targets", "T", "", "需要扫描的url目标文件")
	RootCmd.PersistentFlags().IntVarP(&options.Threads, "threads", "t", 50, "单个任务的扫描线程数")
	RootCmd.PersistentFlags().IntVarP(&options.Concurrency, "concurrency", "c", 4, "并行任务数")
	RootCmd.PersistentFlags().StringVarP(&options.FilterCode, "fcode", "x", "400,404,406,416,501,502,503,302", "需要过滤的状态码")
	RootCmd.PersistentFlags().StringVarP(&options.Wordlist, "wordlist", "w", "", "字典文件路径")
	RootCmd.PersistentFlags().StringVarP(&options.Logfile, "log", "l", "", fmt.Sprintf("日志存储路径"))
	RootCmd.PersistentFlags().StringVarP(&options.Method, "method", "m", "GET", "扫描请求方法")
	RootCmd.PersistentFlags().IntVarP(&options.Timeout, "timeout", "s", 2, "单次请求超时时间/s")
	RootCmd.PersistentFlags().StringVarP(&options.OutPutFile, "outfile", "o", "", fmt.Sprintf("输出文件路径"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	//fmt.Println("initConfig")
	utils.InitLog(&options)
}
