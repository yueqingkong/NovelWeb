package net

import (
	"github.com/go-resty/resty"
	"log"
)

func Get(url string, header map[string]string) string {
	resp, err := resty.SetHeaders(header).R().Get(url)
	if err != nil {
		log.Print(err)
	}

	return resp.String()
}
