package net

import (
	"NovelWeb/util"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

// b: 是否UTF-8格式
func GoQuery(url string, b bool) *goquery.Document {
	var doc *goquery.Document
	var err error

	log.Print(url)
	html := Get(url, nil)

	// 兼容GBK
	if b {
		html = util.GbkToUtf8(html)
	}
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("GoQuery err: ", err)
	}
	return doc
}
