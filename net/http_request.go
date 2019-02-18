package net

import (
	"NovelWeb/orm"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var uri = "http://119.28.68.41:8989"

func UpBook(ebook orm.Book) string {
	api := "/novel/data/crawler/v1/books"
	return Post(uri+api, nil, ebook)
}

func UpChapter(echapter orm.Chapter) string {
	api := "/novel/data/crawler/v1/chapters"
	return Post(uri+api, nil, echapter)
}

// 上传文件
func UploadFile(path string, result interface{}) {
	api := "/novel/api/upload"

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Print(err)
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		log.Print(err)
	}
	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		log.Print(err)
	}

	request, err := http.NewRequest("POST", uri+api, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	resp, err := client.Do(request)
	resultContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	err = json.Unmarshal(resultContent, result)
	if err != nil {
		log.Print(err)
	}
}

func Get(url string, header map[string]string) string {
	return request("GET", url, header, nil)
}

func Post(url string, header map[string]string, value interface{}) string {
	return request("POST", url, header, value)
}

func request(method string, api string, header map[string]string, content interface{}) string {
	proxyIp := ProxyIp()
	temp := fmt.Sprintf("http://%s:%d", proxyIp.IP, proxyIp.Port)
	urlproxy, _ := url.Parse(temp)
	transport := &http.Transport{Proxy: http.ProxyURL(urlproxy)}

	client := &http.Client{
		Transport: transport,
	}

	var reader io.Reader
	var err error

	if content != nil {
		switch content.(type) {
		case string:
			reader = strings.NewReader(content.(string))
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

	req.Header.Set("user-agent", UserAgent())
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

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
