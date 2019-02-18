package main

import (
	"NovelWeb/source"
	"log"
)

func main() {
	log.Print("start...")

	//err := router.HttpServer().Run(":8090")
	//if err != nil {
	//	log.Print(err)
	//}

	links := []string{"http://www.huanyue123.com/book/50/50083/"}
	hy := source.HuanYue{}

	for _, link := range links {
		hy.BookAll(link)
	}
}
