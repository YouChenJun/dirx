package common

import (
	"encoding/json"
	"github.com/YouChenJun/dirx/libs"
	"github.com/YouChenJun/dirx/utils"
	"os"
)

// Filter 过滤扫描结果 302 301等情况
func Filter(result []utils.Result, opt libs.Options) {
	// 定义两个 map 来统计 Location 和 Size 的出现次数
	locationCount := make(map[string]int)
	sizeCount := make(map[string]int)

	// 第一次遍历：统计 Location 和 Size 的出现次数
	for _, r := range result {
		if r.Location != "" {
			locationCount[r.Location]++
		}
		if r.Size != "" {
			sizeCount[r.Size]++
		}
	}

	// 创建新的结果切片
	var filteredResults []utils.Result

	// 第二次遍历：过滤掉需要删除的记录
	for _, r := range result {
		// 检查是否满足删除条件
		if locationCount[r.Location] >= 10 || sizeCount[r.Size] >= 10 && r.Ctype != "application/octet-stream" {
			continue // 跳过需要删除的记录
		}
		filteredResults = append(filteredResults, r)
	}

	// 如果需要保存过滤后的结果，可以调用 SaveJsonData 函数
	if len(filteredResults) > 0 {
		//utils.InforF("%v [%v] %v %v [%v]", data["url"], data["code"], data["ctype"], data["location"], data["size"])
		err := SaveJsonData(filteredResults, opt)
		if err != nil {
			// 处理错误
			return
		}
	}
}

func SaveJsonData(result []utils.Result, opt libs.Options) error {
	// 以追加和写入模式打开文件
	file, err := os.OpenFile(opt.OutPutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 为每个结果单独编码并追加写入文件
	for _, r := range result {
		jsonData, _ := json.Marshal(r)
		if _, err := file.WriteString(string(jsonData) + "\n"); err != nil {
			return err
		}
		utils.InforF("%s [%s] %s %s [%s]", r.Url, r.Code, r.Ctype, r.Location, r.Size)
	}
	return nil
}
