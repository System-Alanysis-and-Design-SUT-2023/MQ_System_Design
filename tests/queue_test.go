package tests

import (
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
)

const errorMessage = `expected "%v", got "%v"`

func TestPushAndPull(t *testing.T) {
	queue := models.NewQueue()

	emptyData := models.NewData("", "", 0)
	iran := models.NewData("Iran", "Tehran", 10)
	mashhad := models.NewData("Iran", "Mashhad", 15)
	france := models.NewData("France", "Paris", 20)
	usa := models.NewData("USA", "Washington", 30)

	pullQueue(t, queue, emptyData, models.ErrQueueIsEmpty)
	pushQueue(t, queue, iran, nil)
	pushQueue(t, queue, mashhad, models.ErrKeyAlreadyExistsInQueue)
	pushQueue(t, queue, france, nil)
	pullQueue(t, queue, iran, nil)
	pushQueue(t, queue, usa, nil)
	pullQueue(t, queue, france, nil)
	pushQueue(t, queue, iran, nil)
	pushQueue(t, queue, usa, models.ErrKeyAlreadyExistsInQueue)
	pullQueue(t, queue, usa, nil)
	pullQueue(t, queue, iran, nil)
	pullQueue(t, queue, emptyData, models.ErrQueueIsEmpty)
}

func pushQueue(t *testing.T, queue *models.Queue, data models.Data, expectedErr error) {
	if err := queue.Push(data); err != expectedErr {
		t.Errorf(errorMessage, expectedErr, err)
	}
}

func pullQueue(t *testing.T, queue *models.Queue, expectedData models.Data, expectedErr error) {
	if data, err := queue.Pull(); err != expectedErr {
		t.Errorf(errorMessage, nil, err)
	} else if data != expectedData {
		t.Errorf(errorMessage, expectedData, data)
	}
}
