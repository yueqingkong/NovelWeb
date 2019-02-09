package util

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

type Device struct {
	Id    int32
	Name  string `xorm:"varchar(255)"`
	Token string `xorm:"varchar(255) notnull unique"`
}

var engine *xorm.Engine

func init() {
	connect()
}

func connect() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@tcp(localhost:3306)/token?charset=utf8")
	if err != nil {
		log.Print(err)
	}

	engine.ShowSQL(true)
	err = engine.Sync2(new(Device))
	if err != nil {
		log.Print(err)
	}
}

func Insert(i interface{}) {
	_, err := engine.Insert(i)
	if err != nil {
		log.Print(err)
	}
}

func Deveices() []Device {
	devices := make([]Device, 0)
	err := engine.Find(&devices)
	if err != nil {
		log.Print(err)
	}
	return devices
}
