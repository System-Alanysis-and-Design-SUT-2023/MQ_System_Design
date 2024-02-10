package internals

import (
	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
	"github.com/gorilla/websocket"
)

type Repository struct {
	Queue      *models.Queue
	Subscriber *models.Subscriber
	CreateData func(key, value string) func() models.Data
}

func (r *Repository) Push(key, value string) error {
	data := r.CreateData(key, value)()
	if r.Subscriber.HasSubscriber() {
		r.Subscriber.SendData(data)
		return nil
	}

	return r.Queue.Push(data)
}

func (r *Repository) Pull() (models.Data, error) {
	return r.Queue.Pull()
}

func (r *Repository) Subscribe(conn *websocket.Conn) error {
	return r.Subscriber.Subscribe(conn)
}

func (r *Repository) Unsubscribe(conn *websocket.Conn) error {
	return r.Subscriber.Unsubscribe(conn)
}

func NewRepository() *Repository {
	return &Repository{
		Queue:      models.NewQueue(),
		Subscriber: models.NewSubscriber(),
		CreateData: func(key, value string) func() models.Data {
			index := uint64(0)
			return func() models.Data {
				index++
				return models.NewData(key, value, index)
			}
		},
	}
}
