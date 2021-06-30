package logger

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
)

var isDev = false

type Level int8

var beelog *logs.BeeLogger

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

func CheckError(err error) error {
	if err != nil {
		write(ERROR, err)
		return err
	}
	return nil
}

func InitLogger(log *logs.BeeLogger, isD bool) {
	beelog = log
	isDev = isD
}

func Debug(v ...interface{}) {

	write(DEBUG, v)
}

func Debugf(format string, v ...interface{}) {

	writef(DEBUG, format, v...)
}

func Info(v ...interface{}) {
	write(INFO, v)
}

func Infof(format string, v ...interface{}) {
	writef(INFO, format, v...)
}

func Warn(v ...interface{}) {
	write(WARN, v)
}

func Warnf(format string, v ...interface{}) {
	writef(WARN, format, v...)
}

func Error(v ...interface{}) {
	write(ERROR, v)
}

func Errorf(format string, v ...interface{}) {
	writef(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	write(FATAL, v)
}

func Fatalf(format string, v ...interface{}) {
	writef(FATAL, format, v...)
}

func write(level Level, v ...interface{}) {

	content := ""
	if isDev {
		pc, _, line, ok := runtime.Caller(2)
		if ok {
			content = fmt.Sprintf("[%s:%d] ", runtime.FuncForPC(pc).Name(), line)
		}
	}

	format := content + "%s"

	switch level {
	case DEBUG:
		beelog.Debug(format, v)
	case INFO:
		beelog.Info(format, v)
	case WARN:
		beelog.Warn(format, v)
	case FATAL:
		beelog.Critical(format, v)
	case ERROR:
		beelog.Error(format, v)
	}
}

func writef(level Level, format string, v ...interface{}) {

	content := ""
	if isDev {
		pc, _, line, ok := runtime.Caller(2)
		if ok {
			content = fmt.Sprintf("[%s:%d] ", runtime.FuncForPC(pc).Name(), line)
		}
	}

	format = content + format

	switch level {
	case DEBUG:
		beelog.Debug(format, v...)
	case INFO:
		beelog.Info(format, v...)
	case WARN:
		beelog.Warn(format, v...)
	case FATAL:
		beelog.Critical(format, v...)
	case ERROR:
		beelog.Error(format, v...)
	}
}
