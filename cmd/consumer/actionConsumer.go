package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go-zentao-task-api/model/zentao"
	"go-zentao-task-api/pkg/config"
	"go-zentao-task-api/pkg/elasticsearch"
)

type Action struct {
	*Consumer
}

func main()  {
	config.Setup("development")
	elasticsearch.Setup()
	a := &Action{
		&Consumer{
			name:"action-consumer",
			topic:"zt_action_topic",
		},
	}
	a.run(a)
}
var es = zentao.NewEsAction()

func (a *Action) DealMessage(msg *kafka.Message) {
	m:= make(map[string]interface{})
	fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	fmt.Println("child")
	json.Unmarshal(msg.Value, &m) // json 转map
	data:= m["data"]
	t:=m["type"]
	if t.(string) == "INSERT" && data!=nil{ //行为插入的动作才处理，其他都废弃
		for _,v:=range data.([]interface{}) {
			_,err := es.Create(v)
			if err!= nil {
				fmt.Println("insert error"+err.Error())
			}
		}
	} else {
		fmt.Println("Message deprecated")
	}
}