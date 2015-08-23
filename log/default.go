package log

import (
	"github.com/agtorre/gocolorize"
	"io"
	"log"
	"os"
)

var Logger ILogger

func init() {
	SetLogger(newDefaultLogger())
}

type ILogger interface {
	Fatal(string, Any)
	Error(string, Any)
	Warn(string, Any)
	Info(string, Any)
}

func SetLogger(ilogger ILogger) {
	Logger = ilogger
}

func Fatal(label string, value Any) {
	Logger.Fatal(label, value)
}

func Error(label string, value Any) {
	Logger.Error(label, value)
}

func Warn(label string, value Any) {
	Logger.Warn(label, value)
}

func Info(label string, value Any) {
	Logger.Info(label, value)
}

func newColorLogger(out io.Writer, severity string, color gocolorize.Color) *log.Logger {
	c := gocolorize.Colorize{Fg: color}
	return log.New(out, c.Paint("severity:"+severity)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
}

type DefaultLogger struct {
	loggers map[string]*log.Logger
}

func newDefaultLogger() *DefaultLogger {
	loggers := make(map[string]*log.Logger)
	loggers["fatal"] = newColorLogger(os.Stderr, "FATAL", gocolorize.Blue)
	loggers["error"] = newColorLogger(os.Stderr, "ERROR", gocolorize.Red)
	loggers["warn"] = newColorLogger(os.Stderr, "WARN", gocolorize.Yellow)
	loggers["info"] = newColorLogger(os.Stdout, "INFO", gocolorize.Cyan)
	return &DefaultLogger{loggers}
}

func (l *DefaultLogger) Fatal(label string, value Any) {
	l.PrintLog(l.loggers["fatal"], label, value)
}

func (l *DefaultLogger) Error(label string, value Any) {
	l.PrintLog(l.loggers["error"], label, value)
}

func (l *DefaultLogger) Warn(label string, value Any) {
	l.PrintLog(l.loggers["warn"], label, value)
}

func (l *DefaultLogger) Info(label string, value Any) {
	l.PrintLog(l.loggers["info"], label, value)
}

func (l *DefaultLogger) PrintLog(logger *log.Logger, label string, value Any) {
	logger.Printf("label:%s\tvalue:%v\t", label, value)
}

type Any interface{}
