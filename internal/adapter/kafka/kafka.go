package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/mathiasscroccaro/event-driven-go-query/internal/adapter"
	"github.com/mathiasscroccaro/event-driven-go-query/internal/domain"

	kafkaLib "github.com/segmentio/kafka-go"
)

type kafkaQueueBroker struct{}

func NewKafkaQueueBroker() adapter.IQueue {
	return &kafkaQueueBroker{}
}

func (broker *kafkaQueueBroker) InsertMessage(message domain.WebScrappingEvent) error {
	topic := domain.GetEnvOrDefault("KAFKA_TOPIC", "scrapping-topic")
	kafkaHost := domain.GetEnvOrDefault("KAFKA_HOSTS", "localhost:9092")

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatal("failed to serialize message:", err)
	}

	partition := 0

	conn, err := kafkaLib.DialLeader(context.Background(), "tcp", kafkaHost, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafkaLib.Message{Value: serializedMessage},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Printf("inserted message at partition %d to topic %s\n", partition, topic)

	return nil
}

func (broker *kafkaQueueBroker) GetChannelOfMessages() (<-chan domain.WebScrappingEvent, error) {
	topic := domain.GetEnvOrDefault("KAFKA_TOPIC", "scrapping-topic")
	kafkaHosts := domain.GetEnvOrDefault("KAFKA_HOSTS", "localhost:9092")

	reader := kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   []string{kafkaHosts},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e3, // 10kB
		GroupID:   "scrapping-group",
	})

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	eventChannel := make(chan domain.WebScrappingEvent)

	go func() {
		for {
			select {

			case <-ctx.Done():
				close(eventChannel)
				if err := reader.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				log.Println("closing kafka reader")
				return

			default:
				message, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Fatalf("failed to read message: %v", err)
					break
				}

				fmt.Printf("message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value))

				var eventMessage domain.WebScrappingEvent
				if err := json.Unmarshal(message.Value, &eventMessage); err != nil {
					log.Fatalf("failed to unmarshal message: %v", err)
					break
				}

				eventChannel <- eventMessage
			}
		}
	}()

	return eventChannel, nil
}

func (broker *kafkaQueueBroker) GetMessages() ([]domain.WebScrappingEvent, error) {
	topic := domain.GetEnvOrDefault("KAFKA_TOPIC", "scrapping-topic")
	kafkaHosts := domain.GetEnvOrDefault("KAFKA_HOSTS", "localhost:9092")

	reader := kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   []string{kafkaHosts},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e3, // 10kB
		GroupID:   "scrapping-group",
	})

	ctx, cancelContext := context.WithTimeout(context.Background(), 10*time.Second)

	var messages []domain.WebScrappingEvent

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if err := reader.Close(); err != nil {
				log.Fatal("failed to close reader:", err)
			}
			log.Println("closing kafka reader")

			cancelContext()

			if len(messages) > 0 {
				break
			} else {
				return nil, ctx.Err()
			}
		}

		fmt.Printf("message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value))

		var eventMessage domain.WebScrappingEvent
		if err := json.Unmarshal(message.Value, &eventMessage); err != nil {
			log.Fatalf("failed to unmarshal message: %v", err)
			break
		}

		messages = append(messages, eventMessage)
	}

	return messages, nil
}
