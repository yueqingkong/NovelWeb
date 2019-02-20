package net

import (
	"NovelWeb/orm"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"log"
)

var uri = "http://119.28.68.41:8989"

// 上传小说信息
func UploadBook(book orm.Book) string {
	api := "/novel/data/crawler/v1/books"
	resp, _ := resty.R().
		SetBody(book).
		Post(api)

	return resp.String()
}

// 上传章节
func UploadChapter(chapter orm.Chapter) string {
	api := "/novel/data/crawler/v1/chapters"
	resp, _ := resty.R().
		SetBody(chapter).
		Post(api)

	return resp.String()
}

// 上传文件
func UploadFile(path string, result interface{}) {
	api := fmt.Sprintf("%s%s", uri, "/novel/api/upload")

	resp, _ := resty.R().
		SetFile("file", path).
		Post(api)

	err := json.Unmarshal(resp.Body(), result)
	if err != nil {
		log.Print(err)
	}
}
