package source

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
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

// 关键字查询
func (hy HuanYue) SearchKeyword(keyword string, page int) (int, []orm.Book) {
	var bookInfos = make([]orm.Book, 0)

	var api = "/modules/article/search.php?searchkey=%s&page=%d"
	var url = fmt.Sprintf(hy.Url+api, util.Utf8ToGbk(keyword), page)
	log.Print("url: ", url)

	var doc = net.GoQuery(url)
	doc.Find("tbody").Find("tr").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			return
		}

		var bookInfo orm.Book
		sec.Find("td").Each(func(i int, sec *goquery.Selection) {
			if i == 0 {
				a := sec.Find("a")
				name := a.Text()
				source := a.AttrOr("href", "")

				bookInfo.Name = name
				bookInfo.Source = source
			} else if i == 1 {
				//a := sec.Find("a")
				//index, title := util.SepatateTitle(a.Text())
				//href := a.AttrOr("href", "")
				//
				//bookInfo.Chapter = ChapterInfo{
				//	Index:  index,
				//	Title:  title,
				//	Source: href,
				//}
			} else if i == 2 {
				author := sec.Text()
				bookInfo.Author = author
			} else if i == 3 {

			} else if i == 4 {
				update := sec.Text()
				bookInfo.Last_update = update
			}
		})

		bookInfos = append(bookInfos, bookInfo)
	})

	pages := doc.Find("div.pagelink").Find("em").Text()
	max, _ := strconv.Atoi(strings.Split(pages, "/")[1])
	return max, bookInfos
}

// 书本简介及章节列表
func (huan HuanYue) Book(url string) (orm.Book, []orm.Chapter) {
	var book orm.Book
	var chapterInfos []orm.Chapter

	book.Domain = url
	book.Source = "crawler"
	book.Language = "zh"
	book.Source_ctr = 3
	book.Score = 3.0

	var doc = net.GoQuery(url)

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
func (huan HuanYue) Chapter(url string) orm.Chapter {
	var chapter orm.Chapter

	chapter.Domain = url
	chapter.Source = "crawler"
	chapter.LastUpdate = time.Now().Unix()

	var doc = net.GoQuery(url)
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
func (hy HuanYue) QuanBenTop(tp int) (int, []orm.Book) {
	var bookInfos []orm.Book
	var api = "/book/quanbu/default-0-0-0-0-2-0-%d.html"
	var url = fmt.Sprintf(hy.Url+api, tp)

	var doc = net.GoQuery(url)
	doc.Find("div.sitebox").Find("dl").Each(func(i int, sec *goquery.Selection) {
		var bookInfo orm.Book
		bookInfo.Domain = url
		bookInfo.Type = "unknow"

		img := sec.Find("dt").Find("a").Find("img")
		bookInfo.Cover = img.AttrOr("src", "")

		sec.Find("dd").Each(func(i int, sec *goquery.Selection) {
			if i == 0 {
				update := sec.Find("span").Text()
				bookInfo.Last_update = update

				a := sec.Find("a")
				name := a.Text()
				bookInfo.Name = name
				href := a.AttrOr("href", "")
				bookInfo.Source = href
			} else if i == 1 {
				author := sec.Find("span").Text()
				bookInfo.Author = author
			} else if i == 2 {
				describe := sec.Text()
				bookInfo.Describe = describe
			} else if i == 3 {
				//a := sec.Find("a")
				//index, title := util.SepatateTitle(a.Text())
				//href := a.AttrOr("href", "")
				//bookInfo.Chapter = orm.Chapter{
				//	Index:  index,
				//	Title:  title,
				//	Source: url + href,
				//}
			}
		})
		bookInfos = append(bookInfos, bookInfo)
	})

	pages := doc.Find("div.pagelink").Find("em#pagestats").Text()
	var max int
	var err error

	if pages != "" {
		splarr := strings.Split(pages, "/")
		max, err = strconv.Atoi(splarr[1])
		if err != nil {
			log.Fatal(tp, splarr, err)
		}
	}

	log.Print("bookInfos: ", tp, " max: ", max, " bookInfos: ", bookInfos)
	return max, bookInfos
}

//func (hy HuanYue) Parser() {
//	var startTop = 1 //解析页数
//
//	mongo := orm.Mongo{}
//
//	maxpage, _ := hy.QuanBenTop(1)
//	website := mongo.WebSite("http://www.huanyue123.com")
//	if website.WebsiteURL != "" {
//		if maxpage != website.LastTop {
//			startTop = website.LastTop
//		} else { //解析完成
//			return
//		}
//	}
//
//	for ; startTop <= maxpage; startTop++ {
//		hy.Top(startTop)
//		time.Sleep(time.Second)
//	}
//}

func (hy HuanYue) Top(top int) {
	log.Print("huanyue [Top]: ", top)
	_, books := hy.QuanBenTop(top)
	for _, book := range books {
		bookinfo, chapters := hy.Book(book.Source)            //每本书籍信息 章节信息
		identify := util.MD5(bookinfo.Domain + bookinfo.Name) //网站+书名 hash
		bookName := bookinfo.Name

		xorm := orm.XOrm{}
		exist := xorm.BookExist(identify)
		if !exist { //本地没有保存该书
			var book = orm.Book{
				Identifier:  identify,
				Domain:      bookinfo.Domain,
				Name:        bookName,
				Cover:       bookinfo.Cover,
				Source:      "crawler",
				Describe:    book.Describe,
				Author:      bookinfo.Author,
				Type:        bookinfo.Type,
				Last_update: bookinfo.Last_update,
				Language:    "zh",
				Source_ctr:  1,
				Ctr:         1,
				Score:       1,
			}

			if book.Name == "" || book.Cover == "" || book.Describe == "" || book.Author == "" {
				log.Println("HuanYue book is null")
				continue
			} else {
				xorm.Insert(book)
				log.Print("HuanYue Parser", book)
			}
		}

		for key, chapter := range chapters {
			exist := xorm.ChapterExist(identify, chapter.Title)
			if !exist { //本地没有保存该章节
				chapterDetail := hy.Chapter(chapter.Source)

				var mChapter = orm.Chapter{
					Identifier: identify,
					Idx:        key + 1,
					Idx_name:   string(chapterDetail.Index),
					Title:      chapterDetail.Title,
					Content:    chapterDetail.Content,
					Source:     "crawler",
					Domain:     chapterDetail.Domain,
				}

				if mChapter.Title == "" {
					log.Println("HuanYue chapter is null", chapterDetail)
					continue
				} else {
					xorm.Insert(mChapter)
					log.Print("HuanYue Parser", mChapter)
				}
			}
		}
	}
}

// 单本书籍
func (hy HuanYue) BookAll(url string) {
	book, chapters := hy.Book(url)

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
			chapter := hy.Chapter(simpleChapter.Domain)

			transIndexName := net.Translate(chapter.Idx_name)
			transTitle := net.Translate(chapter.Title)
			transContent := net.Translate(chapter.Content)

			// 章节翻译失败
			if transIndexName == "" || transTitle == "" || transContent == "" {
				log.Print("[小说章节信息翻译失败]", simpleChapter, simpleChapter.Domain, transIndexName, transTitle, transContent)

				if transIndexName == "" {
					log.Print("[章节名称为空]", chapter.Idx_name, "==", transIndexName)
				} else if transTitle == "" {
					log.Print("[章节标题为空]", chapter.Title, "==", transTitle)
				} else if transContent == "" {
					log.Print("[章节内容为空]", chapter.Content, "==", transContent)
				}
			} else {
				chapter.Idx = index
				chapter.Identifier = identify
				chapter.Idx_name = transIndexName
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
