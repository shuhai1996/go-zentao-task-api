package db

import (
	"errors"
	"go-zentao-task-api/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var Orm *gorm.DB

func Setup() {
	Orm = mysqlConnectPool("")
}

func mysqlConnectPool(store string) *gorm.DB {
	dsn, err := getDBConfig(store)
	if err != nil {
		log.Fatalln(err)
	}

	orm, er := gorm.Open(mysql.Open(dsn))
	if er != nil {
		log.Fatalln(er)
	}
	sqlDB, e := orm.DB()
	if e != nil {
		log.Fatalln(e)
	}
	sqlDB.SetMaxIdleConns(50)                   //设置连接池最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                  //设置打开连接个数
	sqlDB.SetConnMaxLifetime(300 * time.Second) //设置连接最大生命周期

	return orm
}

func getDBConfig(store string) (string, error) {
	host := config.Get("db.host" + store)
	port := config.Get("db.port" + store)
	username := config.Get("db.username" + store)
	password := config.Get("db.password" + store)
	database := config.Get("db.database" + store)

	if host == "" {
		return "", errors.New("数据库host不能为空")
	}
	if port == "" {
		return "", errors.New("数据库port不能为空")
	}
	if username == "" {
		return "", errors.New("数据库username不能为空")
	}
	if password == "" {
		return "", errors.New("数据库password不能为空")
	}
	if database == "" {
		return "", errors.New("数据库名称不能为空")
	}

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=true&loc=Local"
	return dsn, nil
}
