package net

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Get(url string, header map[string]string) string {
	return request("GET", url, header, nil)
}

func Post(url string, header map[string]string, value interface{}) string {
	return request("POST", url, header, value)
}

func request(method string, api string, header map[string]string, content interface{}) string {
	//proxyIp := ProxyIp()
	//temp := fmt.Sprintf("http://%s:%d", proxyIp.IP, proxyIp.Port)
	//urlproxy, _ := url.Parse(temp)
	//transport := &http.Transport{Proxy: http.ProxyURL(urlproxy)}
	//
	//client := &http.Client{
	//	Transport: transport,
	//}


	client := &http.Client{}

	var reader io.Reader
	var err error

	if content != nil {
		switch content.(type) {
		case string:
			str:= content.(string)
			reader =strings.NewReader(str)
			break
		default:
			arr, _ := json.Marshal(content)
			reader = bytes.NewReader(arr)
		}
	}

	req, err := http.NewRequest(method, api, reader)
	if err != nil {
		log.Println("NewRequest: ", err)
		return ""
	}

	// req.Header.Set("user-agent", UserAgent())
	// req.Header.Set("content-type", "application/x-www-form-urlencoded")

	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll: ", err)
		return ""
	}

	//var html = util.GbkToUtf8(string(body))
	return string(body)
}
