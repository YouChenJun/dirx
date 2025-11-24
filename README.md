<div align="center">

```
      _   _               
     | | (_)              
   __| |  _   _ __  __  __
  / _' | | | | '__| \ \/ /
 |(_|  | | | | |     >  <
 \_,___| |_| |_|    /_/\_\
```

# 🚀 DirX - 智能目录扫描工具

[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-0.0.2-blue.svg)](https://github.com/YouChenJun/dirx/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Stars](https://img.shields.io/github/stars/YouChenJun/dirx?style=social)](https://github.com/YouChenJun/dirx/stargazers)

**一款高性能、智能化的 Web 目录扫描工具，支持并发扫描、流式响应检测、大文件智能处理**

[特性](#-核心特性) •
[安装](#-安装) •
[快速开始](#-快速开始) •
[使用文档](#-使用文档) •
[技术原理](#-技术原理)

</div>

---

## ✨ 核心特性

### 🎯 高性能扫描
- ⚡ **并行任务调度**：支持多目标并发扫描，可自定义并行任务数
- 🔥 **多线程扫描**：单个目标支持多线程并发请求
- 📊 **智能过滤**：自动过滤重复响应，减少误报

### 🛡️ 智能检测
- 🌊 **流式响应识别**：自动检测 SSE、gRPC、视频流等 16+ 种流式类型
- 📦 **大文件处理**：智能限制读取大小，防止内存溢出
- ⏱️ **超时优化**：分层超时机制，避免误判

### 🎨 用户友好
- 📝 **多种输出格式**：支持 JSON Lines 格式输出
- 🎨 **彩色日志**：清晰的命令行输出
- 🔧 **灵活配置**：丰富的命令行参数

---

## 📦 安装

### 方式一：编译安装（推荐）

```bash
# 克隆仓库
git clone https://github.com/YouChenJun/dirx.git
cd dirx

# 编译
go build -o dirx

# 运行
./dirx --help
```

### 方式二：直接使用 Go

```bash
go install github.com/YouChenJun/dirx@latest
```

### 方式三：下载预编译版本

从 [Releases](https://github.com/YouChenJun/dirx/releases) 页面下载对应平台的二进制文件。

---

## 🚀 快速开始

### 基础扫描

```bash
# 扫描单个目标
./dirx scan -u http://example.com -w wordlist.txt -o result.json

# 从文件读取多个目标
./dirx scan -T targets.txt -w wordlist.txt -o result.json
```

### 高级用法

```bash
# 使用 4 个并行任务，每个任务 50 个线程
./dirx scan -T targets.txt -w wordlist.txt -c 4 -t 50 -o result.json

# 自定义超时和过滤状态码
./dirx scan -u http://example.com -w wordlist.txt -s 5 -x "404,403,500" -o result.json

# 使用 POST 方法扫描
./dirx scan -u http://example.com -w wordlist.txt -m POST -o result.json
```

---

## 📖 使用文档

### 命令行参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--target` | `-u` | - | 单个扫描目标 URL |
| `--targets` | `-T` | - | 目标文件路径（每行一个 URL） |
| `--wordlist` | `-w` | - | 字典文件路径 |
| `--threads` | `-t` | `20` | 单个任务的扫描线程数 |
| `--concurrency` | `-c` | `1` | 并行任务数 |
| `--timeout` | `-s` | `2` | 单次请求超时时间（秒） |
| `--method` | `-m` | `GET` | HTTP 请求方法 |
| `--fcode` | `-x` | `400,404,406,416,501,502,503,302` | 需要过滤的状态码 |
| `--outfile` | `-o` | - | 输出文件路径 |
| `--log` | `-l` | - | 日志存储路径 |

### 使用示例

#### 1. 基础目录扫描

```bash
./dirx scan \
  -u http://example.com \
  -w /path/to/wordlist.txt \
  -o scan_result.json
```

#### 2. 批量目标扫描

创建目标文件 `targets.txt`：
```
http://example1.com
http://example2.com
http://example3.com
```

执行扫描：
```bash
./dirx scan -T targets.txt -w wordlist.txt -c 3 -t 30 -o results.json
```

#### 3. 高并发扫描

```bash
./dirx scan \
  -T targets.txt \
  -w wordlist.txt \
  -c 5 \          # 5 个并行任务
  -t 50 \         # 每个任务 50 线程
  -s 3 \          # 3 秒超时
  -o results.json
```

**总并发数 = 并行任务数 × 单任务线程数 = 5 × 50 = 250**

#### 4. 自定义过滤

```bash
# 只过滤 404
./dirx scan -u http://example.com -w wordlist.txt -x "404" -o results.json

# 不过滤任何状态码
./dirx scan -u http://example.com -w wordlist.txt -x "" -o results.json
```

### 输出格式

结果以 JSON Lines 格式保存，每行一条记录：

```json
{"url":"http://example.com/admin","code":"200","location":"","ctype":"text/html","server":"nginx","status":"200 OK","size":"1024","body":"...","time":"2024-01-01 12:00:00"}
{"url":"http://example.com/login","code":"200","location":"","ctype":"text/html","server":"nginx","status":"200 OK","size":"2048","body":"...","time":"2024-01-01 12:00:01"}
```

---

## 🔧 技术原理

### 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        DirX 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────┐      ┌──────────────┐     ┌──────────────┐│
│  │ 命令行解析   │─────▶│  任务调度器   │────▶│  结果处理    ││
│  │  (Cobra)    │      │ (Goroutine)  │     │  (Filter)    ││
│  └─────────────┘      └──────────────┘     └──────────────┘│
│                              │                               │
│                              ▼                               │
│                    ┌──────────────────┐                      │
│                    │  并行任务池       │                      │
│                    │ (Concurrency)    │                      │
│                    └──────────────────┘                      │
│                    ┌──────┬──────┬──────┐                   │
│                    ▼      ▼      ▼      ▼                   │
│              ┌────────────────────────────┐                  │
│              │    HTTP 请求引擎 (Httpx)   │                  │
│              ├────────────────────────────┤                  │
│              │ • 智能流式检测              │                  │
│              │ • 大文件处理                │                  │
│              │ • 多线程扫描                │                  │
│              │ • 超时控制                  │                  │
│              └────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### 核心技术

#### 1. 并发模型

采用**两层并发架构**：

```go
// 第一层：并行任务调度
for i := 0; i < concurrency; i++ {
    go func() {
        for url := range urlChan {
            scanSingleTarget(url)
        }
    }()
}

// 第二层：单任务多线程扫描
for thread := 0; thread < threads; thread++ {
    go func() {
        for target := range targets {
            httpRequest(target)
        }
    }()
}
```

**优势**：
- ✅ 资源利用率最大化
- ✅ 避免对单个目标造成过大压力
- ✅ 灵活控制并发粒度

#### 2. 流式响应检测

使用**正则匹配 + 关键词识别**的双重机制：

```go
// 支持的流式类型（部分）
var streamingPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)text/event-stream`),        // SSE
    regexp.MustCompile(`(?i)application/.*stream`),     // 通用流
    regexp.MustCompile(`(?i)application/x-ndjson`),     // NDJSON
    regexp.MustCompile(`(?i)video/.*`),                 // 视频流
    regexp.MustCompile(`(?i)audio/.*`),                 // 音频流
    regexp.MustCompile(`(?i)application/grpc`),         // gRPC
    // ... 16+ 种模式
}
```

**检测流程**：

```
HTTP Response
     │
     ├─► 提取 Content-Type 和 Transfer-Encoding
     │
     ├─► 正则匹配流式模式
     │      ├─ 匹配 ─► 流式响应
     │      └─ 不匹配 ─► 继续
     │
     ├─► 关键词模糊匹配
     │      ├─ stream, ndjson, grpc 等
     │      └─ 匹配 ─► 流式响应
     │
     └─► 常规响应
```

#### 3. 大文件处理

智能检测并限制读取：

```go
// 1. 检测文件大小
contentLength := response.ContentLength
isLargeFile := contentLength > MaxResponseSize || contentLength == -1

// 2. 限制读取
if isLargeFile {
    limitReader := io.LimitReader(response.Body, MaxResponseSize)
    body, _ := io.ReadAll(limitReader)
    // 添加截断标记
    body = append(body, []byte(fmt.Sprintf(
        "\n[Truncated: %d/%d bytes]", len(body), contentLength
    ))...)
}
```

**好处**：
- 🛡️ 防止内存溢出
- ⚡ 提高扫描速度
- 📊 保留关键信息

#### 4. 超时机制

采用**分层超时**而非整体超时：

```go
Transport: &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   2 * time.Second,  // 连接超时
        KeepAlive: 30 * time.Second,
    }).DialContext,
    ResponseHeaderTimeout: 2 * time.Second,  // 响应头超时
    // 不设置整体 Client.Timeout
}
```

| 场景 | 传统超时 | 分层超时 |
|------|---------|---------|
| 小文件 | ✅ 成功 | ✅ 成功 |
| 大文件下载 | ❌ 超时失败 | ✅ 成功（限制读取） |
| SSE 长连接 | ❌ 超时失败 | ✅ 成功（带超时读取） |
| 慢速接口 | ❌ 超时失败 | ✅ 成功（限制读取） |

#### 5. 智能过滤

自动过滤重复响应：

```go
// 统计相同 Location 和 Size 的出现次数
locationCount := make(map[string]int)
sizeCount := make(map[string]int)

// 过滤规则
if locationCount[r.Location] >= 10 || 
   sizeCount[r.Size] >= 10 {
    // 过滤掉高频重复的响应
    continue
}
```

**场景示例**：
- 大量 302 重定向到同一页面 → 自动过滤
- 大量相同大小的 404 页面 → 自动过滤
- 减少 90% 的无效结果

### 性能指标

| 指标 | 数值 |
|------|------|
| 单任务最大线程数 | 无限制（建议 ≤ 100） |
| 最大并行任务数 | 无限制（建议 ≤ 10） |
| 内存占用 | < 100MB（大规模扫描） |
| 扫描速度 | 1000+ 请求/秒（取决于网络） |
| 支持流式类型 | 16+ 种 |

---

## 🎯 应用场景

### 1. 安全测试
- ✅ Web 应用漏洞挖掘
- ✅ 隐藏目录发现
- ✅ 敏感文件检测

### 2. 资产发现
- ✅ 批量目标扫描
- ✅ 子域名探测
- ✅ API 端点枚举

### 3. 渗透测试
- ✅ 信息收集阶段
- ✅ 路径爆破
- ✅ 后台管理页面发现

---

## 📝 字典推荐

推荐使用以下字典项目：

- [SecLists](https://github.com/danielmiessler/SecLists) - 全面的安全测试字典
- [fuzzdb](https://github.com/fuzzdb-project/fuzzdb) - 综合模糊测试字典
- [dirbuster](https://github.com/digination/dirbuster-ng) - 经典目录扫描字典

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 开源协议

本项目采用 MIT 协议 - 详见 [LICENSE](LICENSE) 文件

---

## 👨‍💻 作者

**@Chen_Dark**

- GitHub: [@YouChenJun](https://github.com/YouChenJun)

---

## ⭐ Star History

如果这个项目对你有帮助，请给一个 Star ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=YouChenJun/dirx&type=Date)](https://star-history.com/#YouChenJun/dirx&Date)

---

## 📧 联系方式

如有问题或建议，欢迎：
- 提交 [Issue](https://github.com/YouChenJun/dirx/issues)
- 发起 [Discussion](https://github.com/YouChenJun/dirx/discussions)

---

<div align="center">

**如果这个项目帮到了你，请点个 Star ⭐ 支持一下！**

Made with ❤️ by @Chen_Dark

</div>
