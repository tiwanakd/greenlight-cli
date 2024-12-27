package client

import (
	"bytes"
	"encoding/json"
)

type JsonMap struct {
	kv map[string]string
}

func NewJSONMap() *JsonMap {
	return &JsonMap{make(map[string]string)}
}

func (j *JsonMap) Add(key, value string) {
	if _, ok := j.kv[key]; !ok {
		j.kv[key] = value
	}
}

func (j *JsonMap) CreateJSONReader() (*bytes.Reader, error) {
	js, err := json.Marshal(j.kv)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(js), nil
}
