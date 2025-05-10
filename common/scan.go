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
			Method:  "GET",
			Threads: opt.Threads,
			FCodes:  strings.Split(opt.FilterCode, ","), //需要过滤的状态码
		}

		// 生成字典拼接的url
		targest := spliceUrl(url, wordlist)

		//	输出日志信息

		results := httpx.Reset().Runner(url, targest)
		fmt.Println(results)
		//all.WriteString(results)
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
