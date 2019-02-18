package net

import (
	"Novel/net"
	"NovelWeb/orm"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

var uri = "http://119.28.68.41:8989"

func UpBook(ebook orm.Book) string {
	api := "/novel/data/crawler/v1/books"
	return net.Post(uri+api, ebook)
}

func UpChapter(echapter orm.Chapter) string {
	api := "/novel/data/crawler/v1/chapters"
	return net.Post(uri+api, echapter)
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

func Get(url string) string {
	return request("GET", url, nil)
}

func Post(url string, value interface{}) string {
	return request("POST", url, value)
}

func request(method string, request string, content interface{}) string {
	//client := http.Client{}

	// 代理地址
	//userProxy := UserProxy()
	//if userProxy.IP != "" {
	//	log.Print("http: 代理地址 ", userProxy)
	//	temp := fmt.Sprintf("http://%s:%s", "111.177.166.132", "9999")
	//	urlproxy, _ := url.Parse(temp)
	//	transport := &http.Transport{Proxy: urlproxy}
	//
	//	client.Transport = &http.Transport{
	//		Proxy: http.ProxyURL(transport),
	//	}
	//}
	temp := fmt.Sprintf("http://%s:%s", "109.69.75.5", "46347")
	urlproxy, _ := url.Parse(temp)
	transport := &http.Transport{Proxy: http.ProxyURL(urlproxy)}
	//client.Transport = transport

	client := &http.Client{
		Transport: transport,
	}

	log.Print("url: ", request)

	var arr []byte
	var err error

	if content != nil {
		arr, err = json.Marshal(content)
		log.Print("json: ", string(arr))
		if err != nil {
			log.Println(err)
			return ""
		}
	}

	req, err := http.NewRequest(method, request, bytes.NewReader(arr))
	if err != nil {
		log.Println("NewRequest: ", err)
		return ""
	}

	//req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("user-agent", UserAgent())
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

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
