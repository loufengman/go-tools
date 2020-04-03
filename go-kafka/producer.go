package kafka

import (
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

// var (
// 	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
// )

func main() {
	sarama.Logger = &log.Logger{}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	msg := &sarama.ProducerMessage{}
	msg.Topic = "hello"
	msg.Partition = int32(-1)
	msg.Key = sarama.StringEncoder("key")
	msg.Value = sarama.ByteEncoder("你好, 世界!")

	producer, err := sarama.NewSyncProducer(strings.Split("localhost:9092", ","), config)
	if err != nil {
		sarama.Logger.Println("Failed to produce message: %s", err)
		os.Exit(500)
	}
	defer producer.Close()

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		sarama.Logger.Println("Failed to produce message: ", err)
	}
	sarama.Logger.Printf("partition=%d, offset=%d\n", partition, offset)
}
