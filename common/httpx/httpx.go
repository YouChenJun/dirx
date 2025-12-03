package httpx

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/YouChenJun/dirx/utils"
)

type Httpx struct {
	Targets     chan string
	Results     []map[string]string
	Errnum      int
	Checks      []string
	Timeout     int
	Method      string
	MaxRespone  int
	Threads     int
	FCodes      []string
	FBody       []string
	MinBodySize int
}

// 流式响应的 Content-Type 正则匹配模式
var streamingContentTypePatterns = []*regexp.Regexp{
	// Server-Sent Events
	regexp.MustCompile(`(?i)text/event-stream`),

	// 通用流式
	regexp.MustCompile(`(?i)application/.*stream`),
	regexp.MustCompile(`(?i)text/.*stream`),

	// JSON 流式
	regexp.MustCompile(`(?i)application/.*json.*stream`),
	regexp.MustCompile(`(?i)application/x-ndjson`),
	regexp.MustCompile(`(?i)application/jsonlines`),
	regexp.MustCompile(`(?i)application/x-json-stream`),

	// 视频流
	regexp.MustCompile(`(?i)video/.*`),
	regexp.MustCompile(`(?i)application/x-mpegURL`),
	regexp.MustCompile(`(?i)application/vnd\.apple\.mpegurl`),
	regexp.MustCompile(`(?i)application/dash\+xml`),

	// 音频流
	regexp.MustCompile(`(?i)audio/.*`),

	// gRPC 和其他 RPC 流
	regexp.MustCompile(`(?i)application/grpc`),
	regexp.MustCompile(`(?i)application/grpc\+.*`),

	// Chunked transfer (通过 Transfer-Encoding)
	regexp.MustCompile(`(?i)multipart/x-mixed-replace`),

	// WebSocket 升级（虽然通常是 101 状态码）
	regexp.MustCompile(`(?i)application/websocket`),

	// 其他流式协议
	regexp.MustCompile(`(?i)application/octet-stream.*stream`),
}

// Reset 初始化httpx
func (h *Httpx) Reset() *Httpx {
	h.Results, h.Errnum = []map[string]string{}, 0
	h.Checks = []string{}
	return h
}

// Runner 扫描运行runner
func (h *Httpx) Runner(url string, targets []string) []map[string]string {
	//扫描前判断站点是否存活
	if !h.checkOnline(url) {
		utils.WarnF("%s 无法访问 pass...", url)
		return h.Results
	}
	var wg sync.WaitGroup
	for thread := 0; thread < h.Threads; thread++ {
		wg.Add(1)
		go h.threader(&wg)
	}
	h.send_targets(targets).close_targets()
	wg.Wait()
	return h.Results
}

// checkOnline 判断站点是否存活
func (h *Httpx) checkOnline(url string) bool {
	urls := utils.ConcatURLAndWord(url, "index")
	_, flag := h.requester(urls)
	return flag
}

