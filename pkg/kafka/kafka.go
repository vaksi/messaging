package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// Kafka
type Kafka struct {
	addrs          []string
	producerConfig *sarama.Config
	consumerConfig *cluster.Config
}

// NewKafka creates a new kafka
func NewKafka(addrs []string) *Kafka {
	consumerConfig := cluster.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Group.Return.Notifications = true
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	return &Kafka{
		producerConfig: producerConfig,
		consumerConfig: consumerConfig,
		addrs:          addrs,
	}
}

// SendMessage produces a given message, and returns only when it either has
// succeeded or failed to produce. It will return an error if the message failed to produce.
func (k *Kafka) SendMessage(topic string, key []byte, data []byte) error {
	// Creates new sync producer
	producer, err := sarama.NewSyncProducer(k.addrs, k.producerConfig)
	if err != nil {
		return err
	}
	defer producer.Close()

	// Send message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

// Consume methods use this callback function to handle Messages
// from the read channel that are returned by the broker.
type MessageHandler func(key []byte, value []byte)

// Consume method use this callback function to handle errors that occur during
// offset management. This allow you to implement any custom error handling.
// Usually people log the error.
type ErrorHandler func(error)

// Consume methods use this callback function to handle Notifications
// that occur during consumer re-balancing. This allows you to handle
// custom notification handling. Usually people log the notification if necessary.
type NotificationHandler func(interface{})

// Consume consumes messages from the read channel that are returned by
// the broker. The message will be handled by the MessageHandler, error
// will be handled by the ErrorHandler and the notification will be
// handled by NotificationHandler
func (k *Kafka) Consume(groupID string, topics []string, msgHandler MessageHandler, errHandler ErrorHandler, ntfHandler NotificationHandler, done chan bool) error {
	// Creates new consumer
	consumer, err := cluster.NewConsumer(k.addrs, groupID, topics, k.consumerConfig)
	if err != nil {
		return err
	}
	defer consumer.Close()

	// Handle errors
	go func() {
		for err := range consumer.Errors() {
			errHandler(err)
		}
	}()

	// Handle notifications
	go func() {
		for ntf := range consumer.Notifications() {
			ntfHandler(ntf)
		}
	}()

	// Handle messages
	for {
		select {
		case msg := <-consumer.Messages():
			msgHandler(msg.Key, msg.Value)
			consumer.MarkOffset(msg, "")
		case <-done:
			return nil
		}
	}

	return nil
}
