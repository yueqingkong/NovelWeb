package net

type UpFileResult struct {
	Code int `json:"code"`
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}