func (h *Httpx) requester(url string) (map[string]string, bool) {
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	var result = make(map[string]string)

	client := &http.Client{
		// 对于大文件下载，需要更长的超时时间
		// 这里设置为连接超时，而不是整体超时
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			// 设置连接和响应头超时，而不是整体超时
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(h.Timeout) * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: time.Duration(h.Timeout) * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest(h.Method, url, nil)
	if err != nil {
		return result, false
	}
	h.setHeaders(request)
	respone, err := client.Do(request)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			h.Errnum += 1
			return result, false
		}
		return result, false
	}

	defer respone.Body.Close()

	// 检测是否为流式响应类型
	contentType := respone.Header.Get("Content-Type")
	transferEncoding := respone.Header.Get("Transfer-Encoding")

	// 多维度判断是否为流式连接
	isLongConnection := h.isStreamingResponse(contentType, transferEncoding)

	// 获取内容大小
	contentLength := respone.ContentLength
	isLargeFile := contentLength > int64(h.MaxRespone) || contentLength == -1

	var body []byte
	if isLongConnection {
		// 对于长连接，只读取部分数据或设置读取超时
		// 使用带超时的读取，避免一直阻塞
		bodyReader := io.LimitReader(respone.Body, int64(h.MaxRespone))
		readChan := make(chan []byte, 1)
		errChan := make(chan error, 1)

		go func() {
			data, readErr := io.ReadAll(bodyReader)
			if readErr != nil {
				errChan <- readErr
			} else {
				readChan <- data
			}
		}()

		select {
		case body = <-readChan:
			// 读取成功
		case <-errChan:
			body = []byte{}
		case <-time.After(time.Duration(h.Timeout) * time.Second):
			// 超时，但这是正常的长连接行为
			body = []byte("[Long Connection Detected]")
		}
	} else if isLargeFile {
		// 大文件或未知大小：只读取前 MaxRespone 字节
		limitReader := io.LimitReader(respone.Body, int64(h.MaxRespone))
		body, err = io.ReadAll(limitReader)
		if err != nil {
			body = []byte{}
		}
		// 标记这是部分内容
		if contentLength > int64(h.MaxRespone) {
			body = append(body, []byte(fmt.Sprintf("\n[Truncated: %d/%d bytes]", len(body), contentLength))...)
		} else if contentLength == -1 {
			body = append(body, []byte(fmt.Sprintf("\n[Truncated: %d bytes read, total size unknown]", len(body)))...)
		}
	} else {
		// 常规小文件，正常读取
		body, err = io.ReadAll(respone.Body)
		if err != nil {
			body = []byte{}
		}
	}

	result["url"] = url
	result["body"] = strings.TrimSpace(string(body))
	result["code"] = strings.ReplaceAll(strconv.Itoa(respone.StatusCode), "206", "200")
	result["location"] = respone.Header.Get("Location")
	result["ctype"] = contentType
	result["server"] = respone.Header.Get("Server")
	result["status"] = respone.Status
	result["size"] = strconv.Itoa(len(body)) // 使用读取后的 body 字节长度
	// 设置content-length，包括0的情况，以便filter函数可以检查
	if contentLength >= 0 {
		result["content-length"] = strconv.FormatInt(contentLength, 10)
	}
	result["time"] = time.Now().Format("2006-01-02 15:04:05")
	return result, true
}

func (h *Httpx) setHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "common.GetRandUserAgent()")
	req.Header.Add("Range", fmt.Sprintf("bytes=0-%d", h.MaxRespone))
	req.Header.Add("Connection", "close")
}

// isStreamingResponse 判断是否为流式响应
func (h *Httpx) isStreamingResponse(contentType, transferEncoding string) bool {
	// 1. 检查 Transfer-Encoding 是否为 chunked（分块传输）
	if strings.Contains(strings.ToLower(transferEncoding), "chunked") {
		// chunked 本身不一定是流式，但结合 Content-Length 为 -1 时通常是流式
		// 这里暂时不单独判断，留给后续逻辑
	}

	// 2. 使用正则匹配 Content-Type
	if contentType != "" {
		for _, pattern := range streamingContentTypePatterns {
			if pattern.MatchString(contentType) {
				return true
			}
		}
	}

	// 3. 特殊关键词匹配（兜底策略）
	contentTypeLower := strings.ToLower(contentType)
	streamKeywords := []string{
		"stream",
		"ndjson",
		"jsonlines",
		"event-stream",
		"grpc",
		"websocket",
	}

	for _, keyword := range streamKeywords {
		if strings.Contains(contentTypeLower, keyword) {
			return true
		}
	}

	return false
}

// threader 接受发送的扫描任务，同时调度扫描过滤程序
func (h *Httpx) threader(wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range h.Targets {
		result, flag := h.requester(url)
		if flag && h.filter(result) {
			//utils.InforF("%v [%v] %v %v [%v]", result["url"], result["code"], result["ctype"], result["location"], result["size"])
			h.Results = append(h.Results, result)
		}
	}
}

func (h *Httpx) send_targets(targets []string) *Httpx {
	for _, url := range targets {
		h.Targets <- url
	}
	return h
}

func (h *Httpx) close_targets() {
	time.Sleep(time.Duration(h.Timeout) * time.Second)
	close(h.Targets)
}

// filter 过滤器 false则过滤掉 不保留 过滤无效响应（状态码、空内容等）、特定的 body 内容过滤
func (h *Httpx) filter(result map[string]string) bool {
	//判断是否过滤状态码
	if exclude_codes(result["code"], h.FCodes) || exclude_body(result["body"]) {
		return false
	}
	//判断body大小是否为0
	if result["size"] == "0" {
		return false
	}
	//判断content-length是否为0
	if contentLength, exists := result["content-length"]; exists {
		if contentLength == "0" {
			return false
		}
	}
	//判断是否过滤特定body内容（根据body大小和关键字）
	bodySize, _ := strconv.Atoi(result["size"])
	if exclude_body_content(result["body"], h.FBody, bodySize, h.MinBodySize) {
		return false
	}
	return true
}
func extractTitleWithGoquery(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("title").Text())
}
