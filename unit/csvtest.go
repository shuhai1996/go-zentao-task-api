package unit

import (
	"fmt"
	"go-zentao-task/pkg/util"
	"time"
)

func TestCsv(f string) { // 测试一个1g左右100w行的csv文件
	startT := time.Now()
	count, _ := util.LineCounter(f)
	tc := time.Since(startT) //计算耗时
	fmt.Println(count)
	fmt.Printf("time1 cost = %v\n", tc)
	startT = time.Now()
	c1 := make(chan interface{}) //创建通道c1
	go util.ReadCsv(f, c1)       //启动一个goroutine, 读取csv
	count = 0
	for _ = range c1 { // 从通道 c 中接收,遍历 通道c
		count++
	}
	tc = time.Since(startT) //计算耗时
	fmt.Println(count)
	fmt.Printf("time2 cost = %v\n", tc)
	startT = time.Now()
	all := util.ReadAll(f) //一次性读取
	count = len(all)
	tc = time.Since(startT) //计算耗时
	fmt.Println(count)
	fmt.Printf("time3 cost = %v\n", tc)
	startT = time.Now()
	all = util.ReadCsv2(f) //不用通道读取
	count = len(all)
	tc = time.Since(startT) //计算耗时
	fmt.Println(count)
	fmt.Printf("time4 cost = %v\n", tc)
	startT = time.Now()              //双通道读取
	c1 = make(chan interface{})      //创建通道c1
	c2 := make(chan interface{})     //创建通道c2
	users := make(map[string]string) // 必须初始化 map, nil map 不能用来存放键值对
	go util.ReadCsv(f, c1)           //启动一个goroutine, 读取csv
	go util.ReadCsv(f, c2)           //启动另一个goroutine, 读取相同的csv
	// 用 select 实现多线程
	for len(users) < count { //循环直到切片有四个元素
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
		default:
			//fmt.Println("没获取到值")
		}
	}
	tc = time.Since(startT) //计算耗时
	fmt.Println(len(users))
	fmt.Printf("time5 cost = %v\n", tc)
}
