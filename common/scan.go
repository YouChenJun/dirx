package common

import (
	"strings"
	"sync"

	"github.com/YouChenJun/dirx/common/httpx"
	"github.com/YouChenJun/dirx/libs"
	"github.com/YouChenJun/dirx/utils"
)

var all strings.Builder

func DirbScan(urls []string, wordlist []string, opt libs.Options) {
	utils.InforF("扫描资产数:%v", len(urls))
	utils.InforF("并行任务数:%v | 单任务线程数:%v", opt.Concurrency, opt.Threads)

	// 创建任务通道和等待组
	urlChan := make(chan string, len(urls))
	var wg sync.WaitGroup

	// 启动并行任务
	for i := 0; i < opt.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for url := range urlChan {
				scanSingleTarget(url, wordlist, opt)
			}
		}(i)
	}

	// 发送所有URL到通道
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// 等待所有任务完成
	wg.Wait()
}

// scanSingleTarget 扫描单个目标
func scanSingleTarget(url string, wordlist []string, opt libs.Options) {
	// 分割并去除状态码的空格
	fcodes := strings.Split(opt.FilterCode, ",")
	for i, code := range fcodes {
		fcodes[i] = strings.TrimSpace(code)
	}
	httpx := httpx.Httpx{
		Targets:    make(chan string),
		Method:     opt.Method,
		Threads:    opt.Threads,
		FCodes:     fcodes, //需要过滤的状态码
		Timeout:    opt.Timeout,
		MaxRespone: 1024 * 1024 * 10,
	}

	// 生成字典拼接好的url
	targets := spliceUrl(url, wordlist)
	utils.BlockF("Target", url)
	utils.InforF("扫描路径数:%v", len(targets))
	utils.TSPrintF("Method: %s | Threads: %d | Filter Code: %v | TimeOut: %v", opt.Method, opt.Threads, httpx.FCodes, httpx.Timeout)
	datas := httpx.Reset().Runner(url, targets)
	var results []utils.Result

	for _, data := range datas {
		res := utils.Result{
			Url:      data["url"],
			Code:     data["code"],
			Location: data["location"],
			Ctype:    data["ctype"],
			Server:   data["server"],
			Status:   data["status"],
			Size:     data["size"],
			Body:     data["body"],
			Time:     data["time"],
		}
		results = append(results, res)
	}
	Filter(results, opt)
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
