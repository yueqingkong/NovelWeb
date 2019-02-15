package controller

import (
	"github.com/gin-gonic/gin"
)

type Novel struct {
}

func (_ *Novel) Home(context *gin.Context) {
	context.HTML(200, "kline.html", "")
}
