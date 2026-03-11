package main

import (
	"context"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "ticket-created",
		GroupID: "worker-group",
	})

	fmt.Println("Worker started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		fmt.Println("Received event:", string(msg.Value))

		// contoh: kirim email / logging / notif
	}
}