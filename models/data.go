package models

import (
	"encoding/json"
)

type Data struct {
	Key   string
	Value string
	Index uint64
}

func (d Data) String() string {
	str, _ := json.Marshal([]string{d.Key, d.Value})
	return string(str)
}

func NewData(key, value string, index uint64) Data {
	return Data{
		Key:   key,
		Value: value,
		Index: index,
	}
}
