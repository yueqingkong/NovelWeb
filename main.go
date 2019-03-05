package main

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/web"
	"log"
	"strings"
	"time"
)

func main() {
	log.Print("[小说任务] 启动...")
	go func() {
		log.Print("[小说] 爬取网站...")
		bookUpDown()

		bookSerializate()
	}()

	//ticker := time.NewTicker(time.Second * 10)
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			log.Print("[定时器] ", t)
			bookUpload()
		}
	}()

	ch := make(chan string)
	<-ch
}

// 定时任务
func bookUpDown() {
	hy := web.NewHuanYue()
	hy.Pull()
	bi := web.NewBiquge()
	bi.Pull()
	dd:=web.NewDingDian()
	dd.Pull()
<<<<<<< HEAD
=======
	q3:=web.NewQ3()
	q3.Pull()
>>>>>>> a2ff52f4b4d3adee8e8fb57294df5fa156c307b0
}

// 同步连载最新章节
func bookSerializate() {
	xorm := orm.NewXOrm()

	hy := web.NewHuanYue()
	bi:=web.NewBiquge()

	serializes := xorm.Serialize()
	for _, book := range serializes {
		if strings.Contains(book.Domain, hy.Url) {
			hy.BookAll(book.Domain)
		}
		if strings.Contains(book.Domain, bi.Url) {
			bi.BookAll(book.Domain)
		}
	}
}

// 上传到server
func bookUpload() {
	xorm := orm.NewXOrm()

	//上传书本
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
