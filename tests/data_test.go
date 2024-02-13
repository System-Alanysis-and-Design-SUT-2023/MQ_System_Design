package tests

import (
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
)

func TestCreateData(t *testing.T) {
	createData(t, nil, "", "", 0)               //Empty Data
	createData(t, nil, "Iran", "Tehran", 10)    //Iran Data
	createData(t, nil, "Iran", "Mashhad", 15)   //Mashhad Data
	createData(t, nil, "France", "Paris", 20)   //France Data
	createData(t, nil, "USA", "Washington", 30) //USA Data
}

func createData(t *testing.T, expectedErr error, key, value string, ind uint64) {
	if data := models.NewData(key, value, ind); value != data.String() {
		t.Errorf(errorMessage, expectedErr, nil)
	}
}
