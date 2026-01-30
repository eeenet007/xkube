package models

import (
	//"database/sql"
	"fmt"
	"log"

	//"os"

	lib "xkube/xadmin/src/lib"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var o orm.Ormer

// 数据库连接
func Connect() {
	var dsn string
	db_type, _ := beego.AppConfig.String("db_type")
	db_host, _ := beego.AppConfig.String("db_host")
	db_port, _ := beego.AppConfig.String("db_port")
	db_user, _ := beego.AppConfig.String("db_user")
	db_pass, _ := beego.AppConfig.String("db_pass")
	db_name, _ := beego.AppConfig.String("db_name")
	db_path, _ := beego.AppConfig.String("db_path")
	db_sslmode, _ := beego.AppConfig.String("db_sslmode")
	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	case "postgres":
		orm.RegisterDriver("postgres", orm.DRPostgres)
		dsn = fmt.Sprintf("dbname=%s host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_name, db_host, db_user, db_pass, db_port, db_sslmode)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
		if db_path == "" {
			db_path = "./"
		}
		dsn = fmt.Sprintf("%s%s.db", db_path, db_name)
	default:
		//beego.Critical("Database driver is not allowed:", db_type)
		log.Printf("[ERROR] mysql Connect ErrorType:%s\n", db_type)
	}
	orm.RegisterDataBase("default", db_type, dsn)
}

func SyncdbInfo() {
	lib.ShowDbInfo()
}

func GetOrmObject() orm.Ormer {
	return o
}
