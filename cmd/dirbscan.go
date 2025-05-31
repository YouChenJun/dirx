package cmd

import (
	"github.com/YouChenJun/dirx/common"
	"github.com/YouChenJun/dirx/utils"
	"github.com/spf13/cobra"
)

func init() {
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "作为目录扫描",
		RunE:  runScan,
	}
	RootCmd.AddCommand(scanCmd)
}

func runScan(_ *cobra.Command, _ []string) error {
	//如果输入的目标为文本列表，读取后加载
	if options.TargetFile != "" {
		if utils.FileExists(options.TargetFile) {
			options.Urls = append(options.Urls, utils.ReadingFileUnique(options.TargetFile)...)
		}
	}
	if options.Target != "" {
		options.Urls = append(options.Urls, options.Target)
	}

	if options.Wordlist == "" {
		utils.WarnF("请配置字典信息!")
		//当没输入字典信息时 自动加载字典
		return nil
	}
	wordList := utils.ReadingFileUnique(options.Wordlist)
	common.DirbScan(options.Urls, wordList, options)
	//common.DirbScan(options.Urls, []string{"11", "ch/milite/", "ccc"}, options)
	return nil
}
