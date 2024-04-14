package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jay-SP/movieapplication/rating/pkg/model"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer, topic}, nil
}

//Ingest starts ingestion from Kafka and returns a chanel containing rating events
//representing the data consumed from the topic.

func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}
			//The value of -1 is specific to Kafka and means that we will always
			//consume from the beginning of the topic, reading all existing messages.
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Println("Consumer error: " + err.Error())
				continue
			}
			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Unmarshal error: " + err.Error())
				continue
			}
			ch <- event
		}
	}()
	return ch, nil
}
