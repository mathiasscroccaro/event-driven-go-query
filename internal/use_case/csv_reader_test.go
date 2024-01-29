package use_case

import (
	"strings"
	"testing"
)

func TestGetListOfReferencesCase(t *testing.T) {
	csvReader := strings.NewReader(`reference
C01
C02
C03`)

	listOfReference, err := NewGetListOfReferencesUseCase(csvReader).Execute()
	if err != nil {
		t.Error(err)
	}

	if len(listOfReference) != 3 {
		t.Errorf("Expected 3 references, got %d", len(listOfReference))
	}

	if listOfReference[0] != "C01" {
		t.Errorf("Expected C01, got %s", listOfReference[0])
	}
	if listOfReference[1] != "C02" {
		t.Errorf("Expected C02, got %s", listOfReference[1])
	}
	if listOfReference[2] != "C03" {
		t.Errorf("Expected C03, got %s", listOfReference[2])
	}
}
