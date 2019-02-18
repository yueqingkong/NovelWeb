package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func FileDownload(local string, url string) {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	out, _ := os.Create(local)
	_, err := io.Copy(out, bytes.NewReader(body))
	if err != nil {
		log.Print(err)
	}
}
