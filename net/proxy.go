package net

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
)

//数据来源: https://www.xicidaili.com/

type Proxy struct {
	IP   string
	Port string
}

var userAgent = []string{
	"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
}

func UserAgent() string {
	random := rand.Intn(len(userAgent))
	return userAgent[random]
}

var ProxyArr []Proxy

func UserProxy() Proxy {
	var index Proxy

	if len(ProxyArr) != 0 {
		random := rand.Intn(len(ProxyArr))
		index = ProxyArr[random]
	}
	return index
}

func init() {
	Parser()
}

func Parser() {
	var url = "https://www.xicidaili.com/wn/1"
	doc := GoQuery(url)

	var maxPage int
	doc.Find("div.pagination").Find("a").Each(func(i int, sec *goquery.Selection) {
		value := sec.Text()
		b, err := regexp.Match("[0-9]+", []byte(value))
		if err != nil {
			log.Print(err)
		}

		if b {
			page, err := strconv.Atoi(value)
			if err != nil {
				log.Print(err)
			} else {
				maxPage = page
			}
		}
	})

	log.Print("proxy:   maxpage ", maxPage)
	//for index := 1; index <= 1; index++ {
	//	PageDetail(index)
	//}
}

func PageDetail(index int) {
	api := "https://www.xicidaili.com/wn/%d"
	url := fmt.Sprintf(api, index)

	doc := GoQuery(url)
	doc.Find("div#body").Find("table#ip_list").Find("tr[class]").Each(func(i int, sec *goquery.Selection) {
		var proxy Proxy
		sec.Find("td").Each(func(i int, sectd *goquery.Selection) {

			if i == 1 { //ip
				proxy.IP = sectd.Text()
			}

			if i == 2 { //port
				proxy.Port = sectd.Text()
			}
		})

		log.Println(proxy)
		ProxyArr = append(ProxyArr, proxy)
	})
	log.Println(ProxyArr)
}

func ProxyIp() ProxyInfo {
	client := http.Client{}
	res, err := client.Get("http://localhost:8090/get")
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("body: ", err)
	}

	var info ProxyInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Print(err)
	}
	return info
}
