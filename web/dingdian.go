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

type DingDian struct {
	Url string
}

func NewDingDian() DingDian {
	return DingDian{
		Url: "https://www.x23us.com",
	}
}

// 网站小说下载
func (dd DingDian) Pull() {
	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("/class/%d_", i) + "%d.html"
		log.Print("[url]", url)
		dd.normal(url)
	}

	// 全本
	dd.normal("/quanben/%d")
}

// 通用
func (dd DingDian) normal(format string) {
	api := fmt.Sprintf("%s%s", dd.Url, fmt.Sprintf(format, 1))
	doc := net.GoQuery(api, true)

	lastA := doc.Find("div.bdsub").Find("dd.pages").Find("div.pagelink").Find("a.last")
	totalPage := util.StringToInt(lastA.Text())
	for i := 1; i <= totalPage; i++ {
		api = fmt.Sprintf("%s%s", dd.Url, fmt.Sprintf(format, i))
		doc = net.GoQuery(api, true)

		doc.Find("div.bdsub").Find("tbody").Find("tr[bgcolor='#FFFFFF']").Each(func(i int, sec *goquery.Selection) {
			sec.Find("td").Each(func(i int, sec *goquery.Selection) {
				if i == 0 {
					bookUrl := sec.Find("a").First().AttrOr("href", "")
					dd.BookAll(bookUrl)
				}
			})
		})
	}
}

///////////////////////////////////////////////////   功能  /////////////////////////////////////////////////////
// 小说信息及列表
func (dd DingDian) BookAll(url string) {
	book, chapters := dd.book(url)

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
			chapter := dd.chapter(simpleChapter.Domain)

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

// 小说
func (dd DingDian) book(url string) (orm.Book, []orm.Chapter) {
	var book orm.Book
	var chapters []orm.Chapter

	book.Domain = url

	doc := net.GoQuery(url, true)
	doc.Find("div.bdsub").Find("dd").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			h1Txt := sec.Find("h1").Text()
			name := strings.Split(h1Txt, " ")[0]
			book.Name = name
		} else if i == 1 {
			// 封面
			src := sec.Find("a.hst").Find("img").AttrOr("src", "")
			src = dd.Url + src
			book.Cover = src

			// 小说属性
			sec.Find("tbody").Find("tr").Each(func(i int, sec *goquery.Selection) {
				if i == 0 {
					sec.Find("td").Each(func(j int, sec *goquery.Selection) {
						if j == 0 {
							bookType := sec.Find("a").Text()
							book.Type = bookType
						} else if j == 1 {
							author := sec.Text()
							book.Author = author
						} else if j == 2 {
							statusTxt := sec.Text()

							var status string
							if strings.Contains(statusTxt, "连载中") {
								status = "2"
							} else {
								status = "1"
							}
							book.Status = status
						}
					})
				} else if i == 1 {
					sec.Find("td").Each(func(j int, sec *goquery.Selection) {
						if j == 2 {
							update := sec.Text()
							book.Last_update = update
						}
					})
				}
			})

			chapterUrl := sec.Find("p.btnlinks").Find("a").AttrOr("href", "")
			chapterDoc := net.GoQuery(chapterUrl, true)
			chapterDoc.Find("div.bdsub").Find("tbody").Find("td.L").Each(func(i int, sec *goquery.Selection) {
				a := sec.Find("a")
				txt := a.Text()
				index, title := util.TitleSepatate(txt)
				href := a.AttrOr("href", "")
				url := chapterUrl + href

				chapter := orm.Chapter{
					Idx:      i + 1,
					Idx_name: index,
					Domain:   url,
					Title:    title,
				}
				chapters = append(chapters, chapter)
			})
		} else if i == 3 {
			sec.Find("p").Each(func(i int, sec *goquery.Selection) {
				if i == 1 {
					describe := sec.Text()
					book.Describe = describe
				}
			})
		}
	})
	return book, chapters
}

// 章节
func (dd DingDian) chapter(url string) orm.Chapter {
	doc := net.GoQuery(url, true)
	var chapter orm.Chapter

	bdsub := doc.Find("div.bdsub")

	// 标题
	bdsub.Find("dd").Each(func(i int, sec *goquery.Selection) {
		if i == 0 {
			h1Txt := sec.Find("h1").Text()
			index, title := util.TitleSepatate(h1Txt)
			chapter.Idx_name = index
			chapter.Title = title
		} else if i == 2 {
			content := strings.Replace(sec.Text(), "顶 点 小 说 Ｘ ２３ Ｕ Ｓ．Ｃ ＯＭ", "", -1)
			chapter.Content = content
		}
	})
	return chapter
}
