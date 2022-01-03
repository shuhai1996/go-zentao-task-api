package unit

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	broker = "127.0.0.1:9092"//实例地址
	topic = "zt_action_topic"
)

func TestConsumer()  {
	log.Println("Starting a new kafka consumer")
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request": "true",
		"auto.offset.reset": "earliest",
		"heartbeat.interval.ms": 3000,
		"session.timeout.ms": 30000,
		"max.poll.interval.ms": 120000,
		"fetch.max.bytes": 1024000,
		"max.partition.fetch.bytes": 256000,
	}
	kafkaconf.SetKey("bootstrap.servers", broker)
	kafkaconf.SetKey("group.id", "test-consumer-group")
	kafkaconf.SetKey("broker.address.family", "v4")
	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		log.Panicf("Error creating consumer: %v", err)
		return
	}

	defer consumer.Close()
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Panicf("Error subscribe consumer: %v", err)
		return
	}

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Consumer error: %v (%v)", err, msg)
			} else {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			}
		}
	}()

	//sigterm := make(chan os.Signal, 1)
	//signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	//select {
	//case <-sigterm:
	//	log.Println("terminating: via signal")
	//}
	if err = consumer.Close(); err != nil {
		log.Panicf("Error closing consumer: %v", err)
	}
}

func Test2()  {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              "test-consumer-group",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
	}

	defer c.Close()

	fmt.Printf("Created Consumer %v\n", c)

	err = c.Subscribe(topic, nil)
	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
