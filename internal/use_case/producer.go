package use_case

import (
	"log"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
)

type InsertEventInQueueUseCase struct {
	queue adapter.IQueue
}

func NewInsertEventInQueueUseCase(queue adapter.IQueue) *InsertEventInQueueUseCase {
	return &InsertEventInQueueUseCase{
		queue: queue,
	}
}

func (use_case *InsertEventInQueueUseCase) Execute(event domain.WebScrappingEvent) error {
	if error := use_case.queue.InsertMessage(event); error != nil {
		log.Println("Error inserting message in queue", error)
		return error
	}
	return nil
}
