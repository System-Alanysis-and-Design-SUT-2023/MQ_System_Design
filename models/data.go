package models

type Data struct {
	Key   string
	Value string
	Index uint64
}

func CreateData(key, value string) func() Data {
	index := uint64(0)

	return func() Data {
		index++
		return Data{
			Key:   key,
			Value: value,
			Index: index,
		}
	}
}
