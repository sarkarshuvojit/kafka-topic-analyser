package kafka

import (
	"log"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"com.github/sarkarshuvojit/kafka-topic-scout/pkg/utils"
	"github.com/IBM/sarama"
)

type AnalyseConfig struct {
	Brokers    []string
	TopicRegex string
}

type UnitComp struct {
	Min float64
	Max float64
	Avg float64
}

type TopicReport struct {
	RequestTopic  string
	ResponseTopic string
	CorrelationId string

	// Units
	ReqTimes   UnitComp
	Throughput string
}

type AnalysisReport struct {
	Topics map[string]TopicReport

	ReqTimes   UnitComp
	Throughput string
}

func Analyse(analyseConfig *AnalyseConfig) (*AnalysisReport, error) {
	//_consumer, _err := createConsumer(analyseConfig.Brokers, analyseConfig.TopicRegex)
	return nil, nil
}

type Message struct {
	Key     string
	Topic   string
	Payload []byte
	Headers map[string]string
}

// ConsumeUntilWaitingFor connects to brokers, consumes topics matching with topicRegex
// Keep consuming until the loop waits for timeoutSeconds without any new message
// Then returns the messages as a list of Message struct
func ConsumeUntilWaitingFor(brokers []string, topicRegex string, timeout time.Duration) []Message {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka client: %v", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("Error fetching Kafka topics: %v", err)
	}

	var matchedTopics []string
	for _, topic := range topics {
		slog.Debug("Checking topic against regex", "topic", topic, "regex", topicRegex)
		matched, _ := regexp.MatchString(topicRegex, topic)
		if matched {
			slog.Info("Matched", "topic", topic, "regex", topicRegex)
			matchedTopics = append(matchedTopics, topic)
		} else {
			slog.Warn("Didn't Match", "topic", topic, "regex", topicRegex)
		}
	}

	var wg sync.WaitGroup
	messageChannel := make(chan *sarama.ConsumerMessage)
	stopChannel := make(chan struct{})

	for _, topic := range matchedTopics {
		slog.Debug("Listening to topic", "topic", topic)
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			log.Fatalf("Error fetching partitions for topic %s: %v", topic, err)
		}

		for _, partition := range partitions {
			slog.Debug("Listening to topic.partition",
				"topic", topic,
				"partition", partition,
			)
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				log.Fatalf("Error starting partition consumer for topic %s partition %d: %v", topic, partition, err)
			}

			wg.Add(1)
			go func(
				pc sarama.PartitionConsumer,
				topic string,
				parition int32,
			) {
				defer wg.Done()
				defer pc.Close()
				for {
					select {
					case msg := <-pc.Messages():
						slog.Debug("Got new message",
							"msg", msg,
							"topic", topic,
							"partition", "parition",
						)
						messageChannel <- msg
					case <-stopChannel:
						return
					}
				}
			}(partitionConsumer, topic, partition)
		}
	}

	var messages []Message
	timeoutTimer := time.NewTimer(timeout)
	messageReceived := false

	go func() {
		for {
			select {
			case msg := <-messageChannel:
				slog.Debug("New message recieved", "key", string(msg.Key))
				messageReceived = true
				timeoutTimer.Reset(timeout)
				messages = append(messages, Message{
					Topic:   msg.Topic,
					Payload: msg.Value,
					Headers: utils.SaramaHeaderToMap(msg.Headers),
				})
			case <-timeoutTimer.C:
				slog.Debug("Timeout happened")
				if !messageReceived {
					close(stopChannel)
					wg.Wait()
					return
				}
				messageReceived = false
				timeoutTimer.Reset(timeout)
			}
		}
	}()

	wg.Wait()
	return messages
}

func createConsumer(brokers []string, topicRegex string) (bool, error) {
	panic("unimplemented")
}
