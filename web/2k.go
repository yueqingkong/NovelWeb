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

type K2 struct {
	Url string
}

func NewK2() K2 {
	return K2{
		Url: "https://www.fpzw.com",
	}
}

// 网站小说下载
func (k2 K2) Pull() {
	// 总点击榜
	visitFormat := "/top/allvisit%d_1.html"
	// 总推荐榜
	voteFormat := "/top/allvote%d_1.html"

	for i := 1; i <= 10; i++ {
		k2.normal(k2.Url + fmt.Sprintf(visitFormat, i))
		k2.normal(k2.Url + fmt.Sprintf(voteFormat, i))
	}
}

// 通用
func (k2 K2) normal(url string) {
	doc := net.GoQuery(url, true)

	doc.Find("table.sf-grid").Find("tbody").Find("tr.odd").Each(func(i int, sec *goquery.Selection) {
		a := sec.Find("a.STYLElvx")
		href := k2.Url + a.AttrOr("href", "")

		k2.BookAll(href)
	})
}

///////////////////////////////////////////////////   功能  /////////////////////////////////////////////////////
// 小说信息及列表
func (k2 K2) BookAll(url string) {
	book, chapters := k2.book(url)

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
			chapter := k2.chapter(simpleChapter.Domain)

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
func (k2 K2) book(url string) (orm.Book, []orm.Chapter) {
	var book orm.Book
	var chapters []orm.Chapter

	doc := net.GoQuery(url, true)

	// 封面
	work := doc.Find("div.work")
	bookCover := work.Find("div.bortable").Find("img").AttrOr("src", "")
	book.Cover = k2.Url + bookCover

	// 简介
	wright := work.Find("div.wright")

	// 标题
	title := wright.Find("div#title")
	h2 := title.Find("h2")
	bookName := h2.Find("a").First().Text()
	book.Name = bookName
	bookAuthor := h2.Find("em").Find("a").First().Text()
	book.Author = bookAuthor
	statusId := title.Find("div").AttrOr("id", "")
	if statusId == "lzico" {
		book.Status = "2"
	} else if statusId == "wjico" {
		book.Status = "1"
	}

	// 类型
	wright.Find("div.winfo").Find("ul").Find("li").Each(func(i int, sec *goquery.Selection) {
		if 0 == i {
			bookType := sec.Find("span").Text()
			book.Type = bookType
		} else if 3 == i {
			bookUpdate := sec.Find("span").Text()
			book.Last_update = bookUpdate
		}
	})

	// 简介
	p := wright.Find("p.Text")
	bookDescribe := p.Text()
	book.Describe = bookDescribe

	// box
	aBT := wright.Find("div#box4").Find("div#opt").Find("li#bt_1").Find("a")
	chapterHref := aBT.AttrOr("href", "")

	chapterDoc := net.GoQuery(chapterHref, true)

	var dtTotal int
	chapterDoc.Find("dl.book").Find("dt,dd").Each(func(i int, sec *goquery.Selection) {
		if sec.Is("dt") {
			dtTotal += 1
		}

		if sec.Is("dd") {
			if dtTotal == 2 {
				a := sec.Find("a")
				text := a.Text()
				href := a.AttrOr("href", "")
				idxName, title := util.TitleSepatate(text)

				chapter := orm.Chapter{Idx_name: idxName,
					Title:  title,
					Domain: chapterHref + href}
				chapters = append(chapters, chapter)
			}
		}
	})

	return book, chapters
}

// 章节
func (k2 K2) chapter(url string) orm.Chapter {
	var chapter orm.Chapter

	doc := net.GoQuery(url, true)

	// 标题
	text := doc.Find("body").Find("h2").Text()
	idxName, title := util.TitleSepatate(text)
	chapter.Idx_name = idxName
	chapter.Title = title

	// 内容
	p := doc.Find("p.Text")
	p.Find("a").Remove()
	p.Find("font").Remove()
	p.Find("strong").Remove()
	p.Find("script").Remove()
	chapter.Content = p.Text()

	return chapter
}
