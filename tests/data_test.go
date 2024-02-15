package tests

import (
	"fmt"
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
)

func TestDataToString(t *testing.T) {
	datas := []models.Data{
		models.NewData("Iran", "Tehran", 1),
		models.NewData("France", "Paris", 2),
		models.NewData("USA", "Washington", 3),
		models.NewData("Iran", "Tehran", 2),
		models.NewData("France", "Antibes", 5),
	}

	for _, data := range datas {
		expected := fmt.Sprintf(`["%s","%s"]`, data.Key, data.Value)
		if data.String() != expected {
			t.Errorf("Expected %s but got %s", expected, data)
		}
	}
}
