## go-zentao-task-api
跟go-zentao-task 项目相比，gorm从1.9升为2.0,新增redis、elasticsearch，es版本为7.16
#### 配置
根目录下创建conf.ini文件，添加如下配置
```
[development]
#mysql
db.host =localhost
db.port =3306
db.username =root
db.password =
db.database =
#redis
redis.host =127.0.0.1
redis.port =6379
#es
es.host=127.0.0.1
es.port=9200
es.user=
es.pass=
#kafka
kafka.host=127.0.0.1
kafka.port=9092
```
#### 包管理
新的项目，复制完代码以后，全局替换go-zentao-task-api为你的项目名称，删除go.mod和go.sum文件，执行下列命令：
```
go mod init go-zentao-task-api
go mod vendor //将下载到GOPATH的包复制到当前项目的vendor目录下
go mod tidy //整理包依赖
```
#### 运行本项目
执行如下命令即可运行
```
go build main.go
./main
```

> 参考链接
> 
> [GORM 2.0 使用教程(中文文档)](https://www.bookstack.cn/read/gorm-2.0/docs-index.md)
> 
> [https://gin-gonic.com/](https://gin-gonic.com/)
> 
> [elasticsearch 7.16 guide](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/index.html)
> 
> [confluent-kafka-go-demo](https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/consumer_example/consumer_example.go)
