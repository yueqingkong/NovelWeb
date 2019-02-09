package router

import (
	"NovelWeb/controller"
	"github.com/gin-gonic/gin"
)

func v1(engine *gin.Engine) {
	novel := new(controller.Novel)

	v1 := engine.Group("/novel/v1/")
	{
		v1.GET("/home", novel.Home)
	}
}

func HttpServer() *gin.Engine {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/**")

	v1(engine)
	return engine
}
