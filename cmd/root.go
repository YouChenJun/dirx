package cmd

import (
	"fmt"
	"github.com/YouChenJun/dirx/libs"
	"github.com/spf13/cobra"
	"os"
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
	RootCmd.PersistentFlags().IntVarP(&options.Threads, "threads", "t", 20, "扫描线程数")
	RootCmd.PersistentFlags().StringVarP(&options.FilterCode, "fcode", "x", "400,404,406,416,501,502,503", "需要过滤的状态码 default:400,404,406,416,501,502,503")
	RootCmd.PersistentFlags().StringVarP(&options.Wordlist, "wordlist", "w", "", "字典文件路径")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	//fmt.Println("initConfig")

}
