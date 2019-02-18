package translate

import (
	"NovelWeb/net"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type YouDao struct {
	Url string
}

type YouDaoResult struct {
	TranslateResult [][]struct {
		Tgt string `json:"tgt"`
		Src string `json:"src"`
	} `json:"translateResult"`
	ErrorCode   int    `json:"errorCode"`
	Type        string `json:"type"`
	SmartResult struct {
		Entries []string `json:"entries"`
		Type    int      `json:"type"`
	} `json:"smartResult"`
}

func NewYouDao() YouDao {
	return YouDao{
		Url: "http://fanyi.youdao.com",
	}
}

// 翻译，支持超过翻译字数限制
func (youdao YouDao) Translate(source string) string {
	reg := regexp.MustCompile("\\s+")
	source = reg.ReplaceAllString(source, "")

	var content bytes.Buffer

	for i := 0; i < len(source); i += 100 {
		var endPoint int
		if i+100 < len(source) {
			endPoint = i + 100
		} else {
			endPoint = len(source)
		}

		var part string
		if endPoint == len(source)-1 {
			part = youdao.TranslateLimit(source[i:])
		} else {
			part = youdao.TranslateLimit(source[i:endPoint])

			log.Print("***  ", source[i:endPoint])
			log.Print("***  ", source[endPoint:])
		}
		content.WriteString(part)
	}
	return content.String()
}

func (youDao YouDao) TranslateLimit(source string) string {
	var api = "http://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule"

	salt := youDao.Salt()
	sign := youDao.Sign("fanyideskweb" + source + salt + "ebSeFb%=XZ%T[KZ)c(sy!")

	values := url.Values{
		"i":           {source},
		"from":        {"AUTO"},
		"to":          {"AUTO"},
		"smartresult": {"dict"},
		"client":      {"fanyideskweb"},
		"salt":        {salt},
		"sign":        {sign},
		"doctype":     {"json"},
		"version":     {"2.1"},
		"keyfrom":     {"fanyi.web"},
		"action":      {"FY_BY_CLICKBUTTION"},
		"typoResult":  {"typoResult"}}

	var headers = make(map[string]string)
	headers["Cookie"] = "OUTFOX_SEARCH_USER_ID=1799185238@10.169.0.83;"
	headers["Referer"] = "http://fanyi.youdao.com/"

	var back = net.Post(api, headers, values.Encode())

	var result YouDaoResult
	err := json.Unmarshal([]byte(back), &result)
	if err != nil {
		log.Print("Unmarshal", err)
	}

	resultTranslate := result.TranslateResult[0][0].Tgt
	log.Print("[翻译]", source, "   [结果]", resultTranslate)
	return resultTranslate
}

func (youDao YouDao) Salt() string {
	var mils = time.Now().UnixNano() / 1e6
	var random = rand.Intn(10)
	return strconv.FormatInt(mils, 10) + strconv.Itoa(random)
}

func (youDao YouDao) Sign(source string) string {
	data := []byte(source)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
