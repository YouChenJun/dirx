package httpx

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/YouChenJun/dirx/utils"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Httpx struct {
	Targets    chan string
	Results    []map[string]string
	Errnum     int
	Checks     []string
	Timeout    int
	Method     string
	MaxRespone int
	Threads    int
	FCodes     []string
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
	if h.checkOnline(url) == false {
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
		Timeout: time.Duration(h.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
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

	if err, ok := err.(net.Error); ok && err.Timeout() {
		h.Errnum += 1
		return result, false
	}

	defer respone.Body.Close()

	body, err := ioutil.ReadAll(respone.Body)

	if err != nil {
		result["body"] = ""
	}
	result["url"] = url
	result["body"] = strings.TrimSpace(string(body))
	result["code"] = strings.ReplaceAll(strconv.Itoa(respone.StatusCode), "206", "200")
	result["location"] = respone.Header.Get("Location")
	result["ctype"] = respone.Header.Get("Content-Type")
	result["server"] = respone.Header.Get("Server")
	result["status"] = respone.Status
	result["size"] = strconv.Itoa(len(body)) // 使用读取后的 body 字节长度
	result["time"] = time.Now().Format("2006-01-02 15:04:05")
	return result, true
}

func (h *Httpx) setHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "common.GetRandUserAgent()")
	req.Header.Add("Range", fmt.Sprintf("bytes=0-%d", h.MaxRespone))
	req.Header.Add("Connection", "close")
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

// filter 过滤器
func (h *Httpx) filter(result map[string]string) bool {
	//判断是否过滤状态码
	if exclude_codes(result["code"], h.FCodes) {
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
