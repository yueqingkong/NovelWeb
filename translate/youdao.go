package translate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
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

var YD YouDao

func init() {
	YD = YouDao{
		Url: "http://fanyi.youdao.com",
	}
}

func (youDao YouDao) Translate(source string) string {
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

	param := Param{
		Api:   api,
		Value: values.Encode(),
		Head:  headers,
	}

	var back = Post(param)

	var result YouDaoResult
	err := json.Unmarshal([]byte(back), &result)
	if err != nil {
		log.Fatal(err)
	}
	return result.TranslateResult[0][0].Tgt
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
