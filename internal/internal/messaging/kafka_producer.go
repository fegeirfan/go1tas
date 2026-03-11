package messaging

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic: "ticket-created",
	}
}

func PublishTicketCreated(writer *kafka.Writer, message []byte) error {
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		})
}