package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	DEBUG      = "DEBUG"
	INFO       = "INFO"
	WARN       = "WARN"
	ERROR      = "ERROR"
	TIMEFORMAT = "2006-01-02 15:04:05"
)

var cronLog *CronLog

type CronLog struct {
	file   *os.File
	prefix string
	mu     sync.Mutex
}

func (c *CronLog) SetPrefix(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prefix = s
}

func (c *CronLog) Output(s string) {
	var file string
	var line int
	c.mu.Lock()
	defer c.mu.Unlock()
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	buf := fmt.Sprintf(`[%s] [%s]: %s:%d %s`, time.Now().Format(TIMEFORMAT),
		c.prefix, file, line, s,
	)
	os.Stdout.WriteString(buf)
	c.file.WriteString(buf)
}

func InitLogger() *os.File {
	file, err := os.OpenFile("logs/cronjob.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	cronLog = &CronLog{file: file}
	return file
}

func Info(v ...interface{}) {
	cronLog.SetPrefix(INFO)
	cronLog.Output(fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	cronLog.SetPrefix(INFO)
	cronLog.Output(fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	cronLog.SetPrefix(WARN)
	cronLog.Output(fmt.Sprintln(v...))
}

func Warnf(format string, v ...interface{}) {
	cronLog.SetPrefix(WARN)
	cronLog.Output(fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	cronLog.SetPrefix(DEBUG)
	cronLog.Output(fmt.Sprintln(v...))
}

func Debugf(format string, v ...interface{}) {
	cronLog.SetPrefix(DEBUG)
	cronLog.Output(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	cronLog.SetPrefix(ERROR)
	cronLog.Output(fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	cronLog.SetPrefix(ERROR)
	cronLog.Output(fmt.Sprintf(format, v...))
}
