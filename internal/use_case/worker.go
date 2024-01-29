package use_case

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter/kafka"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
)

type WorkerUseCase struct {
}

func NewWorkerUseCase() *WorkerUseCase {
	return &WorkerUseCase{}
}

func (useCase *WorkerUseCase) Execute() {
	kafkaAdapter := kafka.NewKafkaQueueBroker()

	getEventUseCase := NewGetEventsFromQueueUseCase(kafkaAdapter)

	messages, err := getEventUseCase.Execute()
	if err != nil {
		log.Fatal(err)
	}

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
	datasheetReader, err := NewGetDatasheetByUrlUseCase(http.DefaultClient).Execute(message.Url)
	if err != nil {
		fmt.Println("Error downloading", message.Url, err)
		wg.Done()
		return
	}

	outputDir := "./output"
	createDirectoryIfNotExists(outputDir)

	path := filepath.Join(outputDir, fmt.Sprintf("%s.pdf", message.Code))
	NewWriteToFileUseCase(datasheetReader).Execute(path)

	wg.Done()
}

func createDirectoryIfNotExists(directoryPath string) error {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(directoryPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating directory: %v", err)
		}
		fmt.Printf("Directory created: %s\n", directoryPath)
	} else if err != nil {
		return fmt.Errorf("Error checking directory: %v", err)
	} else {
		fmt.Printf("Directory already exists: %s\n", directoryPath)
	}

	return nil
}
