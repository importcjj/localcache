package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// level
const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
)

var levelNames = [4]string{"DEBUG", "INFO", "WARN", "ERROR"}

var colors = map[string]int{
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
}

var levelColors = map[int]string{
	DEBUG: "blue",
	INFO:  "green",
	WARN:  "yellow",
	ERROR: "red",
}

var defaultLogger = New(os.Stderr)

// Logger ...
type Logger struct {
	enabled bool
	level   int
	colored bool
	w       io.Writer
	mutex   sync.Mutex
}

// New ...
func New(w io.Writer) *Logger {
	return &Logger{
		enabled: true,
		level:   DEBUG,
		colored: true,
		w:       w,
	}
}

// SetLevel 设置日志而级别
func SetLevel(l int) {
	defaultLogger.level = l % len(levelNames)
}

// SetColored 是否开启颜色
func SetColored(b bool) {
	defaultLogger.colored = b
}

// SetWriter 设置writer
func SetWriter(writer io.Writer) {
	defaultLogger.w = writer
}

// Disable 关闭日志
func Disable() {
	defaultLogger.enabled = false
}

// Enable 开启日志
func Enable() {
	defaultLogger.enabled = true
}

// Colored 为字符串加上颜色
func Colored(color string, text string) string {
	return fmt.Sprintf("\033[3%dm%s\033[0m", colors[color], text)
}

// header 构造日志输出格式的Header
func (logger *Logger) header(time string, level int, filepath string, line int) string {
	levelName := fmt.Sprintf("%s", levelNames[level])
	levelColor := levelColors[level]

	if logger.colored {
		levelName = Colored(levelColor, levelName)
	}

	return fmt.Sprintf("%s [%s] [%s:%d]", time, levelName, filepath, line)
}

func (logger *Logger) println(l int, msg string) error {
	if logger.enabled && l >= logger.level {
		_, filename, line, _ := runtime.Caller(2)
		pkgName := path.Base(path.Dir(filename))
		filepath := path.Join(pkgName, path.Base(filename))
		now := time.Now().Format("2006/02/02 15:04:05")

		header := logger.header(now, l, filepath, line)
		logger.mutex.Lock()
		defer logger.mutex.Unlock()
		_, err := fmt.Fprintf(logger.w, "%s %s\n", header, msg)
		return err
	}
	return nil
}

// Debugf 调试
func (logger *Logger) Debugf(format string, a ...interface{}) error {
	return logger.println(DEBUG, fmt.Sprintf(format, a...))
}

// Infof 普通
func (logger *Logger) Infof(format string, a ...interface{}) error {
	return logger.println(INFO, fmt.Sprintf(format, a...))
}

// Warnf 警告
func (logger *Logger) Warnf(format string, a ...interface{}) error {
	return logger.println(WARN, fmt.Sprintf(format, a...))
}

// Errorf 错误
func (logger *Logger) Errorf(format string, a ...interface{}) error {
	return logger.println(ERROR, fmt.Sprintf(format, a...))
}

// Fatalf 错误并退出
func (logger *Logger) Fatalf(format string, a ...interface{}) {
	logger.println(ERROR, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Debug 调试
func (logger *Logger) Debug(a ...interface{}) error {
	return logger.println(DEBUG, fmt.Sprint(a...))
}

// Info 普通
func (logger *Logger) Info(a ...interface{}) error {
	return logger.println(INFO, fmt.Sprint(a...))
}

// Warn 普通
func (logger *Logger) Warn(a ...interface{}) error {
	return logger.println(WARN, fmt.Sprint(a...))
}

// Error 错误
func (logger *Logger) Error(a ...interface{}) error {
	return logger.println(ERROR, fmt.Sprint(a...))
}

// Fatal 错误并退出
func (logger *Logger) Fatal(a ...interface{}) {
	logger.println(ERROR, fmt.Sprint(a...))
	os.Exit(1)
}
