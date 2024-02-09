package models

import "github.com/gorilla/websocket"

type Repository struct {
	Queue      *Queue
	Subscriber *Subscriber
	CreateData func(key, value string) func() Data
}

func NewRepository() Repository {
	return Repository{
		Queue:      NewQueue(),
		Subscriber: NewSubscriber(),
		CreateData: func(key, value string) func() Data {
			index := uint64(0)
			return func() Data {
				index++
				return NewData(key, value, index)
			}
		},
	}
}

func (r *Repository) Push(key, value string) error {
	data := r.CreateData(key, value)()
	return r.Queue.Push(data)
}

func (r *Repository) Pull() (Data, error) {
	return r.Queue.Pull()
}

func (r *Repository) Subscribe(conn *websocket.Conn, key string) error {
	return r.Subscriber.Subscribe(conn, key)
}

func (r *Repository) Unsubscribe(key string) error {
	return r.Subscriber.Unsubscribe(key)
}
