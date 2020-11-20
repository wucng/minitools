package logging

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Lever ...
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

// 默认每隔半小时 创建一个新的文件
func getNums() (int, time.Time) {
	//每隔一定时间创建一个新的log文件
	now := time.Now()
	Nums := now.Minute() / 30 // 每隔半小时
	return Nums, now
}

func (l *Logger) getOutput() error {
	if l.FilePath != "" {
		Nums, now := getNums()
		if Nums != l.Nums { // 要切换
			// 先关闭原文件
			l.Writer.Close()
			// 将原来的文件重命名
			FilePath := fmt.Sprintf("%s-%s-%d.txt", l.FilePath, now.Format("2006-01-02 15:04"), Nums)
			os.Rename(l.FilePath, FilePath)
			// 再重新打开
			file, err := os.OpenFile(l.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
			if err != nil {
				return err
			}
			l.Writer = file
			l.Nums = Nums
		}
	}
	return nil
}

// Logger ...
type Logger struct {
	Class    string
	FilePath string   // 输出到文件的路径
	Nums     int      // 记录是否是要切换文件
	Writer   *os.File //io.Writer
}

// NewLogger 创建logger
func NewLogger(ss ...string) (*Logger, error) {
	var l Logger
	var err error
	var class, FilePath string
	if len(ss) == 0 {
		fmt.Println("至少传入一个参数")
		os.Exit(-1)
	} else if len(ss) < 2 {
		class = ss[0]
		FilePath = ""
	} else {
		class = ss[0]
		FilePath = ss[1]
	}

	l.Class = class
	l.FilePath = FilePath

	if FilePath != "" {
		var file *os.File
		file, err = os.OpenFile(FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if err != nil {
			return &l, err
		}
		l.Writer = file
		Nums, _ := getNums()
		l.Nums = Nums
	} else {
		l.Writer = os.Stdout
	}

	return &l, err
}

// Debug ...
func (l *Logger) Debug(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= DEBUG {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}

// Trace ...
func (l *Logger) Trace(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= TRACE {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}

// Info ...
func (l *Logger) Info(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= INFO {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}

// Warning ...
func (l *Logger) Warning(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= WARNING {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}

// Error ...
func (l *Logger) Error(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= ERROR {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}

// Fatal ...
func (l *Logger) Fatal(s string, v ...interface{}) {
	Class, err := getLever(l.Class)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(1)
	tmp := strings.Split(file, "/")

	if Class <= FATAL {
		str := fmt.Sprintf("[%s] [%s] [%d] [Debug] %s\n", now, tmp[len(tmp)-1], line, s)
		err = l.getOutput()
		if err != nil {
			return
		}
		fmt.Fprintf(l.Writer, str, v...)
	}
}
