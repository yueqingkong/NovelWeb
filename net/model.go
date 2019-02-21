package net

type UpFileResult struct {
	Code int `json:"code"`
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}

type ProxyInfo struct {
	IP       string  `json:"ip"`
	Port     int     `json:"port"`
	Location string  `json:"location"`
	Source   string  `json:"source"`
	Speed    float32 `json:"speed"`
}

type TranslateReq struct {
	Text string `json:"text"`
}

type TranslateRes struct {
	Code int `json:"code"`
	Data string `json:"data"`
}

type BookRes struct {
	Code int `json:"code"`
	Message string `json:"message"`
}


type ChapterRes struct {
	Code int `json:"code"`
	Message string `json:"message"`
}