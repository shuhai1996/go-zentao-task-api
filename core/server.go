package core

import (
	"fmt"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/db"
	"go-zentao-task/pkg/logging"
	"go-zentao-task/pkg/util"
	"go-zentao-task/pkg/zentaouser"
	"go-zentao-task/service/notification"
	"go-zentao-task/service/zentao"
	"math/rand"
	"os"
	"time"
)

var service = zentao.InitializeService()
var notify = notification.NewNotification() // 创建报警实体

func setup(env string) {
	config.Setup(env)
	logging.Setup(env, logging.Stdout)
	db.Setup()                        //初始化数据库
	service.User = zentaouser.Setup() //用户配置
	//rbac.Setup()
}
func getUsersByCsv() map[string]string {
	c1 := make(chan interface{})     //创建通道c1
	c2 := make(chan interface{})     //创建通道c2
	users := make(map[string]string) // 必须初始化 map, nil map 不能用来存放键值对
	file := "xxx.csv"
	lent, _ := util.LineCounter(file)
	go util.ReadCsv(file, c1) //启动一个goroutine, 读取csv
	go util.ReadCsv(file, c2) //启动另一个goroutine, 读取相同的csv
	// 用 select 实现多线程
	for len(users) < lent { //循环直到切片有四个元素
		select {
		// 接收通道 c1 的结果
		case r := <-c1:
			if r != nil {
				n := r.([]string) // 读取csv的 每一行都是[]string, 故可以 转换成 []string,
				users[n[0]] = n[0]
			}
		// 接收通道 c2 的结果
		case r := <-c2:
			if r != nil {
				n := r.([]string)
				users[n[0]] = n[0]
			}
			//default:
			//	fmt.Println("没获取到值")
		}
	}
	delete(users, "user") //去掉表头
	return users
}

func consume() {
	service.UserLogin() //模拟用户登陆
	rand.Seed(time.Now().UnixNano())
	ra := rand.Intn(300)
	fmt.Println(ra)
	time.Sleep(time.Duration(ra) * time.Second) //休眠0～300s
	es, _ := service.GetEstimateToday()         // 已用工时
	count, ids := service.ConsumeRecord(8 - es) //记录工时
	fmt.Println(count, ids)
	notify.SendNotification(service.User.Account, fmt.Sprintf("%.2f", es), fmt.Sprintf("%.2f", count), ids, nil) // 发送报警，工时转成字符串

}

func notice(users map[string]string) {
	//只输出工时信息
	consumers := make(map[string]float64)
	for _, u := range users {
		service.User.Account = u
		es, _ := service.GetEstimateToday() // 已用工时
		fmt.Println(es)
		consumers[u] = es
	}
	notify.SendNotification(service.User.Account, "", "", nil, consumers) // 发送报警

}
func RunServer(env string, quit chan os.Signal) { //服务运行
	setup(env)
	consume()                //自动填写工时
	users := getUsersByCsv() //通过csv获取需要打印工时的用户
	notice(users)
	quit <- nil
}
