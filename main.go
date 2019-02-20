package main

import (
	"NovelWeb/net"
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

	//links := []string{"http://www.huanyue123.com/book/50/50083/",
	//	"http://www.huanyue123.com/book/52/52260/",
	//	"http://www.huanyue123.com/book/49/49221/",
	//	"http://www.huanyue123.com/book/5/5544/"}
	hy := source.NewHuanYue()

	for _, link := range links {
		hy.BookAll(link)
	}

}

func TestTranslate()  {
	log.Print(net.Translate(" 老鹰吃小鸡结果"))
}
