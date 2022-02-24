# log
simple log for Go.

### Features

1. structured
2. compatible with std log
3. log on demand: real-time view, not stored
4. logging with extra data: WithFields()
5. level control: Only log the specified severity or above
6. count limit
7. query

### Example

```Go
func TestLog(t *testing.T) {
	dbname := "test_log.db"
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
```

```bash
% cat test_log.db
```

```json
{"level":"error","module":"my_module","msg":"error: 1","time":"2022-02-15T16:44:47+08:00","version":"my_version1"}
{"level":"error","module":"my_module1","msg":"error: 3","time":"2022-02-15T16:44:47+08:00","version":"my_version"}
```

# TODO

- [x] logondemand模式需要确保任何级别均可查看
- [x] 远程模式：区分客户端，为不同客户端保存不同日志文件
- [x] 无需写入到本地文件（暂不开放此能力）
- [ ] 用test hook做测试
- [x] count limit在server端支持
