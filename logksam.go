package logksam

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var pFile *os.File

// LogContrl @Example
//  var logs logksam.LogContrl = &logksam.LogAttribute{FilePrefixName: "log"}
//	logs.Init()                         初始化
//	defer logs.Close()                  关闭
//	logs.Info("我爱你", "\t", "你知道吗")  正常信息
//	logs.Err("我爱你", "\t", "你知道吗")   错误信息
//	logs.Warn("我爱你", "\t", "你知道吗")  警告信息
type LogContrl interface {
	Init() bool    // Init 初始化日志
	Info(a ...any) // Init 正常日志
	Err(a ...any)  // Init 错误日志
	Warn(a ...any) // Init 警告日志
	Close()        // Init 关闭日志 一般用延迟关闭 defer logs.Close()
}

type LogAttribute struct {
	FilePrefixName string
}

func formatTimeDetailed() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 PM")
}

func formatTime() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
}

func (l *LogAttribute) Info(a ...any) {
	//TODO implement me
	var text string = formatTimeDetailed() + "\t" + "[INFO   ]" + "\t"
	for _, v := range a {
		text += fmt.Sprintf("%s", v)
	}
	text += "\n"
	write := bufio.NewWriter(pFile)
	_, _ = write.WriteString(text)
	_ = write.Flush()
}

func (l *LogAttribute) Err(a ...any) {
	//TODO implement me
	var text string = formatTimeDetailed() + "\t" + "[ERROR  ]" + "\t"
	for _, v := range a {
		text += fmt.Sprintf("%s", v)
	}
	text += "\n"
	write := bufio.NewWriter(pFile)
	_, _ = write.WriteString(text)
	_ = write.Flush()
}

func (l *LogAttribute) Warn(a ...any) {
	//TODO implement me
	var text string = formatTimeDetailed() + "\t" + "[WARNING]" + "\t"
	for _, v := range a {
		text += fmt.Sprintf("%s", v)
	}
	text += "\n"
	write := bufio.NewWriter(pFile)
	_, _ = write.WriteString(text)
	_ = write.Flush()
}

// Init 初始化日志
func (l *LogAttribute) Init() bool {
	dir, _ := filepath.Split(os.Args[0])
	if strings.Contains(dir, "\\Local\\Temp\\") {
		dir, _ = os.Getwd()
	} else {
		dir = dir[:len(dir)-1]
	}
	os.Chdir(dir)
	go goLogTick(l.FilePrefixName)
	for {
		if pFile != nil {
			return true
		}
		time.Sleep(200)
	}
}
func (l *LogAttribute) Close() {
	pFile.Close()
}

func goLogTick(fileName string) {
	dir, _ := filepath.Split(os.Args[0])
	if strings.Contains(dir, "\\Local\\Temp\\") {
		dir, _ = os.Getwd()
	} else {
		dir = dir[:len(dir)-1]
	}

	LogFileName := fmt.Sprintf("%s\\logs\\%s%s.log", dir, fileName, formatTime())

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			panic(err)
		}
	}

	// Create a ticker that runs once per day
	ticker := time.Tick(1 * time.Hour)
	// Set the output of the log package to a file
	logFile, err := os.OpenFile(LogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	pFile = logFile

	for {
		// Wait for the next tick
		<-ticker

		// Create a new log file with a filename that includes the current date
		pFile.Close()
		logFile, err := os.OpenFile(fmt.Sprintf("%s\\logs\\%s%s.log", dir[:len(dir)-1], fileName, formatTime()), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println(err)
		}
		pFile = logFile
	}
}
