package main

import (
	"NovelWeb/router"
	"log"
)

func main() {
	log.Print("start...")
	router.HttpServer().Run(":8090")
}
