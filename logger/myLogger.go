// $GOPATH/src/day01/myLoger/myLoger.go
package minitools

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Lever uint8

const (
	UNKONW Lever = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func getLever(lever string) (Lever, error) {
	switch strings.ToLower(lever) {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		return UNKONW, errors.New("指定级别错误")
	}
}

func getOutput(path string) (writer io.Writer) {
	// var writer io.Writer
	//每隔一定时间创建一个新的log文件
	now := time.Now()
	path = fmt.Sprintf("%s-%s-%d.txt", path, now.Format("2006-01-02 15:04"), now.Minute()/30) // 每个30分钟 保存到一个新的文件

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		writer = os.Stdout
	} else {
		writer = file
	}
	return
}

type logger struct {
	lever    string
	filePath string // 输出到文件的路径
}

func NewLogger(ss ...string) logger {
	var lever, path string
	if len(ss) == 0 {
		fmt.Println("至少传入一个参数")
		os.Exit(-1)
	} else if len(ss) < 2 {
		lever = ss[0]
		path = ""
	} else {
		lever = ss[0]
		path = ss[1]
	}
	return logger{
		lever:    lever,
		filePath: path,
	}
}

func NewLogger2(lever string, path string) logger {
	return logger{
		lever:    lever,
		filePath: path,
	}
}

func (l logger) Debug(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= DEBUG {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}

func (l logger) Trace(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= TRACE {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [trace] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}

func (l logger) Info(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= INFO {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [info] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}

func (l logger) Warning(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= WARNING {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [warning] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}

func (l logger) Error(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= ERROR {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [error] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}

func (l logger) Fatal(s string) {
	lever, err := getLever(l.lever)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if lever <= FATAL {
		fmt.Fprintf(getOutput(l.filePath), "[%s] [%s] [%d] [fatal] %s\n", now, tmp[len(tmp)-1], line, s)
	}
}
