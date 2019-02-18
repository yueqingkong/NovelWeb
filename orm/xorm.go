package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

type XOrm struct {
}

type Device struct {
	Id    int32
	Name  string `xorm:"varchar(255)"`
	Token string `xorm:"varchar(255) notnull unique"`
}

type Book struct {
	Identifier  string  `xorm:"varchar(255)" json:"identifier"`
	Name        string  `xorm:"varchar(255)" json:"name"`
	Domain      string  `xorm:"varchar(255)" json:"domain"`
	Cover       string  `xorm:"varchar(255)" json:"cover"`
	Source      string  `xorm:"varchar(255)" json:"source"`
	Describe    string  `xorm:"varchar(5000)" json:"describe"`
	Author      string  `xorm:"varchar(255)" json:"author"`
	Type        string  `xorm:"varchar(255)" json:"type"`
	Last_update string  `xorm:"varchar(255)" json:"last_update"`
	Language    string  `xorm:"varchar(255)" json:"language"`
	Source_ctr  int64   `xorm:"bigint" json:"source_ctr"`
	Ctr         int64   `xorm:"bigint" json:"ctr"`
	Score       float32 `xorm:"float" json:"score"`
	Finish      string  `xorm:"varchar(255)" json:"finish"` //章节解析完成时间
	UpTime      string  `xorm:"varchar(255)" json:"uptime"` //上传时间,保存的时候该时间为空。上传成功后，设置为更新章节时间
}

type Chapter struct {
	Identifier string `xorm:"varchar(255)" json:"identifier"`
	Idx        int    `xorm:"int" json:"idx"`               //索引序列号
	Idx_name   string `xorm:"varchar(255)" json:"idx_name"` //索引名，第一章，第二章等
	Title      string `xorm:"varchar(255)" json:"title"`    //标题
	Content    string `xorm:"text" json:"content"`  //内容
	Source     string `xorm:"varchar(255)" json:"source"`   //来源 crawler
	Domain     string `xorm:"varchar(255)" json:"domain"`
	UpTime     string `xorm:"varchar(255)" json:"uptime"` //上传时间,保存的时候该时间为空。上传成功后，设置为更新章节时间
}

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@tcp(localhost:3306)/book?charset=utf8")
	if err != nil {
		log.Print(err)
	}

	engine.ShowSQL(true)
	err = engine.Sync2(new(Device), new(Book), new(Chapter))
	if err != nil {
		log.Print(err)
	}
}

func (xorm XOrm) Insert(i interface{}) {
	_, err := engine.Insert(i)
	if err != nil {
		log.Print(err)
	}
}

// 小说是否存在
func (xorm XOrm) BookExist(identify string) bool {
	_, err := engine.Where("identifier=?", identify).Exec(&Book{})
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

// 章节是否存在
func (xorm XOrm) ChapterExist(identify string, title string) bool {
	_, err := engine.Where("identifier=? and title=?", identify, title).Exec(&Chapter{})
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

func (xorm XOrm) Deveices() []Device {
	devices := make([]Device, 0)
	err := engine.Find(&devices)
	if err != nil {
		log.Print(err)
	}
	return devices
}
