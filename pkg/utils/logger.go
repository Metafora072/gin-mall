package utils

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

// LogrusObj logrus 是一个流行的 Go 日志库，用于记录应用程序的日志信息。
var LogrusObj *logrus.Logger

// InitLog 初始化日志对象 LogrusObj，并设置日志文件输出及格式.
func InitLog() {
	if LogrusObj != nil {
		src, _ := setOutputFile()

		LogrusObj.Out = src
		return
	}

	logger := logrus.New()
	src, _ := setOutputFile()
	logger.Out = src
	// 设置日志的级别为 logrus.DebugLevel，表示记录所有级别（Debug、Info、Warn、Error、Fatal、Panic）的日志。
	logger.SetLevel(logrus.DebugLevel)
	// 使用 TextFormatter 格式化日志输出，并设置时间戳的格式为 "2006-01-02 15:04:05".
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}

// setOutputFile 设置日志文件的路径和文件名，确保日志文件存在，如果文件不存在则创建文件和目录。
func setOutputFile() (*os.File, error) {
	// 获取当前时间
	now := time.Now()
	logFilePath := ""

	// 通过 os.Getwd() 获取当前工作目录，并将日志文件存放在该目录下的 logs 子目录中。
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs"
	}
	//fmt.Printf("logFilePath = %s\n", logFilePath)
	// 使用 os.Stat() 检查目录是否存在。如果目录不存在，则使用 os.MkdirAll() 创建该目录。目录权限设置为 0777，表示所有用户都可以读写执行。
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) { // 目录不存在
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	// 根据当前日期生成文件名（例如：2025-02-24.log），并将其与日志目录路径结合成完整的文件路径。
	// now.Format("2006-01-02") 这段代码是 Go 语言中的日期时间格式化函数，它将当前的时间 now 格式化为一个特定的字符串，格式为 YYYY-MM-DD。
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)

	// 如果日志文件不存在，使用 os.Create() 创建新文件。
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	// 使用 os.OpenFile() 以追加模式打开文件（os.O_APPEND|os.O_WRONLY），如果文件不存在，则会创建文件。
	// 以 os.O_APPEND 标志打开文件时，文件的写入操作会自动附加到文件的末尾，而不是覆盖现有的内容
	// os.O_WRONLY 表示以 写模式 打开文件。这意味着你只能对文件进行写操作，不能进行读取操作。
	// 文件模式为 os.ModeAppend，表示追加写入而不是覆盖文件内容。
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return src, nil
}
