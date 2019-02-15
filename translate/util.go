package translate

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Translate interface {
	Translate(source string)
	Salt() string
	Sign(s string) string
}

type Param struct {
	Api   string
	Value string
	Head  map[string]string
}

func Get(param Param) string {
	return request("GET", param)
}

func Post(param Param) string {
	return request("POST", param)
}

func request(method string, param Param) string {
	client := http.Client{}

	api := param.Api
	value := param.Value
	req, err := http.NewRequest(method, api, strings.NewReader(value))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	for k, v := range param.Head {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	value = string(body)
	//log.Print(value)
	return value
}
