package use_case

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

type GetListOfReferencesUseCase struct {
	fileReader io.Reader
}

func NewGetListOfReferencesUseCase(fileReader io.Reader) *GetListOfReferencesUseCase {
	return &GetListOfReferencesUseCase{
		fileReader: fileReader,
	}
}

func (useCase *GetListOfReferencesUseCase) Execute() ([]string, error) {
	records, err := csv.NewReader(useCase.fileReader).ReadAll()
	if err != nil {
		log.Panicf("not was able to read the CSV file: %v", err)
	}

	if len(records) <= 1 {
		return nil, fmt.Errorf("CSV has no records")
	}

	references := make([]string, len(records[1:]))

	for idx, csvLine := range records[1:] {
		references[idx] = csvLine[0]
	}

	return references, nil
}
