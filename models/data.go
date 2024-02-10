package models

type Data struct {
	Key   string
	Value string
	Index uint64
}

func (d Data) String() string {
	return d.Value
}

func NewData(key, value string, index uint64) Data {
	return Data{
		Key:   key,
		Value: value,
		Index: index,
	}
}
