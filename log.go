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
		req := &WriteArgs{
			Client:     l.options.Name,
			CountLimit: l.options.CountLimit,
			Bytes:      p,
		}
		err := l.rpc.Call(l.options.RemoteServer, "LogServer.Write", req, &Reply{}) // TODO async
		if err != nil {
			log.Printf("rpc.Call err: %v", err)
		}
	}
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
