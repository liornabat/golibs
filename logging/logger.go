package logging

import (
	"github.com/sirupsen/logrus"

	//"fmt"
	"os"
)

type Logger struct {
	log *logrus.Entry
	src string
	f   func(msg string)
}

func InitLoggers() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.SetOutput(os.Stdout)

}

func NewLogger(src string) *Logger {
	l := &Logger{
		log: logrus.WithFields(logrus.Fields{"source": src}),
		src: src,
	}
	return l
}
func (l *Logger) NewLogger(src string) *Logger {
	nl := &Logger{
		log: logrus.WithFields(logrus.Fields{"source": l.src + "/" + src}),
		src: l.src + "/" + src,
	}
	return nl
}
func SetDebug(isDebug bool) {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		return
	}
	logrus.SetLevel(logrus.InfoLevel)
}

func (l *Logger) SetDebug(isDebug bool) {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		return
	}
	logrus.SetLevel(logrus.InfoLevel)
}

func (l *Logger) Info(args ...interface{}) {
	if len(args) == 1 {
		l.log.Info(args[0])
		return
	}
	l.log.Info(args)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	l.log.Infof(format, args)
}
func (l *Logger) Error(err error, args ...string) {
	if len(args) == 0 {
		l.log.Error(err)
		return
	}
	if len(args) == 1 {
		if err != nil {
			l.log.Error(args[0], err)
			return
		}
		l.log.Error(args[0])
		return
	}

	if err != nil {
		l.log.Error(args, err)
		return
	}
	l.log.Error(args)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.log.Errorf(format, args)
}
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args)
}

func (l *Logger) Panic(err error, args ...interface{}) {
	if len(args) == 0 {
		l.log.Panic(err)
		return
	}
	if len(args) == 1 {
		l.log.Panic(args[0], err)
		return
	}

	if err != nil {
		l.log.Panic(args, err)
		return
	}
	l.log.Panic(args)

}
