package common

import (
	"encoding/json"
	"os"

	"github.com/YouChenJun/dirx/libs"
	"github.com/YouChenJun/dirx/utils"
)

// Filter 过滤扫描结果：基于统计信息过滤重复的Location和Size
// 注意：状态码过滤（如302、404等）已在httpx.filter中处理
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
		// 如果 Location 或 Size 出现频率 >= 10 次，且不是二进制文件，则过滤
		// 注意：使用括号确保先计算频率条件，再检查 content-type，避免运算符优先级问题
		if (locationCount[r.Location] >= 10 || sizeCount[r.Size] >= 10) && r.Ctype != "application/octet-stream" {
			continue // 跳过需要删除的记录
		}
		filteredResults = append(filteredResults, r)

		// 只有在没有指定输出文件时才显示扫描结果信息
		if opt.OutPutFile == "" {
			utils.InforF("%s [%s] %s %s [%s]", r.Url, r.Code, r.Ctype, r.Location, r.Size)
		}
	}

	// 如果需要保存过滤后的结果，可以调用 SaveJsonData 函数
	if len(filteredResults) > 0 {
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
		// 当指定了输出文件时，不在控制台显示扫描结果信息
		// 因为这里是在SaveJsonData函数中，说明用户指定了输出文件
	}
	return nil
}
