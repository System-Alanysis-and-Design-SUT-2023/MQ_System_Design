package internals

import (
	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
	"github.com/gorilla/websocket"
)

type Repository struct {
	IndexCount uint64
	Queue      *models.Queue
	Subscriber *models.Subscriber
}

func (r *Repository) Push(key, value string) error {
	r.IndexCount++
	data := models.NewData(key, value, r.IndexCount)

	if r.Subscriber.HasSubscriber() {
		return r.Subscriber.SendData(data)
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
		IndexCount: 0,
		Queue:      models.NewQueue(),
		Subscriber: models.NewSubscriber(),
	}
}
