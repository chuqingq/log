package log

import (
	"encoding/json"
	"io"
	"os"
)

// TODO
// https://github.com/multiprocessio/dsq

func Query(filename string, filter Fields) ([]Fields, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	results := []Fields{}
	for {
		r := Fields{}
		err = dec.Decode(&r)
		if err == io.EOF {
			return results, nil
		}
		if err != nil {
			return nil, err
		}
		if compare(filter, r) {
			results = append(results, r)
		}
	}
}

func compare(filter Fields, record Fields) bool {
	// check all field
	for k, v := range filter {
		if record[k] != v {
			return false
		}
	}
	return true
}
