package domain

import "fmt"

func referenceCodeToWebScrappingEvent(referenceCode string) WebScrappingEvent {
	url := fmt.Sprintf("https://www.lcsc.com/search?q=%s", referenceCode)

	return WebScrappingEvent{
		Url:  url,
		Code: referenceCode,
	}
}

func ListOfReferencesToWebScrappingEvents(listOfReferences []string) []WebScrappingEvent {
	var events []WebScrappingEvent
	for _, referenceCode := range listOfReferences {
		events = append(events, referenceCodeToWebScrappingEvent(referenceCode))
	}
	return events
}
