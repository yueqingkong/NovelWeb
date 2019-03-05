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

type Q3 struct {
	Url string
}

func NewQ3() Q3 {
	return Q3{
		Url: "https://www.3qzone.com",
	}
}

// 网站小说下载
func (q3 Q3) Pull() {
	q3.normal(q3.Url + "/xiuzhenxiaoshuo/")
}

// 通用
func (q3 Q3) normal(url string) {
	doc := net.GoQuery(url, true)

	main := doc.Find("div#main")
	// 热门
	main.Find("div#hotcontent").Find("div.item").Each(func(i int, sec *goquery.Selection) {
		aFirst := sec.Find("a").First()
		href := aFirst.AttrOr("href", "")
		q3.BookAll(href)
	})

	// 更新列表
	main.Find("div#newscontent").Find("div.l").Find("li").Each(func(i int, sec *goquery.Selection) {
		aFirst := sec.Find("a")
		href := aFirst.AttrOr("href", "")
		q3.BookAll(href)
	})

	// 排行榜
	main.Find("div#newscontent").Find("div.r").Find("li").Each(func(i int, sec *goquery.Selection) {
		aFirst := sec.Find("a")
		href := aFirst.AttrOr("href", "")
		q3.BookAll(href)
	})
}

///////////////////////////////////////////////////   功能  /////////////////////////////////////////////////////
// 小说信息及列表
func (q3 Q3) BookAll(url string) {
	book, chapters := q3.book(url)

	log.Print(book)
	xorm := orm.XOrm{}

	// 书本信息
	identify := util.MD5(book.Domain + book.Name)
	if xorm.BookExist(identify) {
		log.Print("[小说 Book 已存在]", book.Name)
	} else {
		filePath := "covers/" + identify + ".jpg"
		log.Print("[cover]", book.Cover)
		util.FileDownload(filePath, book.Cover)

		var fileResult net.UpFileResult
		net.UploadFile(filePath, &fileResult)

		if fileResult.Code == 2000 { // 封面上传成功
			book.Cover = fileResult.Data.URL

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

				//j, _ := json.Marshal(book)
				//log.Print(string(j))

				xorm.Insert(book)
			}
		} else {
			log.Print("[封面]上传失败", book.Name)
		}
	}

	// 章节
	for index, simpleChapter := range chapters {
		index++
		if xorm.ChapterExist(identify, util.IntToString(index)) {
			log.Print("[章节已存在]", book.Name, simpleChapter.Title)
		} else {
			chapter := q3.chapter(simpleChapter.Domain)

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

				//j, _ := json.Marshal(chapter)
				//log.Print(string(j))

				xorm.Insert(chapter)
			}
		}
	}
}

// 书本简介及章节列表
func (q3 Q3) book(url string) (orm.Book, []orm.Chapter) {
	var book orm.Book
	var chapters []orm.Chapter

	doc := net.GoQuery(url, true)
	body := doc.Find("body")

	// 类型
	top := body.Find("div.con_top").Text()
	top = strings.TrimSpace(top)
	bookType := strings.Split(top, ">")[1]
	book.Type = bookType

	// info
	maininfo := body.Find("div#maininfo")
	img := maininfo.Find("img")
	bookCover := img.AttrOr("src", "")
	book.Cover = bookCover

	info := maininfo.Find("div#info")
	bookName := info.Find("h1").Text()
	book.Name = bookName
	info.Find("p").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			aTxt := sec.Find("a").Text()
			book.Author = aTxt
		}
	})
	bookDescribe := info.Find("div#intro").Text()
	book.Describe = bookDescribe

	// 列表
	body.Find("div.box_con").Find("div#list").Find("dd").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a")
		href := a.AttrOr("href", "")

		chapter := orm.Chapter{
			Domain: url + href,
		}
		chapters = append(chapters, chapter)
	})

	return book, chapters
}

// 章节
func (q3 Q3) chapter(url string) orm.Chapter {
	var chapter orm.Chapter

	doc := net.GoQuery(url, true)
	box := doc.Find("div.content_read").Find("div#box_con")

	// 标题
	t := box.Find("div.bookname").Find("h1").Text()
	_, title := util.TitleSepatate(t)
	chapter.Domain = url
	chapter.Title = title

	// 内容
	content := box.Find("div#content").Text()
	content = strings.Replace(content, "           一秒记住【3Q中文网 www.3qzone.com】，精彩小说无弹窗免费阅读！", "", -1)
	chapter.Content = content
	return chapter
}
