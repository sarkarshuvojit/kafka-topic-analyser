package main

import (
	"flag"
	"log/slog"
	"strings"
	"time"

	"com.github/sarkarshuvojit/kafka-topic-scout/pkg/kafka"
)

var (
	brokers        string
	topicRegex     string
	timeoutSeconds int
)

func initFlags() {
	flag.StringVar(&brokers, "brokers", "localhost:290902", "comma seperated brokers in the format host:port; eg: localhost:9092")
	flag.StringVar(&topicRegex, "topics", "*", "regex to match which topics should be read from")
	flag.IntVar(&timeoutSeconds, "timeout", 10, "seconds to wait idle and die if no new message comes")

	flag.Parse()
}

func main() {
	initFlags()

	slog.Debug("Running with args",
		"brokers", brokers,
		"topic regex", topicRegex,
		"timeout", timeoutSeconds,
	)

	brokerList := strings.Split(brokers, ",")
	messages := kafka.ConsumeUntilWaitingFor(brokerList, topicRegex, time.Second*time.Duration(timeoutSeconds))

	slog.Info("Total messages", "count", len(messages))
}
