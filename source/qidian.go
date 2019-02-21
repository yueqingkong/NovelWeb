package source

import "NovelWeb/orm"

type QiDian struct {
	Url string
}

func NewQidian() QiDian {
	return QiDian{Url: "https://www.qidian.com/"}
}

// 热门
func (qi QiDian) Hot() []orm.Book {
	return nil
}
