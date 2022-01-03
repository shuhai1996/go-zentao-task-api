package zentao

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go-zentao-task-api/pkg/elasticsearch"
	"strconv"
	"time"
)
var request = make(map[string]string)

type Es struct { //es 结构体
	Index string
	Total int
	Hits interface{}
	Type string
}

type EsAction struct { // 组合方式实现继承
	*Es
}

func NewEsAction() *EsAction {
	return &EsAction{
		Es: &Es{
			Index: "zentao_action",
			Type: "_doc",
		},
	}
}

func (es *EsAction) Create(data interface{}) (interface{}, error) {
	var m = make(map[string]interface{})
	body,err := json.Marshal(data) //json 序列化
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	json.Unmarshal(body, &m) // json序列化之后转成map

	m["date"] = formatData(m["date"].(string))
	body,_ = json.Marshal(m) //再转成json
	request["bodyJson"] = string(body)
	request["index"] = es.Index
	request["id"] = m["id"].(string)
	request["type"] = es.Type
	fmt.Println(request)
	return elasticsearch.EsClient.Create(request), nil
}

func (es *EsAction) Update(data interface{}) (interface{}, error) {
	var m = make(map[string]string)
	body,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	json.Unmarshal(body, &m) // json序列化之后转成map
	request["body"] = string(body)
	request["index"] = es.Index
	request["id"] = m["id"]
	request["type"] = es.Type
	return elasticsearch.EsClient.Update(request), nil
}

func (es *EsAction) Delete(id string) (string, error) {
	request["id"] = id
	request["index"] = es.Index
	request["type"] = es.Type
	return elasticsearch.EsClient.Delete(request), nil
}

func (es *EsAction) Find(id string) (*elastic.GetResult, error) {
	request["id"] = id
	request["index"] = es.Index
	request["type"] = es.Type
	return elasticsearch.EsClient.Gets(request)
}

func formatData(t string) int64 {
	tUtc, _ := time.Parse("2006-01-02 15:04:05", t) //utc时间 string 类型转 成time 类型
	return  tUtc.UTC().Unix()
}

func (es *EsAction) List(page int, size int, query string) *elastic.SearchHits {
	request["index"] = es.Index
	request["queryString"] = query
	request["size"] = strconv.Itoa(size)
	request["page"] = strconv.Itoa(page-1)
	fmt.Println(request)
	return elasticsearch.EsClient.List(request).Hits
}