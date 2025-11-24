package httpx

import "strings"

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
