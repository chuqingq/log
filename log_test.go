package log

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	dbname := "test_log.db"
	defer os.Remove(dbname)
	// start
	options := Options{
		FIFO:       "test_log.fifo",
		DB:         dbname,
		CountLimit: 1000,
		Level:      LevelError,
	}
	logger, err := New(options)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	// level
	logger.WithFields(Fields{"module": "my_module", "version": "my_version1"}).Errorf("error: %v", "1")
	logger.WithFields(Fields{"module": "my_module1", "version": "my_version"}).Warnf("error: %v", "2")
	logger.WithFields(Fields{"module": "my_module1", "version": "my_version"}).Errorf("error: %v", "3")
	logger.WithFields(Fields{"module": "my_module", "version": "my_version1"}).Warnf("error: %v", "4")
	// query all
	rs, err := Query(dbname, Fields{})
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("query len not match: %v", len(rs))
	}
	// query filter
	rs, err = Query(dbname, Fields{"module": "my_module", "version": "my_version1"})
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if len(rs) != 1 {
		t.Fatalf("query len not match: %v", len(rs))
	}
}

func TestRemote(t *testing.T) {
	remoteLogServer := "remote-logger-1"
	dbname := "remote-log.db"
	defer os.Remove(dbname)
	// start remote server
	server := &Server{
		Name: remoteLogServer,
		File: dbname,
	}
	server.Start()
	defer server.Stop()
	// client
	options := Options{
		RemoteServer: remoteLogServer,
		Level:        LevelError,
	}
	logger, err := New(options)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	// level
	logger.WithFields(Fields{"module": "my_module", "version": "my_version1"}).Errorf("error: %v", "1")
	logger.WithFields(Fields{"module": "my_module1", "version": "my_version"}).Warnf("error: %v", "2")
	logger.WithFields(Fields{"module": "my_module1", "version": "my_version"}).Errorf("error: %v", "3")
	logger.WithFields(Fields{"module": "my_module", "version": "my_version1"}).Warnf("error: %v", "4")
	// check result
	// query all
	rs, err := Query(dbname, Fields{})
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("query len not match: %v", len(rs))
	}
	// query filter
	rs, err = Query(dbname, Fields{"module": "my_module", "version": "my_version1"})
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if len(rs) != 1 {
		t.Fatalf("query len not match: %v", len(rs))
	}
}
