package domain

type WebScrappingEvent struct {
	Url  string
	Code string
}

type WebScrappingResult struct {
	Url        string
	StatusCode int
	Body       []byte
}
