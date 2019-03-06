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
	Id          int64   `xorm:"int" json:"id"`                         // 自增长id
	Identifier  string  `xorm:"varchar(255) unique" json:"identifier"` // 唯一标志
	Domain      string  `xorm:"varchar(255)" json:"domain"`            // 来源网站
	Name        string  `xorm:"varchar(255)" json:"name"`              //书名
	Cover       string  `xorm:"varchar(255)" json:"cover"`             //封面
	Source      string  `xorm:"varchar(255)" json:"source"`            //来源 crawler
	Describe    string  `xorm:"varchar(5000)" json:"describe"`         //描述
	Author      string  `xorm:"varchar(255)" json:"author"`            //作者
	Type        string  `xorm:"varchar(255)" json:"type"`              //类型
	Last_update string  `xorm:"varchar(255)" json:"last_update"`       //最近更新时间
	Language    string  `xorm:"varchar(255)" json:"language"`          //语言 en,zh
	Source_ctr  int64   `xorm:"bigint" json:"source_ctr"`              //原网站点击率
	Ctr         int64   `xorm:"bigint" json:"ctr"`                     //本站点击率
	Score       float32 `xorm:"float" json:"score"`                    //评分
	Keywords    string  `xorm:"varchar(255)"  json:"keywords"`         //seo信息
	Index       string  `xorm:"varchar(255)" json:"index"`             //索引序列号，书名索引
	Status      string  `xorm:"varchar(255)" json:"status"`            //状态 1:完成 2：连载
	Translate   string  `xorm:"varchar(255)" json:"translate"`         //状态 1:原文 2：机器翻译
	IsUpload    int     `xorm:"int" json:"is_upload"`                  //上传状态,上传成功后，更新状态位 1
}

type Chapter struct {
	Id         int64  `xorm:"int" json:"id"`                                            // 自增长id
	Identifier string `xorm:"varchar(255) unique(identifier_domain)" json:"identifier"` //书籍唯一标志
	Idx        int    `xorm:"int" json:"idx"`                                           //索引序列号
	Idx_name   string `xorm:"varchar(255)" json:"idx_name"`                             //索引名，第一章，第二章等
	Title      string `xorm:"varchar(255)" json:"title"`                                //标题
	Content    string `xorm:"mediumtext" json:"content"`                                //内容
	Source     string `xorm:"varchar(255)" json:"source"`                               //来源 crawler
	Domain     string `xorm:"varchar(255) unique(identifier_domain)" json:"domain"`     //来源网站
	LastUpdate int64  `json:"last_update"`                                              //最近更新时间
	Keywords   string `xorm:"varchar(255)" json:"keywords"`                             //seo信息
	Index      string `xorm:"varchar(255)" json:"index"`                                //章节名索引索引
	BookIndex  string `xorm:"varchar(255)" json:"book_index"`                           //书名索引
	Translate  string `xorm:"varchar(255)" json:"translate"`                            //状态 1:原文 2：机器翻译
	IsUpload   int    `xorm:"int" json:"is_upload"`                                     //上传状态,上传成功后，更新状态位 1
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

func NewXOrm() XOrm {
	return XOrm{}
}

func (xorm XOrm) Insert(i interface{}) {
	_, err := engine.Insert(i)
	if err != nil {
		log.Print(err)
	}
}

func (xorm XOrm) Books() []Book {
	var books []Book
	err := engine.SQL("select * from book where is_upload != 1;").Find(&books)
	if err != nil {
		log.Print(err)
	}
	return books
}

// 连载
func (xorm XOrm) Serialize() []Book {
	var books []Book
	err := engine.SQL("select * from book where status = 2;").Find(&books)
	if err != nil {
		log.Print(err)
	}
	return books
}

func (xorm XOrm) Chapters() []Chapter {
	var chapters []Chapter
	err := engine.SQL("select * from chapter where is_upload != 1").Find(&chapters)
	if err != nil {
		log.Print(err)
	}
	return chapters
}

// 小说是否存在
func (xorm XOrm) BookExist(identify string) bool {
	result, err := engine.SQL("select * from book where identifier = ?;", identify).Exist()
	if err != nil {
		log.Print(err)
	}
	return result
}

// 章节是否存在
func (xorm XOrm) ChapterExist(identify string, index string) bool {
	result, err := engine.SQL("select * from chapter where identifier = ? and idx = ?;", identify, index).Exist()
	if err != nil {
		log.Print(err)
	}
	return result
}

// 更新书籍 上传成功
func (xorm XOrm) BookUpload(identify string) {
	sql := "update book set is_upload = 1 where identifier = ? ;"
	_, err := engine.Exec(sql, identify)
	if err != nil {
		log.Print("[xorm] BookUpload: ", err)
	}
}

// 更新章节 上传成功
func (xorm XOrm) ChapterUpload(identify string, idx int) {
	sql := "update chapter set is_upload = 1 where identifier = ? and idx = ?;"
	_, err := engine.Exec(sql, identify, idx)
	if err != nil {
		log.Print("[xorm] ChapterUpload: ", err)
	}
}

func (xorm XOrm) Deveices() []Device {
	devices := make([]Device, 0)
	err := engine.Find(&devices)
	if err != nil {
		log.Print(err)
	}
	return devices
}
