package httpx

import (
	"strings"
)

/*过滤器核心代码*/

// exclude_codes 判断状态码是否在排除列表内
func exclude_codes(code string, FCodes []string) bool {
	// 去除状态码的空格，确保匹配
	code = strings.TrimSpace(code)
	for _, fcode := range FCodes {
		if strings.TrimSpace(fcode) == code {
			return true
		}
	}
	return false
}

func exclude_body(body string) bool {
	return body == "Forbidden"
}

// exclude_body_content 判断响应body是否包含需要过滤的关键字
func exclude_body_content(body string, FBody []string, bodySize int, minBodySize int) bool {
	// 只有当body大小小于阈值时才检查内容
	if len(body) < minBodySize {
		// 检查body是否包含任何过滤关键字
		for _, keyword := range FBody {
			if strings.Contains(body, strings.TrimSpace(keyword)) {
				return true
			}
		}
	}
	return false
}
