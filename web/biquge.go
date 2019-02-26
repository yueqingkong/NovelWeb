package web

type Biquge struct {
	Source string
}

func NewBiquge() Biquge {
	return Biquge{
		Source: "http://www.xbiquge.la",
	}
}


