package main

import (
	"NovelWeb/translate"
	"log"
)

func main() {
	log.Print("start...")

	//err := router.HttpServer().Run(":8090")
	//if err != nil {
	//	log.Print(err)
	//}

	//links := []string{"http://www.huanyue123.com/book/50/50083/"}
	//hy := source.HuanYue{}
	//
	//for _, link := range links {
	//	hy.BookAll(link)
	//}

	trans:=translate.BaiDu{}
	//log.Print(trans.TranslateLimit("我说是"))


	content:=`
	更新时间：2019-02-18 13:32:44最新章节：征求一下意见

	《全球高武》简介：
	    今日头条——
	    “大马宗师突破九品，征战全球！”
	    “小马宗师问鼎至高，横扫欧亚！”
	    “乔帮主再次出手，疑似九品大宗师境！”
	    “股神宝刀未老，全球宗师榜再入前十！”
	    “……”
	    看着一条条新闻闪现，方平心好累，这剧本不对啊！
	    各位书友要是觉得《全球高武》还不错的话请不要忘记向您qq群和微博里的朋友推荐哦！

	`
	content  = "“大马宗师突破九品，征战全球！”"
	log.Print(trans.TranslateLimit(content))
}
