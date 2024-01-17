package main

import (
	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter/kafka"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/use_case"
)

func main() {

	kafkaAdapter := kafka.NewKafkaQueueBroker()

	intertEvent := use_case.NewInsertEventInQueueUseCase(kafkaAdapter)

	intertEvent.Execute(domain.WebScrappingEvent{
		Url:  "https://www.lcsc.com/search?q=C99003",
		Code: "C99003",
	})
}
