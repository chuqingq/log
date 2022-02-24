package log

import (
	"log"
	"os"
	"sync"

	"github.com/chuqingq/mrpc"
)

type LogServer struct {
	Server string // RPC node name
	// File string // log file
	rpc *mrpc.RPC
	// file *os.File
	// logFiles sync.Map
	clients      map[string]*logClient
	clientsMutex sync.Mutex
}

type logClient struct {
	client string // client name, aka log file name
	file   *os.File
}

// NewLogServer 创建log server（remote模式）
func NewLogServer(server string) (*LogServer, error) {
	s := &LogServer{
		Server:  server,
		clients: make(map[string]*logClient),
	}
	// rpc
	s.rpc = mrpc.NewRPC()
	s.rpc.RegisterService(s.Server, s)
	return s, nil
}

func (s *LogServer) Stop() {
	if s.rpc != nil {
		s.rpc.Close()
	}
	s.clientsMutex.Lock()
	for _, v := range s.clients {
		v.file.Close()
	}
	s.clientsMutex.Unlock()
}

// rpc

// // StartClient 启动log client
// func (s *LogServer) StartClient(client string) error {
// 	logfile := &logFile{
// 		client: client,
// 	}
// 	// open file
// 	var err error
// 	logfile.file, err = os.OpenFile(logfile.client, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	// rpc
// 	s.rpc = mrpc.NewRPC()
// 	s.rpc.RegisterService(s.Name, s)
// 	return nil
// }

type Args struct {
	Client string
	Bytes  []byte
}

type Reply struct {
}

func (s *LogServer) Write(req *Args, reply *Reply) error {
	log.Printf("logserver.Write %v", req.Client)
	// quick path
	if c, ok := s.clients[req.Client]; ok {
		c.file.Write(req.Bytes)
		return nil
	}
	// new log client
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	log.Printf("logserver.Write %v client not exists", req.Client)
	if c, ok := s.clients[req.Client]; ok {
		c.file.Write(req.Bytes)
		return nil
	} else {
		c := &logClient{client: req.Client}
		var err error
		c.file, err = os.OpenFile("/mnt/d/temp/projects/log/"+c.client, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			log.Printf("open file %v error: %v", c.client, err)
			return err
		}
		c.file.Write(req.Bytes)
		log.Printf("logserver.Write %v client created and logging", req.Client)
	}
	return nil
}
