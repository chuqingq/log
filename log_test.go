package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	options := Options{
		FIFO:       "test_log.fifo",
		DB:         "test_log.db",
		CountLimit: 1000,
	}
	logger, err := New(options)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	//
	logger.WithFields(Fields{"module": "my_module", "version": "my_version"}).Errorf("error: %v", "test error")
}
