package util

import (
	"bytes"
	"github.com/go-resty/resty"
	"io"
	"log"
	"os"
)

func FileDownload(local string, url string) {
	resp, err := resty.R().Get(url)
	if err != nil {
		log.Print(err)
	}

	out, _ := os.Create(local)
	_, err = io.Copy(out, bytes.NewReader(resp.Body()))
	if err != nil {
		log.Print(err)
	}
}
