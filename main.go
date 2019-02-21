package main

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/source"
	"log"
)

func main() {
	log.Print("start...")

	//err := router.HttpServer().Run(":8090")
	//if err != nil {
	//	log.Print(err)
	//}

	//links := []string{"http://www.huanyue123.com/book/50/50083/"}

	links := []string{"http://www.huanyue123.com/book/50/50083/",
		"http://www.huanyue123.com/book/52/52260/",
		"http://www.huanyue123.com/book/49/49221/",
		"http://www.huanyue123.com/book/5/5544/",
		"http://www.huanyue123.com/book/11/11430/",
		"http://www.huanyue123.com/book/20/20125/",
		"http://www.huanyue123.com/book/42/42935/",}
	hy := source.NewHuanYue()

	for _, link := range links {
		hy.BookAll(link)
	}

	// 上传书本
	xorm := orm.NewXOrm()
	books := xorm.Books()
	for _, book := range books {
		net.UploadBook(book)
	}

	// 上传章节
	chapters := xorm.Chapters()
	for _, chapter := range chapters {
		net.UploadChapter(chapter)
	}
}

func TestTranslate() {
	log.Print(net.Translate(" 老鹰吃小鸡结果"))
}
