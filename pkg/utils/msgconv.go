package utils

import "github.com/IBM/sarama"

func SaramaHeaderToMap(headers []*sarama.RecordHeader) map[string]string {
	finalHeaders := make(map[string]string)
	for _, record := range headers {
		finalHeaders[string(record.Key)] = string(record.Value)
	}
	return finalHeaders
}
