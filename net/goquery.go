package net

import (
	"Novel/net"
	"Novel/util"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func GoQuery(url string) *goquery.Document {
	var doc *goquery.Document
	var err error

	log.Print(url)
	html := net.Get(url)

	// 兼容GBK
	html = util.GbkToUtf8(html)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("GoQuery err: ", err)
	}
	return doc
}
