package web

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strings"
	"time"
)

type HuanYue struct {
	Url string
}

func NewHuanYue() HuanYue {
	return HuanYue{
		Url: "http://www.huanyue123.com",
	}
}

// 网站小说下载
func (hy HuanYue) Pull() {
	hy.homePage()
	hy.bookRoom()
}

// 首页
func (hy HuanYue) homePage() {
	var doc = net.GoQuery(hy.Url, true)
	doc.Find("div.books").Find("li").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a")
		text := a.Text()
		href := a.AttrOr("href", "")

		log.Print("[推荐阅读]|[最新入库小说]", text, href)
		hy.bookAll(href)
	})

	doc.Find("div.news").Find("div.bk").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("h3").Find("a")
		text := a.Text()
		href := a.AttrOr("href", "")

		log.Print("[热门小说]", text, href)
		hy.bookAll(href)
	})

	doc.Find("div.novelslist").Each(func(i int, sec *goquery.Selection) {
		sec.Find("a").Each(func(i int, sec *goquery.Selection) {
			text := sec.Text()
			href := sec.AttrOr("href", "")

			log.Print("[小说类型]", text, href)
			if text != "" {
				hy.bookAll(href)
			}
		})
	})

	doc.Find("div.col").Find("li").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("span.s2").Find("a")
		text := a.Text()
		href := a.AttrOr("href", "")

		log.Print("[最近更新小说列表]", text, href)
		hy.bookAll(href)
	})
}

// 书库
func (hy HuanYue) bookRoom() {
	api := hy.Url + "/book/"

	// 总页数
	var doc = net.GoQuery(api, true)
	max := doc.Find("div.pagelink").Find("a.last").Text()
	maxPage := util.StringToInt(max)

	for i := 1; i < maxPage; i++ {
		hy.quanBen(i)
	}
}

///////////////////////////////////////////////////   功能  /////////////////////////////////////////////////////
// 单本书籍及其列表
func (hy HuanYue) bookAll(url string) {
	book, chapters := hy.book(url)

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
			book.Keywords = `wuxia,topNovel,novel, light novel, web novel, chinese novel, korean novel, japanese novel, read light novel, read web novel, read koren novel, read chinese novel, read english novel, read novel for free, novel chapter,free,free novel`

			log.Print("[小说]", book)
			xorm.Insert(book)
		}
	}

	// 章节
	for index, simpleChapter := range chapters {
		if xorm.ChapterExist(identify, util.IntToString(index)) {
			log.Print("[章节已存在]", book.Name, simpleChapter.Title)
		} else {
			chapter := hy.chapter(simpleChapter.Domain)

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

// 书本简介及章节列表
func (huan HuanYue) book(url string) (orm.Book, []orm.Chapter) {
	var book orm.Book
	var chapterInfos []orm.Chapter

	book.Domain = url
	book.Source = "crawler"
	book.Language = "zh"
	book.Source_ctr = 3
	book.Score = 3.0

	var doc = net.GoQuery(url, true)

	// 类型
	bookType := doc.Find("div.title").Find("a").Next().Text()
	book.Type = bookType

	doc.Find("div.book_info").Find("div").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			img := sec.Find("img").AttrOr("src", "")
			book.Cover = img
		} else if i == 1 {
			name := sec.Find("h1").Text()
			book.Name = name

			sec.Find("span.item").Each(func(i int, sec *goquery.Selection) {
				if i == 0 {
					author := strings.Split(sec.Text(), "：")[1]
					book.Author = author
				} else if i == 1 {
					var status string
					value := sec.Text()
					if value == "连载中" {
						status = "2"
					} else {
						status = "1"
					}
					book.Status = status
				}
			})

			options := sec.Find("h3").Find("div.options")
			optionTxt := options.Text()
			update := strings.Split(options.Find("span.hottext").Text(), "：")[1]
			book.Last_update = update

			describe := sec.Find("h3").Text()
			describe = strings.Replace(describe, optionTxt, "", -1)
			reg, _ := regexp.Compile("(^\\s*)|(《\\S*》简介：\\s*)|(\\s*$)")
			describe = reg.ReplaceAllString(describe, "")
			book.Describe = describe
		}
	})

	doc.Find("div.book_list").Find("li").Each(func(i int, sec *goquery.Selection) {
		var chapterInfo orm.Chapter
		a := sec.Find("a")
		indexName, title := util.TitleSepatate(a.Text())
		href := a.AttrOr("href", "")

		chapterInfo = orm.Chapter{
			Idx_name: indexName,
			Title:    title,
			Domain:   href,
		}
		chapterInfos = append(chapterInfos, chapterInfo)
	})

	return book, chapterInfos
}

// 章节详情
func (huan HuanYue) chapter(url string) orm.Chapter {
	var chapter orm.Chapter

	chapter.Domain = url
	chapter.Source = "crawler"
	chapter.LastUpdate = time.Now().Unix()

	var doc = net.GoQuery(url, true)
	body := doc.Find("div.wrapper_main")

	// 章节及标题
	showTitle := body.Find("div.h1title").Text()
	indexName, title := util.TitleSepatate(showTitle)
	chapter.Idx_name = indexName
	chapter.Title = title

	// 内容
	content := body.Find("div#htmlContent").Text()
	chapter.Content = util.ChapterFilter(content)

	return chapter
}

// 每页全本的列表
func (hy HuanYue) quanBen(tp int) {
	var api = "/book/quanbu/default-0-0-0-0-2-0-%d.html"
	var url = fmt.Sprintf(hy.Url+api, tp)

	var doc = net.GoQuery(url, true)
	doc.Find("div.sitebox").Find("dl").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a").First()
		href := a.AttrOr("href", "")

		hy.bookAll(href)
	})
}
