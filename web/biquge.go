package web

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

type Biquge struct {
	Url string
}

func NewBiquge() Biquge {
	return Biquge{
		Url: "http://www.xbiquge.la",
	}
}

// 网站小说下载
func (bi Biquge) Pull() {
	api := fmt.Sprintf("%s%s", bi.Url, "/xiaoshuodaquan/")
	doc := net.GoQuery(api, false)

	doc.Find("div#main").Find("li").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a").First()
		href := a.AttrOr("href", "")

		bi.bookAll(href)
	})
}

///////////////////////////////////////////////////   功能  /////////////////////////////////////////////////////
// 小说信息及列表
func (bi Biquge) bookAll(url string) {
	book, chapters := bi.book(url)

	log.Print(book)
	xorm := orm.XOrm{}

	// 书本信息
	identify := util.MD5(book.Domain + book.Name)
	if xorm.BookExist(identify) {
		log.Print("[小说 Book 已存在]", book.Name)
	} else {
		filePath := "covers/" + identify + ".jpg"
		util.FileDownload(filePath, book.Cover)

		var fileResult net.UpFileResult
		net.UploadFile(filePath, &fileResult)

		if fileResult.Code == 2000 {
			book.Cover = fileResult.Data.URL
		}

		transBookName := net.Translate(book.Name)
		transBookDesc := net.Translate(book.Describe)
		transBookAuthor := net.Translate(book.Author)
		transBookType := net.Translate(book.Type)

		// 翻译失败
		if transBookName == "" || transBookDesc == "" || transBookAuthor == "" {
			log.Print("[小说 Book 翻译失败]", book, transBookName, transBookDesc, transBookAuthor)

			if transBookName == "" {
				log.Print("[书名为空]", book.Name, "==", transBookName)
			} else if transBookDesc == "" {
				log.Print("[简介为空]", book.Describe, "==", transBookDesc)
			} else if transBookAuthor == "" {
				log.Print("[作者为空]", book.Author, "==", transBookAuthor)
			}
		} else {
			book.Identifier = identify
			book.Name = transBookName
			book.Describe = transBookDesc
			book.Author = transBookAuthor
			book.Type = transBookType
			book.Index = strings.Replace(transBookName, " ", "-", -1)
			book.Translate = "2"
			book.Domain = url
			book.Source = "crawler"
			book.Language = "zh"
			book.Source_ctr = 3
			book.Score = 3.0
			book.Keywords = `wuxia,topNovel,novel, light novel, web novel, chinese novel, korean novel, japanese novel, read light novel, read web novel, read koren novel, read chinese novel, read english novel, read novel for free, novel chapter,free,free novel`

			log.Print("[小说]", book)
			xorm.Insert(book)
		}
	}

	// 章节
	for index, simpleChapter := range chapters {
		index++
		if xorm.ChapterExist(identify, util.IntToString(index)) {
			log.Print("[章节已存在]", book.Name, simpleChapter.Title)
		} else {
			chapter := bi.chapter(simpleChapter.Domain)

			transTitle := net.Translate(chapter.Title)
			transContent := net.Translate(chapter.Content)

			// 章节翻译失败
			if transTitle == "" || transContent == "" {
				log.Print("[小说章节信息翻译失败]", simpleChapter, simpleChapter.Domain, transTitle, transContent)

				if transTitle == "" {
					log.Print("[章节标题为空]", chapter.Title, "==", transTitle)
				} else if transContent == "" {
					log.Print("[章节内容为空]", chapter.Content, "==", transContent)
				}
			} else {
				chapter.Idx = index
				chapter.Identifier = identify
				chapter.Idx_name = fmt.Sprintf("Chapter %d", index)
				chapter.Title = transTitle
				chapter.Content = transContent
				chapter.Domain = simpleChapter.Domain
				chapter.LastUpdate = time.Now().Unix()
				chapter.Index = strings.Replace(transTitle, " ", "-", -1)
				chapter.Source = "crawler"
				chapter.Keywords = `wuxia,topNovel,novel, light novel, web novel, chinese novel, korean novel, japanese novel, read light novel, read web novel, read koren novel, read chinese novel, read english novel, read novel for free, novel chapter,free,free novel`
				chapter.Translate = "2"
				transBookName := net.Translate(book.Name)
				chapter.BookIndex = strings.Replace(transBookName, " ", "-", -1)

				log.Print("[章节]", book.Name, chapter)
				xorm.Insert(chapter)
			}
		}
	}
}

// 小说
func (bi Biquge) book(url string) (orm.Book, []orm.Chapter) {
	doc := net.GoQuery(url, false)

	var book orm.Book
	var chapters []orm.Chapter

	book.Domain = url

	doc.Find("div.con_top").Find("a").Each(func(i int, sec *goquery.Selection) {
		if i == 2 {
			book.Type = sec.Text()
		}
	})

	cover := doc.Find("div#sidebar").Find("img").AttrOr("src", "")
	book.Cover = cover

	maininfo := doc.Find("div#maininfo")
	info := maininfo.Find("div#info")

	name := info.Find("h1").Text()
	book.Name = name

	info.Find("p").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			author := strings.Split(sec.Text(), "：")[1]
			book.Author = author
		} else if i == 2 {
			update := strings.Split(sec.Text(), "：")[1]
			book.Last_update = update
		}
	})

	maininfo.Find("div#intro").Find("p").Each(func(i int, sec *goquery.Selection) {
		if i == 1 {
			introTxt := sec.Text()
			book.Describe = introTxt
		}
	})

	doc.Find("div#list").Find("dd").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a")
		href := bi.Url + a.AttrOr("href", "")

		chapter := orm.Chapter{Domain: href}
		chapters = append(chapters, chapter)
	})

	return book, chapters
}

// 章节
func (bi Biquge) chapter(url string) orm.Chapter {
	var chapter orm.Chapter

	doc := net.GoQuery(url, false)
	content := doc.Find("div.content_read")
	h1 := content.Find("div.bookname").Find("h1").Text()
	_, title := util.TitleSepatate(h1)
	chapter.Title = title

	conTxt := content.Find("div#content").Text()
	chapter.Content = conTxt
	return chapter
}
