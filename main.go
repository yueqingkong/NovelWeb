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

	hy := source.HuanYue{}
	hy.BookAll("http://www.huanyue123.com/book/52/52260/")
}
