package db

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/config"
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

	var orm *gorm.DB
	orm, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	orm.DB().SetMaxIdleConns(50)
	orm.DB().SetMaxOpenConns(100)
	orm.DB().SetConnMaxLifetime(300 * time.Second)
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
