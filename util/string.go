package util

import (
	"strconv"
	"strings"
)

// 拆分文章 章节/标题
func SepatateTitle(title string) (string, string) {
	title = strings.TrimSpace(title)
	var arr = make([]string, 2)

	var index = strings.Index(title, " ")
	if index > 0 {
		arr[0] = title[0:index]
		arr[1] = title[index:]
	} else {
		arr[0] = ""
		arr[1] = title
	}
	return arr[0], arr[1]
}

// string 转 int
func StringToInt(str string) int {
	var value int
	i, err := strconv.Atoi(str)
	if err != nil {
		value = 0
	} else {
		value = i
	}
	return value
}
