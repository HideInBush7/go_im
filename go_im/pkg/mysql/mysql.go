package mysql

import (
	"github.com/HideInBush7/go_im/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type mysqlConf struct {
	Dsn          string `json:"dsn"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

var db *sqlx.DB

func init() {
	// 从配置文件(viper)中key为'mysql'获取配置
	var conf = mysqlConf{}
	config.Register(`mysql`, &conf)

	var err error
	db, err = sqlx.Connect(`mysql`, conf.Dsn)
	if err != nil {
		panic(err)
	}

	if conf.MaxIdleConns <= 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns <= 0 {
		conf.MaxIdleConns = 20
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
}

func GetInstance() *sqlx.DB {
	return db
}
