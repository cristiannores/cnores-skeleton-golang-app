package consumer

import (
	"context"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"cnores-skeleton-golang-app/app/application/constant"
	settings_connection "cnores-skeleton-golang-app/app/infrastructure/kafka/settings"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"cnores-skeleton-golang-app/app/shared/utils/log"
)

type ClientInterface interface {
	Consumer(ctx context.Context) <-chan IncomingMessage
	Close() error
}

type Client struct {
	consumer      *kafka.Consumer
	consumerGroup string
	topic         string
	brokers       string
	ssl           settings_connection.SecureConnection
	settings      settings_connection.ConsumerConnection
}

type IncomingMessage struct {
	Message []byte
	Key     string
}

func NewKafkaClient(address []string, topic string, consumerGroup string, ssl settings_connection.SecureConnection, settings settings_connection.ConsumerConnection) ClientInterface {
	log.Info("Creating consumer instance addr: %s  topic: %s consumerGroup: %s", address, topic, consumerGroup)
	return &Client{
		brokers:       strings.Join(address, ","),
		consumerGroup: consumerGroup,
		topic:         topic,
		ssl:           ssl,
		settings:      settings,
	}
}

func (client *Client) getKafkaConsumer() (*kafka.Consumer, error) {
	config := kafka.ConfigMap{}

	if client.ssl.IsEnabled {
		config = kafka.ConfigMap{
			"group.id":                 client.consumerGroup,
			"bootstrap.servers":        client.brokers,
			"security.protocol":        "SSL",
			"ssl.key.location":         client.ssl.KeyLocation,
			"ssl.key.password":         client.ssl.KeyPassword,
			"ssl.certificate.location": client.ssl.CertificateLocation,
			"ssl.ca.location":          client.ssl.CaLocation,
			"session.timeout.ms":       client.settings.SessionTimeout,
			"max.poll.interval.ms":     client.settings.MaxPollInterval,
			"heartbeat.interval.ms":    client.settings.IntervalHeartbeat,
			"fetch.max.bytes":          client.settings.FetchMaxBytes,
			"fetch.message.max.bytes":  client.settings.FetchMessageMaxBytes,
		}
	} else {
		config = kafka.ConfigMap{
			"bootstrap.servers": client.brokers,
			"group.id":          client.consumerGroup,
		}
	}
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (client *Client) Consumer(ctx context.Context) <-chan IncomingMessage {
	log := utils_context.GetLogFromContext(ctx, constant.ApplicationLayer, "consumer.Consumer")

	kafkaConsumer, err := client.getKafkaConsumer()
	if err != nil {
		log.Error("Error in connection consumer Client: %s", err.Error())
		return nil
	}
	client.consumer = kafkaConsumer
	messages := make(chan IncomingMessage)
	go func() {
		defer close(messages)

		for {
			err := client.consumer.Subscribe(client.topic, nil)
			if err != nil {
				log.Error("Error in Consumer Subscribe: %s", err.Error())
			} else {
				incomingMessage, kafkaError := client.consumer.ReadMessage(-1)
				if kafkaError != nil {
					log.Error("Error in Consumer Client: %s", kafkaError.Error())
				} else {
					var e IncomingMessage
					e.Message = incomingMessage.Value
					e.Key = string(incomingMessage.Key)

					messages <- e
				}
			}
		}
	}()

	return messages
}

func (client *Client) Close() error {
	return client.consumer.Close()
}
