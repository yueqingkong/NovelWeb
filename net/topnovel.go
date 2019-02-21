package net

import (
	"NovelWeb/orm"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"log"
)

var bookuri = "http://119.28.68.41:9898"
var fileuri = "http://119.28.68.41:8989"

// 上传小说信息
func UploadBook(book orm.Book) BookRes {
	api := fmt.Sprintf("%s%s", bookuri, "/novel/data/crawler/v1/books")

	resp, _ := resty.R().
		SetBody(book).
		Post(api)

	var bookRes BookRes

	err := json.Unmarshal(resp.Body(), &bookRes)
	if err != nil {
		log.Print(err)
	}
	return bookRes
}

// 上传章节
func UploadChapter(chapter orm.Chapter) ChapterRes {
	api := fmt.Sprintf("%s%s", bookuri, "/novel/data/crawler/v1/chapters")

	resp, _ := resty.R().
		SetBody(chapter).
		Post(api)

	var chapterRes ChapterRes
	err := json.Unmarshal(resp.Body(), &chapterRes)
	if err != nil {
		log.Print(err)
	}
	return chapterRes
}

// 文本翻译
func Translate(source string) string {
	if source == "" {
		return ""
	}

	api := "http://47.52.131.191:3013/api/baidu"
	resp, _ := resty.R().
		SetBody(TranslateReq{Text: source}).
		Post(api)

	var result TranslateRes
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Print(err)
	}

	log.Print("[翻译] 原文: ", source, "结果: ", string(resp.Body()))
	return result.Data
}

// 上传文件
func UploadFile(path string, result interface{}) {
	api := fmt.Sprintf("%s%s", fileuri, "/novel/api/upload")

	resp, _ := resty.R().
		SetFile("file", path).
		Post(api)

	err := json.Unmarshal(resp.Body(), result)
	if err != nil {
		log.Print(err)
	}
}
