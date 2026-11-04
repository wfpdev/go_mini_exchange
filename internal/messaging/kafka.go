package messaging

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

const OrderTopic = "create-order"

type OrderCreatedEvent struct {
	OrderID int64 `json:"order_id"`
}

type KafkaOrderPublisher interface {
	Publish(orderID int64) error
}

type KafkaOrderPublisherImpl struct {
	writer *kafka.Writer
}

func NewKafkaOrderPublisher(brokers []string, topic string) *KafkaOrderPublisherImpl {
	return &KafkaOrderPublisherImpl{writer: &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}}
}

func (p *KafkaOrderPublisherImpl) Publish(orderID int64) error {
	payload, err := json.Marshal(OrderCreatedEvent{OrderID: orderID})
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.FormatInt(orderID, 10)),
		Value: payload,
	})
}

func (p *KafkaOrderPublisherImpl) Close() error {
	return p.writer.Close()
}

func StartKafkaOrderConsumer(brokers []string, topic string, handler func(orderID int64) error) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "mini-exchange-order-consumer",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	go func() {
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka consumer read error: %v", err)
				continue
			}

			var event OrderCreatedEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Kafka consumer decode error: %v", err)
				continue
			}
			if err := handler(event.OrderID); err != nil {
				log.Printf("Kafka consumer process error: %v", err)
			}
		}
	}()

	return reader
}
