package main

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/source"
	"log"
	"strings"
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
	xorm := orm.NewXOrm()

	// 下载小说
	hy := source.NewHuanYue()
	hy.Pull()

	// 本地连载书籍，同步最新章节
	serializes := xorm.Serialize()

	for _, book := range serializes {
		if strings.Contains(book.Domain, hy.Url) {
			hy.BookAll(book.Domain)
		}
	}

	// 上传书本
	books := xorm.Books()
	for _, book := range books {
		bookRes := net.UploadBook(book)
		if 2000 == bookRes.Code || 2400 == bookRes.Code {
			if 2000 == bookRes.Code {
				log.Print("[小说上传成功] ", book.Domain)
			} else {
				log.Print("[小说重复上传] ", book.Domain)
			}
			xorm.BookUpload(book.Identifier)
		}
	}

	// 上传章节
	chapters := xorm.Chapters()
	for _, chapter := range chapters {
		chapterRes := net.UploadChapter(chapter)
		if 2000 == chapterRes.Code || 2400 == chapterRes.Code {
			if 2000 == chapterRes.Code {
				log.Print("[章节上传成功] ", chapter.Domain, " [章节idx] ", chapter.Idx)
			} else {
				log.Print("[章节重复上传] ", chapter.Domain, " [章节idx] ", chapter.Idx)
			}

			xorm.ChapterUpload(chapter.Identifier, chapter.Idx)
		}
	}
}
