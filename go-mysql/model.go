package model

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB     *gorm.DB
	sqlLog log.Logger //可重写gorm.Logger
)

type sqlLogger struct {
	*gorm.Logger
}

type MysqlConfig struct {
	host     string `json:"host"`
	port     string `json:"port"`
	user     string `json:"user"`
	password string `json:"user"`
	dbname   string `json:"dbname"`
}

func InitDB(conf MysqlConfig) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.user, conf.password, "tcp", conf.host, conf.port, conf.dbname)
	DB, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	DB.LogMode(true)

	// 给db设置一个超时时间，时间小于数据库的超时时间即可
	DB.DB().SetConnMaxLifetime(100 * time.Second)

	// 用于设置最大打开的连接数，默认值为0表示不限制。
	DB.DB().SetMaxOpenConns(10)

	// 用于设置闲置的连接数
	DB.DB().SetMaxIdleConns(10)

	//自己实现日志
	DB.SetLogger(sqlLogger{})

	//初始化连接，默认创建的连接池，没有真正连接
	err = DB.DB().Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		panic(err.Error())
	}

	fmt.Println("Success to connect to mysql")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

//Acording
func (sqlLogger) Print(values ...interface{}) {
	//rewrite gorm.Logger
	sqlLog.Print(values)
}
