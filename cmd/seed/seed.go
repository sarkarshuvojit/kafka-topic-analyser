package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/hashicorp/go-uuid"
)

var brokerStr string

func main() {
	flag.StringVar(&brokerStr, "brokers", "localhost:9092", "comma seperated broker config; host:port")
	flag.Parse()
	topics := []string{
		"auth.login.req", "auth.loging.res",
		"pay.verifyPayment.req", "pay.verifyPayment.res",
	}

	brokers := strings.Split(brokerStr, ",")
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		slog.Error("Failed to initiate producer", "err", err)
		os.Exit(1)
	}
	defer producer.Close()

	for i := 1; i <= 10; i++ {
		for _, topic := range topics {
			key, err := uuid.GenerateUUID()
			if err != nil {
				slog.Warn("Failed to GenerateUUID", "err", err)
				key = "could-not-generate"
				continue
			}
			producer.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(key),
				Value: sarama.StringEncoder("{}"),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("someHeader"),
						Value: []byte("somevalue"),
					},
				},
			})
		}
	}
}
