package log

import (
	"path"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
)

var (
	_log *logrus.Logger
	once sync.Once
)

const (
	ROTATION_TIME = 86400   // seconds, roate once a day
	MAX_AGE       = 2592000 // seconds, keep a month of log files
)

func Init(base string) {
	once.Do(
		func() {
			_log = logrus.New()
			//_log.Formatter = new(logrus.JSONFormatter)

			debugf, _ := rotatelogs.New(
				path.Join(base, "debug.log-%Y%m%d%H%M%S"),
				rotatelogs.WithLinkName(path.Join(base, "debug.log")),
				rotatelogs.WithMaxAge(time.Duration(MAX_AGE)*time.Second),
				rotatelogs.WithRotationTime(time.Duration(ROTATION_TIME)*time.Second),
			)

			infof, _ := rotatelogs.New(
				path.Join(base, "info.log-%Y%m%d%H%M%S"),
				rotatelogs.WithLinkName(path.Join(base, "info.log")),
				rotatelogs.WithMaxAge(time.Duration(MAX_AGE)*time.Second),
				rotatelogs.WithRotationTime(time.Duration(ROTATION_TIME)*time.Second),
			)

			warnf, _ := rotatelogs.New(
				path.Join(base, "warn.log-%Y%m%d%H%M%S"),
				rotatelogs.WithLinkName(path.Join(base, "warn.log")),
				rotatelogs.WithMaxAge(time.Duration(MAX_AGE)*time.Second),
				rotatelogs.WithRotationTime(time.Duration(ROTATION_TIME)*time.Second),
			)

			errorf, _ := rotatelogs.New(
				path.Join(base, "error.log-%Y%m%d%H%M%S"),
				rotatelogs.WithLinkName(path.Join(base, "error.log")),
				rotatelogs.WithMaxAge(time.Duration(MAX_AGE)*time.Second),
				rotatelogs.WithRotationTime(time.Duration(ROTATION_TIME)*time.Second),
			)

			_log.Hooks.Add(lfshook.NewHook(lfshook.WriterMap{
				logrus.DebugLevel: debugf,
				logrus.InfoLevel:  infof,
				logrus.WarnLevel:  warnf,
				logrus.ErrorLevel: errorf,
			}))
		})
}

func GetLogger() *logrus.Logger {
	if _log == nil {
		_log = logrus.New()
	}
	return _log
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	GetLogger().Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	GetLogger().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Print(args ...interface{}) {
	GetLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warning(args ...interface{}) {
	GetLogger().Warning(args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

func Debugln(args ...interface{}) {
	GetLogger().Debugln(args...)
}

func Infoln(args ...interface{}) {
	GetLogger().Infoln(args...)
}

func Println(args ...interface{}) {
	GetLogger().Println(args...)
}

func Warnln(args ...interface{}) {
	GetLogger().Warnln(args...)
}

func Warningln(args ...interface{}) {
	GetLogger().Warningln(args...)
}

func Errorln(args ...interface{}) {
	GetLogger().Errorln(args...)
}

func Fatalln(args ...interface{}) {
	GetLogger().Fatalln(args...)
}

func Panicln(args ...interface{}) {
	GetLogger().Panicln(args...)
}
