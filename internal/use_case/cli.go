package use_case

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter/kafka"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"
)

type GetCSVFileFromCLIArgsUseCase struct{}

func NewGetCSVFileFromCLIArgsUseCase() *GetCSVFileFromCLIArgsUseCase {
	return &GetCSVFileFromCLIArgsUseCase{}
}

func (useCase *GetCSVFileFromCLIArgsUseCase) Execute() (io.Reader, error) {
	if len(os.Args) == 1 {
		return nil, fmt.Errorf("Argument -h for help")
	}

	filePathPtr := flag.String("file", "", "Path to the CSV file")

	flag.Parse()

	if *filePathPtr == "" {
		if len(os.Args) < 2 {
			return nil, fmt.Errorf("Argument -h for help")
		}
		filePathPtr = &os.Args[1]
	}

	file, err := os.Open(*filePathPtr)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type CLIApi struct{}

func NewCLIAPIUseCase() *CLIApi {
	return &CLIApi{}
}

func (w *CLIApi) Execute() {
	csvFileReader, err := NewGetCSVFileFromCLIArgsUseCase().Execute()
	if err != nil {
		log.Fatal(err)
	}

	listOfReferences, err := NewGetListOfReferencesUseCase(csvFileReader).Execute()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("List of references: %v \n", listOfReferences)

	events := domain.ListOfReferencesToWebScrappingEvents(listOfReferences)

	log.Printf("List of web scrapping events: %v \n", events)

	eventQueue := kafka.NewKafkaQueueBroker()

	for _, event := range events {
		if err := eventQueue.InsertMessage(event); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("%d events inserted in the queue \n", len(events))
}
