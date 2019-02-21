package main

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/source"
	"log"
)

func main() {
	log.Print("[小说任务] 启动...")

	//err := router.HttpServer().Run(":8090")
	//if err != nil {
	//	log.Print(err)
	//}

	//ticker := time.NewTicker(time.Hour * 1)
	//go func() {
	//	for t := range ticker.C {
	//		log.Print("[定时器]", t)
	//		bookUpDown()
	//	}
	//}()


	bookUpDown()
}

// 定时任务
func bookUpDown() {
	// 下载热门
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
		bookRes := net.UploadBook(book)
		if 2000 == bookRes.Code || 2400 == bookRes.Code {
			xorm.BookUpload(book.Identifier)
		}
	}

	// 上传章节
	chapters := xorm.Chapters()
	for _, chapter := range chapters {
		chapterRes := net.UploadChapter(chapter)
		if 2000 == chapterRes.Code || 2400 == chapterRes.Code {
			xorm.ChapterUpload(chapter.Identifier)
		}
	}
}
