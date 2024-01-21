package use_case

import (
	"testing"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
)

type MockedQueue struct{}

func (m *MockedQueue) GetMessages() ([]domain.WebScrappingEvent, error) {
	return []domain.WebScrappingEvent{
		{
			Url:  "https://www.example.com/1",
			Code: "C1",
		},
		{
			Url:  "https://www.example.com/2",
			Code: "C2",
		},
	}, nil
}

func (m *MockedQueue) InsertMessage(message domain.WebScrappingEvent) error {
	return nil
}

func TestNewGetEventsFromQueueUseCase(t *testing.T) {
	messages, _ := NewGetEventsFromQueueUseCase(&MockedQueue{}).Execute()

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].Url != "https://www.example.com/1" {
		t.Errorf("Expected https://www.example.com/1, got %s", messages[0].Url)
	}
	if messages[1].Url != "https://www.example.com/2" {
		t.Errorf("Expected https://www.example.com/2, got %s", messages[1].Url)
	}
}
