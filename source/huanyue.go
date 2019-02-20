package source

import (
	"NovelWeb/net"
	"NovelWeb/orm"
	"NovelWeb/translate"
	"NovelWeb/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

var (
	url = "http://www.huanyue123.com"
)

type HuanYue struct {
}

func NewHuanYue() HuanYue {
	return HuanYue{}
}

// 关键字查询
func (huan HuanYue) SearchKeyword(keyword string, page int) (int, []BookInfo) {
	var bookInfos = make([]BookInfo, 0)

	var api = "/modules/article/search.php?searchkey=%s&page=%d"
	var url = fmt.Sprintf(url+api, util.Utf8ToGbk(keyword), page)
	log.Print("url: ", url)

	var doc = net.GoQuery(url)
	doc.Find("tbody").Find("tr").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			return
		}

		var bookInfo BookInfo
		sec.Find("td").Each(func(i int, sec *goquery.Selection) {
			if i == 0 {
				a := sec.Find("a")
				name := a.Text()
				source := a.AttrOr("href", "")

				bookInfo.Name = name
				bookInfo.Source = source
			} else if i == 1 {
				a := sec.Find("a")
				index, title := util.SepatateTitle(a.Text())
				href := a.AttrOr("href", "")

				bookInfo.Chapter = ChapterInfo{
					Index:  index,
					Title:  title,
					Source: href,
				}
			} else if i == 2 {
				author := sec.Text()
				bookInfo.Author = author
			} else if i == 3 {

			} else if i == 4 {
				update := sec.Text()
				bookInfo.Update = update
			}
		})

		bookInfos = append(bookInfos, bookInfo)
	})

	pages := doc.Find("div.pagelink").Find("em").Text()
	max, _ := strconv.Atoi(strings.Split(pages, "/")[1])
	return max, bookInfos
}

// 书本简介及章节列表
func (huan HuanYue) Book(url string) (BookInfo, []ChapterInfo) {
	var bookInfo BookInfo
	var chapterInfos []ChapterInfo

	bookInfo.Domain = url
	var doc = net.GoQuery(url)

	booktype := doc.Find("div.main").Find("a").Next().Text()
	bookInfo.Type = booktype

	doc.Find("div.book_info").Find("div").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			img := sec.Find("img").AttrOr("src", "")
			bookInfo.Cover = img
		} else if i == 1 {
			name := sec.Find("h1").Text()
			bookInfo.Name = name

			author := strings.Split(sec.Find("span.red").Text(), "：")[1]
			bookInfo.Author = author

			options := sec.Find("h3").Find("div.options")
			update := strings.Split(options.Find("span.hottext").Text(), "：")[1]
			bookInfo.Update = update

			a := options.Find("a")
			index, title := util.SepatateTitle(a.Text())
			href := a.AttrOr("href", "")
			bookInfo.Chapter = ChapterInfo{
				Index:  index,
				Title:  title,
				Source: url + href,
			}

			describe := sec.Find("h3").Text()
			bookInfo.Describe = describe
		}
	})

	doc.Find("div.book_list").Find("li").Each(func(i int, sec *goquery.Selection) {
		var chapterInfo ChapterInfo
		a := sec.Find("a")
		index, title := util.SepatateTitle(a.Text())
		href := a.AttrOr("href", "")

		chapterInfo = ChapterInfo{
			Index:  index,
			Title:  title,
			Source: href,
		}
		chapterInfos = append(chapterInfos, chapterInfo)
	})

	bookInfo.Type = "unknow"
	return bookInfo, chapterInfos
}

// 章节详情
func (huan HuanYue) Chapter(url string) ChapterDetail {
	var chapterDetail ChapterDetail

	chapterDetail.Domain = url

	var doc = net.GoQuery(url)
	body := doc.Find("div.wrapper_main")
	index, title := util.SepatateTitle(body.Find("div.h1title").Text())
	chapterDetail.Index = index
	chapterDetail.Title = title

	body.Find("div.chapter_Turnpage_1").Find("a").Each(func(i int, sec *goquery.Selection) {
		href := sec.AttrOr("href", "")
		if i == 0 {
			chapterDetail.Last = url + href
		} else if i == 4 {
			chapterDetail.Next = url + href
		}
	})
	content := body.Find("div#htmlContent").Text()

	// 章节错误,点此举报(免注册)
	content = strings.Replace(content, "章节错误,点此举报(免注册)", "", -1)
	chapterDetail.Content = content

	return chapterDetail
}

