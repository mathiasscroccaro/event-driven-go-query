package adapter

import "github.com/mathiasscroccaro/event-driven-go-query/internal/domain"

type IQueue interface {
	InsertMessage(message domain.WebScrappingEvent) error
	GetChannelOfMessages() (<-chan domain.WebScrappingEvent, error)
	GetMessages() ([]domain.WebScrappingEvent, error)
}
