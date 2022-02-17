package log

import (
	"log"
	"os"

	"github.com/chuqingq/mrpc"
)

type Server struct {
	Name string // RPC node name
	File string // log file
	rpc  *mrpc.RPC
	file *os.File
}

// NewServer 创建log server（remote模式）
func (s *Server) Start() error {
	// open file
	var err error
	s.file, err = os.OpenFile(s.File, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	// rpc
	s.rpc = mrpc.NewRPC()
	s.rpc.RegisterService(s.Name, s)
	return nil
}

func (s *Server) Stop() {
	if s.rpc != nil {
		s.rpc.Close()
	}
	if s.file != nil {
		s.file.Close()
	}
}

type Reply struct {
}

func (s *Server) Write(p []byte, reply *Reply) error {
	log.Printf("write: %v", p)
	s.file.Write(p)
	return nil // TODO async
}