// 每页全本的列表
func (huan HuanYue) QuanBenTop(tp int) (int, []BookInfo) {
	var bookInfos []BookInfo
	var api = "/book/quanbu/default-0-0-0-0-2-0-%d.html"
	var url = fmt.Sprintf(url+api, tp)

	var doc = net.GoQuery(url)
	doc.Find("div.sitebox").Find("dl").Each(func(i int, sec *goquery.Selection) {
		var bookInfo BookInfo
		bookInfo.Domain = url
		bookInfo.Type = "unknow"

		img := sec.Find("dt").Find("a").Find("img")
		bookInfo.Cover = img.AttrOr("src", "")

		sec.Find("dd").Each(func(i int, sec *goquery.Selection) {
			if i == 0 {
				update := sec.Find("span").Text()
				bookInfo.Update = update

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
				a := sec.Find("a")
				index, title := util.SepatateTitle(a.Text())
				href := a.AttrOr("href", "")
				bookInfo.Chapter = ChapterInfo{
					Index:  index,
					Title:  title,
					Source: url + href,
				}
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
				Last_update: bookinfo.Update,
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
	baiduTrans := translate.NewBaiDu()

	// 书本信息
	identify := util.MD5(book.Domain + book.Name)
	log.Print("identify", identify)
	if xorm.BookExist(identify) {
		log.Print("[小说已存在]", book.Name)
	} else {
		filePath := identify + ".jpg"
		util.FileDownload(filePath, book.Cover)

		var fileResult net.UpFileResult
		net.UploadFile(filePath, &fileResult)

		if fileResult.Code == 2000 {
			book.Cover = fileResult.Data.URL
		}

		transBookName := baiduTrans.Translate(book.Name)
		transBookDesc := baiduTrans.Translate(book.Describe)
		transBookAuthor := baiduTrans.Translate(book.Author)
		transBookType := baiduTrans.Translate(book.Type)

		// 翻译失败
		if transBookName == "" || transBookDesc == "" || transBookAuthor == "" {
			log.Print("[小说信息翻译失败]", book.Source, transBookName, transBookDesc, transBookAuthor)
		} else {
			ormBook := orm.Book{
				Identifier:  identify,
				Name:        transBookName,
				Domain:      book.Domain,
				Cover:       book.Cover,
				Source:      book.Source,
				Describe:    transBookDesc,
				Author:      transBookAuthor,
				Type:        transBookType,
				Last_update: book.Update,
				Language:    book.Language,
			}

			log.Print("[小说]", ormBook)
			xorm.Insert(ormBook)
		}
	}

	// 章节
	for index, chapter := range chapters {
		if xorm.ChapterExist(identify, string(index)) {
			log.Print("[章节已存在]", book.Name, chapter.Title)
		} else {
			chapterDetail := hy.Chapter(chapter.Source)

			transChapterIndex := baiduTrans.Translate(chapter.Index)
			transChapterTitle := baiduTrans.Translate(chapter.Title)
			transChapterContent := baiduTrans.Translate(chapterDetail.Content)

			// 章节翻译失败
			if transChapterIndex == "" || transChapterTitle == "" || transChapterContent == "" {
				log.Print("[小说章节信息翻译失败]", chapter.Source, transChapterIndex, transChapterTitle, transChapterContent)
			} else {
				ormChapter := orm.Chapter{
					Identifier: identify,
					Idx:        index,
					Idx_name:   transChapterIndex,
					Title:      transChapterTitle,
					Content:    transChapterContent,
					Source:     chapter.Source,
					Domain:     book.Domain,
					UpTime:     chapterDetail.Update,
				}

				log.Print("[章节]", book.Name, ormChapter)
				xorm.Insert(ormChapter)
			}
		}
	}
}
