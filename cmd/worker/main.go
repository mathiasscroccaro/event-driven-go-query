package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter/kafka"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/use_case"
)

func main() {

	kafkaAdapter := kafka.NewKafkaQueueBroker()

	getEvent := use_case.NewGetEventsFromQueueUseCase(kafkaAdapter)

	messages, err := getEvent.Execute()

	fmt.Println(messages, err)

	// if len(messages) > 0 {
	// 	url := messages[0].Url

	// 	fmt.Println("Downloading", url)
	// 	scrapperUseCase := use_case.NewGetDatasheetByUrlUseCase(http.DefaultClient)
	// 	datasheetReader, err := scrapperUseCase.Execute(url)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	fmt.Println("Writing", url)
	// 	writerUseCase := use_case.NewWriteToFileUseCase(datasheetReader)

	// 	writerUseCase.Execute("datasheet.pdf")
	// }

	var wg sync.WaitGroup

	for idx, message := range messages {
		log.Printf("Go Routine %d: %s", idx, message.Code)
		wg.Add(1)
		go processMessage(message, &wg)
	}

	wg.Wait()
	log.Printf("Finished all the tasks")
}

func processMessage(message domain.WebScrappingEvent, wg *sync.WaitGroup) {
	datasheetReader, err := use_case.NewGetDatasheetByUrlUseCase(http.DefaultClient).Execute(message.Url)
	if err != nil {
		fmt.Println("Error downloading", message.Url, err)
	}
	use_case.NewWriteToFileUseCase(datasheetReader).Execute(fmt.Sprintf("%s.pdf", message.Code))
	wg.Done()
}
