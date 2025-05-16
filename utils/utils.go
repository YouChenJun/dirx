package utils

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"
)

// FileExists 判断文件是否存在
func FileExists(filename string) bool {
	filename = NormalizePath(filename)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// NormalizePath the path
func NormalizePath(path string) string {
	if strings.HasPrefix(path, "~") {
		path, _ = homedir.Expand(path)
	}
	return path
}

// ReadingFileUnique 读取文件内容 并且去重返回数组
func ReadingFileUnique(filename string) []string {
	var result []string
	if strings.Contains(filename, "~") {
		filename, _ = homedir.Expand(filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer file.Close()

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		// unique stuff
		if val == "" {
			continue
		}
		if seen[val] {
			continue
		}

		seen[val] = true
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}

// RemoveDuplicateStrings 去重字符串切片
func RemoveDuplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// ConcatURLAndWord 拼接 URL 和 Word，去除 url 末尾的 / 和 word 开头的 /
func ConcatURLAndWord(url string, word string) string {
	// 去除 URL 末尾的斜杠
	url = strings.TrimSuffix(url, "/")

	// 去除 word 左侧的斜杠
	word = strings.TrimLeft(word, "/")

	// 拼接结果
	return fmt.Sprintf("%s/%s", url, word)
}

// FolderExists 检查文件是否存在
func FolderExists(foldername string) bool {
	foldername = NormalizePath(foldername)
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		return false
	}
	return true
}
