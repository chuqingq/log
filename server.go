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

type logClient struct {
	client     string // client name, aka log file name
	file       *os.File
	countlimit int
	count      int
}

func (c *logClient) Write(p []byte) (int, error) {
	// if reach countlimit, backup and reopen
	if c.countlimit != 0 {
		if c.count == 0 {
			// close file
			if c.file != nil {
				c.file.Close()
			}
			// backup file
			os.Rename(c.client, c.client+".bak")
			// reopen db
			var err error
			c.file, err = os.OpenFile(c.client, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			// reset count
			c.count = c.countlimit
			if err != nil {
				// os.Stderr.Write([]byte("open db file error"))
				return 0, err
			}
		} else {
			c.count -= 1
		}
	}
	// write
	if c.file != nil {
		n, err := c.file.Write(p)
		return n, err
	}
	return 0, nil
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

// type ClientOption struct {
// 	Name       string
// 	CountLimit int
// }

// func (s *LogServer) NewClient(req *ClientOption, reply *Reply) error {
// 	s.clientsMutex.Lock()
// 	defer s.clientsMutex.Unlock()
// 	if _, ok := s.clients[req.Name]; ok {
// 		return errors.New("client %s already exists")
// 	}
// 	s.clients[req.Name] = &logClient{
// 		client:     req.Name,
// 		countlimit: req.CountLimit,
// 		count:      req.CountLimit,
// 	}
// 	return nil
// }

type WriteArgs struct {
	Client     string
	CountLimit int
	Bytes      []byte
}

type Reply struct {
}

func (s *LogServer) Write(req *WriteArgs, reply *Reply) error {
	log.Printf("logserver.Write %v", req.Client)
	// quick path
	if c, ok := s.clients[req.Client]; ok {
		c.Write(req.Bytes)
		return nil
	}
	// new log client
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	log.Printf("logserver.Write %v client not exists", req.Client)
	if c, ok := s.clients[req.Client]; ok {
		c.Write(req.Bytes)
		return nil
	} else {
		c := &logClient{
			client:     req.Client,
			countlimit: req.CountLimit,
			count:      req.CountLimit,
		}
		var err error
		c.file, err = os.OpenFile(c.client, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			log.Printf("open file %v error: %v", c.client, err)
			return err
		}
		s.clients[req.Client] = c
		// write
		c.Write(req.Bytes)
		log.Printf("logserver.Write %v client created and logging", req.Client)
	}
	return nil
}
