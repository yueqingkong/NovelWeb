package plat

import (
	"Novel/plat"
	"NovelWeb/orm"
	"NovelWeb/util"
	"log"
	"time"
)

type HuanYue struct {
}

func (hy HuanYue) Parser() {
	var startTop = 1 //解析页数

	mongo:=orm.Mongo{}

	maxpage, _ := plat.Huanyue.Top(1)
	website := mongo.WebSite("http://www.huanyue123.com")
	if website.WebsiteURL != "" {
		if maxpage != website.LastTop {
			startTop = website.LastTop
		} else { //解析完成
			return
		}
	}

	for ; startTop <= maxpage; startTop++ {
		hy.Top(startTop)
		time.Sleep(time.Second)
	}
}

func (hy HuanYue) Top(top int) {
	log.Print("huanyue [Top]: ", top)
	_, books := plat.Huanyue.Top(top)
	for _, book := range books {
		bookinfo, chapters := plat.Huanyue.All(book.Source) //每本书籍信息 章节信息

		//log.Print("bookinfo", bookinfo)
		identify := util.MD5(bookinfo.Domain + bookinfo.Name) //网站+书名 hash
		bookName := bookinfo.Name

		mongo:=orm.Mongo{}
		mgoBook := mongo.BookIdentify(identify)
		if mgoBook.Identifier == "" { //本地没有保存该书
			var mBook = orm.MBook{
				Identifier:  identify,
				Domain:      bookinfo.Domain,
				Name:        bookName,
				Cover:       bookinfo.Cover,
				Source:      "crawler",
				Describe:    book.Describe,
				Author:      bookinfo.Author,
				Type:        bookinfo.Type,
				Last_update: bookinfo.Update,
				Language:    "zh",
				Source_ctr:  1,
				Ctr:         1,
				Score:       1,
			}

			if mBook.Name == "" || mBook.Cover == "" || mBook.Describe == "" || mBook.Author == "" {
				log.Println("HuanYue book is null")
				continue
			} else {
				mongo.Insert(mBook)
				log.Print("HuanYue Parser", mBook)
			}
		} else {
			if mgoBook.Finish != "" {
				continue
			}
		}

		for key, chapter := range chapters {
			mgoChapter := mongo.ChapterIdentifyName(identify, chapter.Title)
			if mgoChapter.Identifier == "" { //本地没有保存该章节
				chapterDetail := plat.Huanyue.Chapter(chapter.Source)

				var mChapter = orm.MChapter{
					Identifier: identify,
					Idx:        key + 1,
					Idx_name:   chapterDetail.Index,
					Title:      chapterDetail.Title,
					Content:    chapterDetail.Content,
					Source:     "crawler",
					Domain:     chapterDetail.Domain,
				}

				if mChapter.Title == "" {
					log.Println("HuanYue chapter is null", chapterDetail)
					continue
				} else {
					mongo.Insert(mChapter)
					log.Print("HuanYue Parser", mChapter)
				}
			}

			mongo.LocalBookFinsih(identify)
		}
		mongo.Insert(orm.Website{WebsiteURL: "http://www.huanyue123.com", LastTop: top})
	}
}
