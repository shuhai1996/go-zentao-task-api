package main

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go-zentao-task-api/pkg/mq"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type Consumer struct {
	confMap  *kafka.ConfigMap
	name string
	topic string
}

func (c *Consumer) init() error{
	c.confMap = mq.GetCfg()
	subMapping := mq.GetTopicSubscribeList()
	if subMapping[c.name] == "" {
		return errors.New("topic subscribe not matched:" + c.name)
	}
	return nil
}

func (c *Consumer) run(child interface{}) {
	ref := reflect.ValueOf(child) //通过反射调用子类对象
	err:=c.init()//先进行初始化
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	log.Println("Starting a new kafka consumer" + c.name)
	consumer, err := kafka.NewConsumer(c.confMap)
	if err != nil {
		log.Panicf("Error creating consumer: %v", err)
		return
	}
	defer consumer.Close()
	err = consumer.Subscribe(c.topic, nil)
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
				method := ref.MethodByName("DealMessage")
				if method.IsValid() { //子类方法存在调用子类方法,否则实现父类方法
					params := make([]reflect.Value, 1) // 设置参数，参数是一个 Value 的 slice
					params[0] = reflect.ValueOf(msg)
					method.Call(params) //调用方法
				} else {
					c.DealMessage(msg)
				}
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	if err = consumer.Close(); err != nil {
		log.Panicf("Error closing consumer: %v", err)
	}
}

//DealMessage 处理消息
func (c *Consumer) DealMessage(msg *kafka.Message) {
	fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	fmt.Println("parent")
}



