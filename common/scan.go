package common

import (
	"fmt"
	"github.com/YouChenJun/dirx/common/httpx"
	"github.com/YouChenJun/dirx/libs"
	"github.com/YouChenJun/dirx/utils"
	"strings"
)

var all strings.Builder

func DirbScan(urls []string, wordlist []string, opt libs.Options) {
	for _, url := range urls {
		httpx := httpx.Httpx{
			Targets: make(chan string),
			Method:  opt.Method,
			Threads: opt.Threads,
			FCodes:  strings.Split(opt.FilterCode, ","), //需要过滤的状态码
			Timeout: opt.Timeout,
		}

		// 生成字典拼接好的url
		targets := spliceUrl(url, wordlist)
		utils.BlockF("Target", url)
		utils.TSPrintF("Method: %s | Threads: %d | Filter Code: %v | TimeOut: %v", opt.Method, opt.Threads, httpx.FCodes, httpx.Timeout)
		utils.InforF("扫描资产数:%v", len(urls))
		utils.InforF("扫描路径数:%v", len(targets))

		results := httpx.Reset().Runner(url, targets)
		for _, result := range results {
			fmt.Println(result["code"])
		}
	}
}

// spliceUrl 拼接扫描路径
func spliceUrl(url string, wordlist []string) []string {
	var targets []string

	// 遍历字典中的每个单词和扩展名，生成目标 URL
	for _, word := range wordlist {
		// 跳过包含 %EXT% 的字典行
		if strings.Contains(word, "%EXT%") {
			continue
		}
		// 去除 URL 末尾的斜杠（如果存在）
		trimmedURL := strings.TrimSuffix(url, "/")

		// 确保 word 以斜杠开头
		paddedWord := word
		if !strings.HasPrefix(word, "/") {
			paddedWord = "/" + word
		}
		targets = append(targets, trimmedURL+paddedWord)
	}
	return utils.RemoveDuplicateStrings(targets)
}
