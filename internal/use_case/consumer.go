package use_case

import (
	"errors"
	"log"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
)

type GetEventsFromQueueUseCase struct {
	queue adapter.IQueue
}

func NewGetEventsFromQueueUseCase(queue adapter.IQueue) *GetEventsFromQueueUseCase {
	return &GetEventsFromQueueUseCase{
		queue: queue,
	}
}

func (use_case *GetEventsFromQueueUseCase) Execute() ([]domain.WebScrappingEvent, error) {
	messages, error := use_case.queue.GetMessages()

	if len(messages) == 0 {
		log.Println("Queue is empty")
		return nil, errors.New("queue is empty")
	}

	if error != nil {
		log.Println("Error getting messages from queue", error)
		return nil, error
	}

	return messages, nil
}
