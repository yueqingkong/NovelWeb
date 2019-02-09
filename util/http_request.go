package util

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func PostForm(url string, bts []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bts))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	log.Print(resp.Header)
	return body, nil
}

func PostJson(url string, bts []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bts))
	req.Header.Set("Content-Type", "application/json;;charset=UTF8")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
