package consumer

import (
	"context"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers"

	"cnores-skeleton-golang-app/app/application/constant"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
)

type Consumer interface {
	Consume(ctx context.Context)
}

type ConsumerBuilder = func(c ClientInterface, fromKafka handlers.FromKafkaInterface) Consumer

type ConsumerTopic string

func (consumerTopic ConsumerTopic) String() string {
	return string(consumerTopic)
}

type consumerClient struct {
	consumer ClientInterface
	input    handlers.FromKafkaInterface
}

func NewConsumer(client ClientInterface, input handlers.FromKafkaInterface) Consumer {
	return &consumerClient{client, input}
}

func (c *consumerClient) Consume(ctx context.Context) {
	topic := ctx.Value(ConsumerTopic("topic")).(string)
	log := utils_context.GetLogFromContext(ctx, constant.ApplicationLayer, "consumer.Consume")

	for messageInfo := range c.consumer.Consumer(ctx) {
		err := c.input.FromKafka(messageInfo.Message)
		if err != nil {
			log.Error("[%s] Error processing input message from kafka: %s", topic, err.Error())
		}
		log.Info("[%s] message was consumed successfully", topic)
	}
}
