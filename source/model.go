package source

type IBase interface {
	Top(tp int) []BookInfo                                    //排行榜
	SearchKeyword(keyword string, page int) (int, []BookInfo) //关键字查询
	All(url string) (BookInfo, []BookInfo)                    //所有章节
	Chapter(url string) ChapterDetail                         //章节详情
}

type BookInfo struct {
	Domain   string      //来源
	Name     string      //书名
	Cover    string      //封面
	Source   string      //链接
	Describe string      //简介
	Author   string      //作者
	Type     string      //类型(都市言情...)
	Update   string      //最后更新时间
	Chapter  ChapterInfo //最后更新章节
	Language string      //语言
}

// 列表章节
type ChapterInfo struct {
	Index  string //章节
	Title  string //标题
	Source string //原文链接
}

// 章节详情
type ChapterDetail struct {
	Index   string //章节
	Title   string //标题
	Content string //内容
	Last    string //上一章
	Next    string //下一章
	Update  string //更新时间
	Domain string //原文链接
}
