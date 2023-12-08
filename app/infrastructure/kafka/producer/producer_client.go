package producer

import (
	"encoding/json"
	"golang.org/x/net/context"
	"cnores-skeleton-golang-app/app/application/constant"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"cnores-skeleton-golang-app/app/infrastructure/kafka/settings"
	"cnores-skeleton-golang-app/app/shared/utils/log"
)

type Closeable interface {
	Close()
}

type ProducerInterface[T any] interface {
	Produce(ctx context.Context, key string, value *T) error
}

type Client[T any] struct {
	producer *kafka.Producer
	topic    string
	brokers  string
	ssl      settings_connection.SecureConnection
	settings settings_connection.ProducerConnection
}

type ClientDisabled[T any] struct {
	Client[T]
}

var producerSingleton *kafka.Producer

func NewKafkaClient[T any](address []string, topic string, ssl settings_connection.SecureConnection, settings settings_connection.ProducerConnection) ProducerInterface[T] {
	log.Info("Creating producer instance addr: %s  topic: %s", address, topic)
	return &Client[T]{
		brokers:  strings.Join(address, ","),
		topic:    topic,
		ssl:      ssl,
		settings: settings,
	}
}
func (client *Client[T]) SetKafkaConfig(address []string, topic string, ssl settings_connection.SecureConnection, settings settings_connection.ProducerConnection) {
	client.brokers = strings.Join(address, ",")
	client.topic = topic
	client.ssl = ssl
	client.settings = settings
}

func (client *Client[T]) getKafkaProducer() (*kafka.Producer, error) {
	config := kafka.ConfigMap{}

	if client.ssl.IsEnabled {
		config = kafka.ConfigMap{
			"bootstrap.servers":        client.brokers,
			"security.protocol":        "SSL",
			"ssl.key.location":         client.ssl.KeyLocation,
			"ssl.key.password":         client.ssl.KeyPassword,
			"ssl.certificate.location": client.ssl.CertificateLocation,
			"ssl.ca.location":          client.ssl.CaLocation,
			"compression.type":         "lz4",
			"go.delivery.reports":      true,
			"go.produce.channel.size":  10,
			"go.events.channel.size":   10,
			"go.batch.producer":        true,
			"request.timeout.ms":       client.settings.RequestTimeout,
		}
	} else {
		config = kafka.ConfigMap{
			"bootstrap.servers": client.brokers,
		}
	}
	var err error
	if producerSingleton == nil {
		producerSingleton, err = kafka.NewProducer(&config)
		if err != nil {
			return nil, err
		}
	}

	return producerSingleton, nil
}

func (client *Client[T]) Produce(ctx context.Context, key string, value *T) error {
	log := utils_context.GetLogFromContext(ctx, constant.ApplicationLayer, "producer.Produce")

	kafkaProducer, err := client.getKafkaProducer()
	if err != nil {
		log.Error("Error in connection producer Client: %s", err.Error())
		return err
	}
	jsonByte, errMarshal := json.Marshal(value)
	if errMarshal != nil {
		return errMarshal
	}
	client.producer = kafkaProducer
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &client.topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          jsonByte,
	}
	if err := client.producer.Produce(kafkaMessage, nil); err != nil {
		log.Error("Error in Producer Client: %s", err.Error())
		return err
	}
	event := <-client.producer.Events()
	switch ev := event.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			log.Error("Failed to deliver message: %v\n", ev.TopicPartition)
		} else {
			log.Info("Successfully produced record to topic %s partition [%d] @ offset %v\n",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
		}
	}
	client.Close()
	return nil
}
func (client *ClientDisabled[T]) Produce(ctx context.Context, key string, value *T) error {
	log := utils_context.GetLogFromContext(ctx, constant.ApplicationLayer, "producer.Produce")
	log.Info("Init in producer disabled key: [%s]  value: %#v", key, value)
	return nil
}

func (client *Client[T]) Close() {
	if client.producer != nil {
		client.producer.Flush(60 * 1000)
	}
}
