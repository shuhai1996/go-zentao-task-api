package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go-zentao-task-api/pkg/config"
	"log"
	"os"
	"strconv"
)

type EsClientType struct {
	EsCon *elastic.Client
}

type esConfig struct {
	host string
	port string
	username string
	password string
}

var Timeout = "1s"        //超时时间
var EsClient *EsClientType //连接类型

func Setup() {
	esConnectPool("")
}

func esConnectPool(store string) {
	conf, err := getEsConfig(store)
	if err != nil {
		log.Fatalln(err)
	}
	EsClient = newEsPool("http://" + conf.username +":"+ conf.password + "@"+conf.host + ":" + conf.port)
}

func getEsConfig(store string) (conf *esConfig, err error) {
	conf = &esConfig{
		host: config.Get("es.host" + store),
		port: config.Get("es.port" + store),
		username: config.Get("es.user" + store),
		password: config.Get("es.pass" + store),
	}
	if conf.host == "" {
		err = errors.New("redis.host" + store + "不能为空")
		return
	}
	if conf.port == "" {
		err = errors.New("redis.port" + store + "不能为空")
		return
	}
	return
}
func newEsPool(host string) (newClient *EsClientType){
	elastic.SetSniff(false) //必须 关闭 Sniffing
	//es 配置
	var err error
	//EsClient.EsCon, err = elastic.NewClient(elastic.SetURL(host))
	esCon, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if err != nil {
		panic(err)
	}
	info, code, err := esCon.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := esCon.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	fmt.Println("conn es succ", esCon)
	newClient = &EsClientType{
		EsCon: esCon,
	}
	return
}

// CreateIndex 创建索引
func (client *EsClientType) CreateIndex(Params map[string]string) string {
	//使用字符串
	var res *elastic.IndicesCreateResult
	var err error
	res, err = client.EsCon.CreateIndex(Params["index"]).
		Body(Params["mappings"]).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	return res.Index
}

// Create 创建
func (client *EsClientType) Create(Params map[string]string) string {
	//使用字符串
	var res *elastic.IndexResponse
	var err error

	res, err = client.EsCon.Index().
		Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).BodyJson(Params["bodyJson"]).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	return res.Result

}

// Delete 删除
func (client *EsClientType) Delete(Params map[string]string) string {
	var res *elastic.DeleteResponse
	var err error

	res, err = client.EsCon.Delete().Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		println(err.Error())
	}
	return res.Result
}

// Update 修改
func (client *EsClientType) Update(Params map[string]string) string {
	var res *elastic.UpdateResponse
	var err error

	res, err = client.EsCon.Update().
		Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).
		Doc(Params["doc"]).
		Do(context.Background())

	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)
	return res.Result

}

// Gets 查找
func (client *EsClientType) Gets(Params map[string]string) (*elastic.GetResult,error) {
	//通过id查找
	var get1 *elastic.GetResult
	var err error
	if len(Params["id"]) < 0 {
		err = errors.New("param error")
		return get1, err
	}

	get1, err = client.EsCon.Get().Index(Params["index"]).Type(Params["type"]).Id(Params["id"]).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return get1, nil
}

//搜索
func (client EsClientType) Query(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = client.EsCon.Search(Params["index"]).Type(Params["type"]).Do(context.Background())
	if len(Params["queryString"]) > 0 {
		//字段相等
		q := elastic.NewQueryStringQuery(Params["queryString"])
		res, err = client.EsCon.Search(Params["index"]).Type(Params["type"]).Query(q).Do(context.Background())
	}
	if err != nil {
		println(err.Error())
	}

	//if res.Hits.TotalHits > 0 {
	//	fmt.Printf("Found a total of %d Employee \n", res.Hits.TotalHits)
	//}
	return res
}

//简单分页 可用
func (client *EsClientType) List(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	size, _ := strconv.Atoi(Params["size"])
	page, _ := strconv.Atoi(Params["page"])
	//q := elastic.NewQueryStringQuery(Params["queryString"])

	//排序类型 desc asc es 中只使用 bool 值  true or false
	//sortType := true
	//if Params["sort_type"] == "desc" {
	//	sortType = false
	//}
	//fmt.Printf(" sort info  %s,%s\n", Params["sort"],Params["sort_type"])
	if size < 0 || page < 0 {
		fmt.Printf("param error")
		return res
	}
	if len(Params["queryString"]) > 0 {
		res, err = client.EsCon.Search(Params["index"]).
			Type(Params["type"]).
			//Query(q).
			Size(size).
			From((page)*size).
			//Sort(Params["sort"], sortType).
			Timeout(Timeout).
			Do(context.Background())

	} else {
		res, err = client.EsCon.Search(Params["index"]).
			Type(Params["type"]).
			Size(size).
			From((page)*size).
			//Sort(Params["sort"], sortType).
			//SortBy(elastic.NewFieldSort("add_time").UnmappedType("long").Desc(), elastic.NewScoreSort()).
			Timeout(Timeout).
			Do(context.Background())
	}

	if err != nil {
		println("func list error:" + err.Error())
	}
	return res

}

//聚合 平均 可用
func (client *EsClientType) Aggregation(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error

	//需要聚合的指标 求平均
	avg := elastic.NewAvgAggregation().Field(Params["avg"])
	//单位时间和指定字段
	aggs := elastic.NewDateHistogramAggregation().
		Interval("day").
		Field(Params["field"]).
		//TimeZone("Asia/Shanghai").
		SubAggregation(Params["agg_name"], avg)

	res, err = client.EsCon.Search(Params["index"]).
		Type(Params["type"]).
		Size(0).
		Aggregation(Params["aggregation_name"], aggs).
		//Sort(Params["sort"],sort_type).
		Timeout(Timeout).
		Do(context.Background())

	if err != nil {
		println("func Aggregation error:" + err.Error())
	}
	println("func Aggregation here 297")

	return res

}
