package util

import (
	"regexp"
	"strconv"
	"strings"
)

// 拆分文章 章节/标题
func TitleSepatate(title string) (string, string) {
	title = strings.TrimSpace(title)

	// 去除标题中的 (*)
	reg, _ := regexp.Compile("\\([\\s\\S]*\\)")
	title = reg.ReplaceAllString(title, "")

	var arr = make([]string, 2)

	reg = regexp.MustCompile("第([0-9]+|[\u4e00-\u9fa5]+)章")
	idx := reg.FindAllString(title, 1)
	if len(idx) > 1 {
		arr[0] = idx[0]
		arr[1] = strings.Replace(title, idx[0], "", -1)
	} else {
		// 按 章|节 拆分
		reg := regexp.MustCompile("(章|节)[\\s]*")
		index := reg.Split(title, -1)
		if len(index) > 1 {
			arr[0] = index[0] + "章"
			arr[1] = index[1]
		} else {
			arr[0] = title
			arr[1] = title
		}
	}
	return arr[0], arr[1]
}

// 章节内容处理
// 去除html标签 章节错误,点此举报(免注册)
func ChapterFilter(content string) string {
	reg := regexp.MustCompile("(<[^>]+>)|(章节错误,点此举报\\(免注册\\))")
	content = reg.ReplaceAllString(content, "")

	return content
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

func IntToString(valye int) string {
	return strconv.Itoa(valye)
}
