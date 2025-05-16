package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/YouChenJun/dirx/libs"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = logrus.New()

func InitLog(options *libs.Options) {
	mwr := io.MultiWriter(os.Stdout)
	logDir := libs.LOGDIR

	if options.Logfile == "" {
		// 如果 Logfile 为空，则不写入文件日志，直接使用标准输出
		logger = &logrus.Logger{
			Out:   mwr,
			Level: logrus.InfoLevel,
			Hooks: make(logrus.LevelHooks),
			Formatter: &prefixed.TextFormatter{
				ForceColors:     false,
				ForceFormatting: true,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			},
		}
		return
	}

	// 否则继续处理用户指定的日志路径
	logDir = filepath.Dir(options.Logfile)
	if !FolderExists(logDir) {
		if err := os.MkdirAll(logDir, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "无法创建日志目录: %v\n", logDir)
			os.Exit(1)
		}
	}

	f, err := os.OpenFile(options.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "打开日志文件出错: %v\n", options.Logfile)
		fmt.Fprintf(os.Stderr, "💡您可能想先通过 %v 命令切换到 %v ", color.HiMagentaString("root user"), color.HiCyanString("sudo su"))
	} else {
		mwr = io.MultiWriter(os.Stdout, f)
	}

	logger = &logrus.Logger{
		Out:   mwr,
		Level: logrus.InfoLevel,
		Hooks: make(logrus.LevelHooks),
		Formatter: &prefixed.TextFormatter{
			ForceColors:     false,
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
}

// PrintLine print seperate line
func PrintLine() {
	dash := color.HiWhiteString("-")
	fmt.Println(strings.Repeat(dash, 40))
}

// GoodF print good message
func GoodF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiBlueString("▶▶"))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// PrefixF print good message
func PrefixF(symbol string, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiGreenString(symbol))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// BannerF print info message
func BannerF(format string, data string) {
	banner := fmt.Sprintf("%v%v%v ", color.WhiteString("["), color.BlueString(format), color.WhiteString("]"))
	fmt.Printf("%v%v\n", banner, color.HiGreenString(data))
}

// BlockF print info message
func BlockF(name string, data string) {
	prefix := fmt.Sprintf("%v ", color.HiGreenString("💬 %v", name))
	message := fmt.Sprintf("%v%v", prefix, data)
	logger.Info(message)
}

// TSPrintF print info message
func TSPrintF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v", color.HiBlueString("▶ "))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// BadBlockF print info message
func BadBlockF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiRedString(" [!] "))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// InforF print info message
func InforF(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

// Infor print info message
func Infor(args ...interface{}) {
	logger.Info(args...)
}

// ErrorF print good message
func ErrorF(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

// Error print good message
func Error(args ...interface{}) {
	logger.Error(args...)
}

// WarnF print good message
func WarnF(format string, args ...interface{}) {
	logger.Warning(fmt.Sprintf(format, args...))
}

// Warn print good message
