package log

import (
	"fmt"
	"io"
	"os"

	"github.com/chuqingq/logondemand"
	"github.com/sirupsen/logrus"
)

type Level = logrus.Level

type Options struct {
	Flags      int // 尽量用默认值
	FIFO       string
	DB         string
	CountLimit int
	Level      Level
}

type Logger struct {
	options Options
	fifo    io.WriteCloser
	db      io.WriteCloser
	count   int
	*logrus.Logger
}

func New(options Options) (*Logger, error) {
	var logger Logger
	// default
	if options.Level == 0 {
		options.Level = logrus.InfoLevel
	}
	if options.FIFO != "" {

		var err error
		logger.fifo, err = logondemand.New(options.FIFO)
		fmt.Printf("logondemand error :%v", err)
		if err != nil {
			return nil, err
		}
	}
	if options.CountLimit == 0 {
		options.CountLimit = 10000
		logger.count = options.CountLimit
	}
	logger.options = options
	logger.Logger = &logrus.Logger{
		Out:       &logger,
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.Level(options.Level),
	}
	return &logger, nil
}

func (l *Logger) Write(p []byte) (int, error) {
	if l.options.CountLimit != 0 {
		if l.count == 0 {
			// 关闭db
			if l.db != nil {
				l.db.Close()
			}
			// backup db
			os.Rename(l.options.DB, l.options.DB+".bak")
			// reopen db
			var err error
			l.db, err = os.OpenFile(l.options.DB, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				// os.Stderr.Write([]byte("open db file error"))
				return 0, err
			}
			// reset count
			l.count = l.options.CountLimit
		} else {
			l.count -= 1
		}
	}
	// write
	if l.fifo != nil {
		l.fifo.Write(p)
	}
	n, err := l.db.Write(p)
	return n, err
}

// func (l *Logger) WithFields(map[string]interface{}) *Logger {
// 	// TODO 是一次性的，还是持久的？
// 	return l
// }

type Fields = logrus.Fields
type Field = map[string]interface{}

const (
	LevelDebug = logrus.DebugLevel
	LevelInfo  = logrus.InfoLevel
	LevelWarn  = logrus.WarnLevel
	LevelError = logrus.ErrorLevel
	LevelFatal = logrus.FatalLevel
)

// Debugf
// Infof
// Warnf
// Errorf
// Fatalf
