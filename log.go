package log

import (
	"errors"
	"log"

	"github.com/chuqingq/mrpc"
	"github.com/sirupsen/logrus"
)

type Level = logrus.Level

type Options struct {
	Name         string // name of log client and log file
	CountLimit   int    // remote log count limit
	Level        Level  // remote level, 不影响fifo level
	RemoteServer string // remote log server
}

type Logger struct {
	options  Options
	fifoHook *fifoHook
	rpc      *mrpc.RPC
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
	// countlimit 默认0，表示不限制
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
	var err error
	logger.fifoHook, err = newFifoHook(options.Name)
	if err != nil {
		return nil, err
	}
	logger.Logger.AddHook(logger.fifoHook)
	return &logger, nil
}

func (l *Logger) Close() {
	if l.fifoHook != nil {
		l.fifoHook.Close()
	}
	if l.rpc != nil {
		l.rpc.Close()
	}
}

func (l *Logger) Write(p []byte) (int, error) {
	if l.rpc == nil {
		return 0, nil
	}
	req := &WriteArgs{
		Client:     l.options.Name,
		CountLimit: l.options.CountLimit,
		Bytes:      p,
	}
	err := l.rpc.Call(l.options.RemoteServer, "LogServer.Write", req, &Reply{}) // TODO async
	if err != nil {
		log.Printf("rpc.Call err: %v", err)
		return 0, err
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
