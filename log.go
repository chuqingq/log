package log

import (
	"errors"
	"log"

	"github.com/chuqingq/mrpc"
	"github.com/sirupsen/logrus"
)

type Level = logrus.Level

type Options struct {
	// Flags        int // 尽量用默认值
	Name string
	// DB           string
	CountLimit   int    // remote log count limit
	Level        Level  // remote level, 不影响fifo level
	RemoteServer string // remote log server
}

type Logger struct {
	options Options
	// fifo    io.WriteCloser
	fifoHook *fifoHook
	// db      io.WriteCloser
	// count   int // remote log server来控制，无需本地控制
	rpc *mrpc.RPC
	*logrus.Logger
}

func New(options Options) (*Logger, error) {
	var logger Logger
	// default
	if options.Level == 0 {
		options.Level = logrus.InfoLevel
	}
	if options.Name == "" {
		return nil, errors.New("log name is invalid")
	}
	// fifo
	// var err error
	// logger.fifo, err = logondemand.New(options.Name + ".fifo")
	// if err != nil {
	// 	return nil, err
	// }
	// count limit TODO
	if options.CountLimit == 0 {
		options.CountLimit = 10000
		// logger.count = options.CountLimit
	}
	// remote logger
	if options.RemoteServer != "" {
		logger.rpc = mrpc.NewRPC()
	}
	logger.options = options
	logger.Logger = &logrus.Logger{
		Out:       &logger,
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.Level(options.Level),
		Hooks:     make(logrus.LevelHooks),
	}
	// fifo hook
	logger.fifoHook = newFifoHook(options.Name)
	logger.Logger.AddHook(logger.fifoHook)
	return &logger, nil
}

func (l *Logger) Close() {
	if l.fifoHook != nil {
		l.fifoHook.Close()
	}
	// if l.db != nil {
	// 	l.db.Close()
	// }
	if l.rpc != nil {
		l.rpc.Close()
	}
}

func (l *Logger) Write(p []byte) (int, error) {
	if l.rpc != nil {
		req := &Args{
			Client: l.options.Name,
			Bytes:  p,
		}
		err := l.rpc.Call(l.options.RemoteServer, "LogServer.Write", req, &Reply{}) // TODO async
		if err != nil {
			log.Printf("rpc.Call err: %v", err)
		}
	}
	// if l.options.CountLimit != 0 {
	// 	if l.count == 0 {
	// 		// 关闭db
	// 		if l.db != nil {
	// 			l.db.Close()
	// 		}
	// 		// backup db
	// 		os.Rename(l.options.DB, l.options.DB+".bak")
	// 		// reopen db
	// 		var err error
	// 		l.db, err = os.OpenFile(l.options.DB, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// 		if err != nil {
	// 			// os.Stderr.Write([]byte("open db file error"))
	// 			return 0, err
	// 		}
	// 		// reset count
	// 		l.count = l.options.CountLimit
	// 	} else {
	// 		l.count -= 1
	// 	}
	// }
	// write
	// if l.fifo != nil {
	// 	l.fifo.Write(p)
	// }
	// if l.db != nil {
	// 	n, err := l.db.Write(p)
	// 	return n, err
	// }
	return 0, nil
}

// func (l *Logger) WithFields(map[string]interface{}) *Logger

type Fields = logrus.Fields

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
