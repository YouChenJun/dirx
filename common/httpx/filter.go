package httpx

import "slices"

/*过滤器核心代码*/

// exclude_codes 判断状态码是否在排除列表内
func exclude_codes(code string, FCodes []string) bool {
	return slices.Contains(FCodes, code)
}
