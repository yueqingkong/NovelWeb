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

// 删除文件
func FileRemove(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Print("[FileRemove] ", err, "  path: ", path)
	} else {
		log.Print("[FileRemove] ", path)
	}
}
