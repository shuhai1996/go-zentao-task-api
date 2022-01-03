package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Action struct {
	*Consumer
}

func main()  {
	a := &Action{
		&Consumer{
			name:"action-consumer",
			topic:"zt_action_topic",
		},
	}
	a.run()
}

func (a *Action) DealMessage(msg *kafka.Message) {
	fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
}