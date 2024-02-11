package models

type Queue struct {
	Map  map[string]Data
	List []Data
}

func (queue *Queue) Push(data Data) error {
	if _, ok := queue.Map[data.Key]; ok {
		return ErrKeyAlreadyExistsInQueue
	}

	queue.List = append(queue.List, data)
	queue.Map[data.Key] = data
	return nil
}

func (queue *Queue) Pull() (Data, error) {
	if len(queue.List) == 0 {
		return Data{}, ErrQueueIsEmpty
	}

	data := queue.List[0]
	queue.List = queue.List[1:]
	delete(queue.Map, data.Key)
	return data, nil
}

func NewQueue() *Queue {
	return &Queue{
		Map:  make(map[string]Data),
		List: make([]Data, 0),
	}
}
