package mq

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go-zentao-task-api/pkg/config"
)

// GetTopicSubscribeList {consumer:topic} consumer以 -consumer命名
func GetTopicSubscribeList() map[string]string{
	var m = make(map[string]string)
	m["action-consumer"] = "test"
	return m
}

func GetCfg() *kafka.ConfigMap{
	broker := config.Get("kafka.host")
	group := config.Get("kafka.group")
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
	kafkaconf.SetKey("group.id", group)
	kafkaconf.SetKey("broker.address.family", "v4")
	return kafkaconf
}


