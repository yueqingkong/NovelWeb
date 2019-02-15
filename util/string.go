package util

import "strings"

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
