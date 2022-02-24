package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chuqingq/logondemand"
	"github.com/sirupsen/logrus"
)

type fifoHook struct {
	fifo io.WriteCloser
}

func newFifoHook(name string) (*fifoHook, error) {
	fifo, err := logondemand.New(name + ".fifo")
	if err != nil {
		log.Printf("logondemand.New(%v.fifo) error: %v", name, err)
		return nil, err
	}
	return &fifoHook{
		fifo: fifo,
	}, nil
}

func (f *fifoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *fifoHook) Fire(e *logrus.Entry) error {
	b, err := e.Bytes()
	if err != nil {
		return err
	}
	f.fifo.Write(b) // 会报epipe，因此不把错误传回
	fmt.Fprintf(os.Stderr, "fifo write: %v", string(b))
	return nil
}

func (f *fifoHook) Close() error {
	return f.fifo.Close() // 关闭，会让读者退出
}
